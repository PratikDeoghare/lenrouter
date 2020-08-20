package lenrouter

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO(pratik): implement
// lenrouter ignores http method, worries only about paths
// implement router that takes method into account on top of lenrouter
// maybe just one lenrouter per http method type
func NewMethodRouter() {

}

// stolen from httprouter https://github.com/julienschmidt/httprouter
// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

type Handle func(http.ResponseWriter, *http.Request, Params)

// Print dumps the state of the router. Useful for debugging.
func Print(h http.Handler) {
	type Guess struct {
		EndpointIdx int    `json:"endpoint_idx"`
		CheckSet    []int  `json:"check_set"`
		Instance    string `json:"instance"`
	}

	r := h.(*router)
	for l, guesses := range r.lenToGuesses {
		if guesses == nil {
			continue
		}
		var gs []Guess
		for _, guess := range guesses {
			gs = append(gs, Guess{
				EndpointIdx: guess.endpointIdx,
				CheckSet:    guess.checkSet,
				Instance:    guess.instance,
			})
		}

		fmt.Println(l)
		data, _ := json.MarshalIndent(gs, "", " ")
		fmt.Println(string(data))
	}
}
