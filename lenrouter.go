package lenrouter

import (
	"errors"
	"net/http"
	"sync"
)

var (
	errNoGuess = errors.New("no guess")
)

type Endpoint struct {
	Method  string
	Pattern string
	Handler Handle

	parts      []string // parts of the url for /foo/bar/:spam it is []string{"foo", "bar", ":spam"}
	isParam    []bool   // isParam[i] = true if parts[i] is a parameter else false.
	paramCount int      // number of parameters in this endpoint e.g. for /foo/:bar/:spam paramCount = 2.
}

type guess struct {
	endpointIdx int
	checkSet    []int
	instance    string
	slashLocs   []int // locations of slashes in the url path.
}

type router struct {
	endpoints []Endpoint
	paramPool sync.Pool

	// This will need locking because multiple requests could be modifying it concurrently.
	mu           sync.Mutex
	lenToGuesses [][]*guess // url.Path length -> possible guesses
}

func New(maxLen, maxParams int, endpoints ...Endpoint) http.Handler {
	r := &router{
		lenToGuesses: make([][]*guess, maxLen),
		paramPool: sync.Pool{
			New: func() interface{} {
				ps := make(Params, 0, maxParams)
				return &ps
			},
		},
	}

	for _, endpoint := range endpoints {
		endpoint.parts, endpoint.isParam, endpoint.paramCount = parts(endpoint.Pattern)
		r.endpoints = append(r.endpoints, endpoint)
	}

	return r
}

func (r *router) returnToPool(params *Params) {
	if params != nil {
		r.paramPool.Put(params)
	}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, params, err := r.guess(req.URL.Path)
	if err == nil {
		if params != nil {
			handler(w, req, *params)
			r.paramPool.Put(params)
		} else {
			handler(w, req, nil)
		}
		return
	}

	handler, params, err = r.brute(req.URL.Path)
	if err == nil {
		if params != nil {
			handler(w, req, *params)
			r.paramPool.Put(params)
		} else {
			handler(w, req, nil)
		}
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (r *router) brute(s string) (Handle, *Params, error) {
	g := &guess{slashLocs: r.slashLocs(s)}
	for i, p := range r.endpoints {
		params, ok := r.match(p, g, s)
		if ok {
			r.update(len(s), i, s, g.slashLocs)
			return p.Handler, params, nil
		}
	}
	return nil, nil, errNoGuess
}

func (r *router) match(e Endpoint, g *guess, s string) (*Params, bool) {
	parts := e.parts
	isParam := e.isParam
	paramCount := e.paramCount
	slashLocs := g.slashLocs

	if len(parts) != len(slashLocs) {
		return nil, false
	}

	var params *Params

	paramIndex := 0
	for i := 0; i < len(parts); i++ {
		var val string
		if i+1 == len(parts) {
			val = s[slashLocs[i]+1:]
		} else {
			val = s[slashLocs[i]+1 : slashLocs[i+1]]
		}

		if isParam[i] {
			if params == nil {
				params = r.paramPool.Get().(*Params)
				*params = ([]Param)(*params)[:paramCount]
			}
			(*params)[paramIndex] = Param{parts[i], val}
			paramIndex++
		} else {
			if parts[i] != val {
				if params != nil {
					r.paramPool.Put(params)
				}
				return nil, false
			}
		}
	}

	return params, true
}

func (r *router) guess(s string) (Handle, *Params, error) {
	l := len(s)
	guesses := r.lenToGuesses[l]
	for i := 0; i < len(guesses); i++ {
		ok := check(guesses[i].checkSet, guesses[i].instance, s)
		if !ok {
			continue
		}

		params, ok := r.match(r.endpoints[guesses[i].endpointIdx], guesses[i], s)
		if ok {
			return r.endpoints[guesses[i].endpointIdx].Handler, params, nil
		}
	}
	return nil, nil, errNoGuess
}

func parts(p string) ([]string, []bool, int) {
	var piSlash []int
	for pi := 0; pi < len(p); pi++ {
		if p[pi] == '/' {
			piSlash = append(piSlash, pi)
		}
	}
	piSlash = append(piSlash, len(p))

	var parts []string
	var isKey []bool
	keyCount := 0
	for k := 0; k < len(piSlash)-1; k++ {
		part := p[piSlash[k]+1 : piSlash[k+1]]
		if part != "" && part[0] == ':' {
			parts = append(parts, part[1:])
			isKey = append(isKey, true)
			keyCount++
		} else {
			parts = append(parts, part)
			isKey = append(isKey, false)
		}
	}
	return parts, isKey, keyCount
}

// check checks if x and y match at the indexes specified in checkSet
// assumes that len(x) == len(y)
func check(checkSet []int, y, x string) bool {
	for _, i := range checkSet {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func (r *router) slashLocs(s string) []int {
	var locs []int
	for si := 0; si < len(s); si++ {
		if s[si] == '/' {
			locs = append(locs, si)
		}
	}
	return locs
}

func (r *router) update(length, endpoint int, instance string, slashlocs []int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	v := new(guess)

	v.endpointIdx = endpoint
	v.instance = mark(r.endpoints[endpoint].Pattern, []byte(instance))
	v.slashLocs = slashlocs // r.slashLocs(instance)

	r.lenToGuesses[length] = append(r.lenToGuesses[length], v)

	updateCheckSets(r.lenToGuesses[length])
}

// diff returns location of first char where x and y differ.
// parts of the path that are params (start with :) are ignored.
func diff(x, y string) int {
	i := 0
	for i < len(x) {
		if x[i] == ':' {
			for i < len(x) && x[i] != '/' {
				i++
			}
		}
		if i < len(x) && x[i] != y[i] {
			return i
		}
		i++
	}
	return -1
}

func updateCheckSets(xs []*guess) {
	m := make(map[int]map[int]struct{})
	for i := 0; i < len(xs); i++ {
		m[i] = make(map[int]struct{})
		for j := 0; j < len(xs); j++ {
			if d := diff(xs[i].instance, xs[j].instance); d != -1 {
				m[i][d] = struct{}{}
			}
		}
		xs[i].checkSet = xs[i].checkSet[:0]
		for v := range m[i] {
			xs[i].checkSet = append(xs[i].checkSet, v)
		}
	}
}

// mark takes a pattern (/foo/:bar) and a string that matches that pattern
// e.g. (/foo/something) and returns a pattern (/foo/:omething).
func mark(p string, s []byte) string {
	i := 0
	j := 0

	for i < len(p) && j < len(s) {
		if p[i] == ':' {
			s[j] = ':'
		}
		for i < len(p) && p[i] != '/' {
			i++
		}
		for j < len(s) && s[j] != '/' {
			j++
		}
		i++
		j++
	}

	return string(s)
}
