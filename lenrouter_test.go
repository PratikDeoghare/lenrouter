package lenrouter_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pratikdeoghare/lenrouter"
)

func handlerAssertParams(t *testing.T, s string, expectedParams lenrouter.Params) lenrouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params lenrouter.Params) {
		patternParts := strings.Split(s, "/")
		reqPathParts := strings.Split(r.URL.Path, "/")
		//fmt.Println(strings.Join(patternParts, ">"))
		//fmt.Println(strings.Join(reqPathParts, ">"))

		if len(patternParts) != len(reqPathParts) {
			t.Errorf("pattern mismatch expected: %s, got: %s", patternParts, patternParts)
			w.WriteHeader(500)
		}

		j := 0
		for i, part := range patternParts {
			if part == "" {
				continue
			}
			if patternParts[i][0] != ':' {
				continue
			}

			param := params[j]
			key := patternParts[i][1:]
			value := reqPathParts[i]

			if key != param.Key {
				t.Errorf("key mismatch expected: %s, got: %s", key, param.Key)
				w.WriteHeader(500)
			}
			if value != param.Value {
				t.Errorf("value mismatch expected: %s, got: %s", value, param.Value)
				w.WriteHeader(500)
			}

			j++
		}

		if len(params) != len(expectedParams) {
			t.Errorf("params mismatch expected: %s, got: %s", expectedParams, params)
		}

		for i, param := range params {
			if param != expectedParams[i] {
				t.Errorf("params mismatch expected: %s, got: %s", expectedParams, params)
			}
		}

		w.WriteHeader(http.StatusOK)

	}
}

func TestRouter(t *testing.T) {

	routerTests := []struct {
		pattern        string
		reqPath        string
		expectedParams lenrouter.Params
		expectedCode   int
	}{
		{
			pattern:        "/",
			reqPath:        "/",
			expectedParams: nil,
			expectedCode:   200,
		},
		{
			pattern:        "/a",
			reqPath:        "/a",
			expectedParams: nil,
			expectedCode:   200,
		},
		{
			pattern:        "/a/b",
			reqPath:        "/a/b",
			expectedParams: nil,
			expectedCode:   200,
		},
		{
			pattern:        "/:a/b",
			reqPath:        "/a/b",
			expectedParams: lenrouter.Params{{"a", "a"}},
			expectedCode:   200,
		},
		{
			pattern:        "/:a/:b",
			reqPath:        "/a/b",
			expectedParams: lenrouter.Params{{"a", "a"}, {"b", "b"}},
			expectedCode:   200,
		},
		{
			pattern:        "/a/:boo",
			reqPath:        "/a/woogie",
			expectedParams: lenrouter.Params{{"boo", "woogie"}},
			expectedCode:   200,
		},
		{
			pattern:        "/a/:boo",
			reqPath:        "/b/woogie",
			expectedParams: nil,
			expectedCode:   404,
		},
	}

	for _, routerTest := range routerTests {
		pattern := routerTest.pattern
		router := lenrouter.New(100, 5, lenrouter.Endpoint{
			Pattern: pattern,
			Handler: handlerAssertParams(t, pattern, routerTest.expectedParams),
		})
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, routerTest.reqPath, nil)
		router.ServeHTTP(w, r)
		if w.Code != routerTest.expectedCode {
			t.Errorf("expected OK status, got: %d, details: %+v", w.Code, routerTest)
		}
	}

}

func handlerAssert(t *testing.T, s string) lenrouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params lenrouter.Params) {
		patternParts := strings.Split(s, "/")
		reqPathParts := strings.Split(r.URL.Path, "/")
		//fmt.Println(strings.Join(patternParts, ">"))
		//fmt.Println(strings.Join(reqPathParts, ">"))

		if len(patternParts) != len(reqPathParts) {
			t.Errorf("pattern mismatch expected: %s, got: %s", patternParts, patternParts)
			w.WriteHeader(500)
		}

		j := 0
		for i, part := range patternParts {
			if part == "" {
				continue
			}
			if patternParts[i][0] != ':' {
				continue
			}

			param := params[j]
			key := patternParts[i][1:]
			value := reqPathParts[i]

			if key != param.Key {
				t.Errorf("key mismatch expected: %s, got: %s", key, param.Key)
				w.WriteHeader(500)
			}
			if value != param.Value {
				t.Errorf("value mismatch expected: %s, got: %s", value, param.Value)
				w.WriteHeader(500)
			}

			j++
		}

		w.WriteHeader(http.StatusOK)

	}
}

func TestDifferentLengths(t *testing.T) {

	routerTests := []struct {
		reqPath      string
		expectedCode int
	}{
		{
			reqPath:      "/api/:foo/bar/:spam",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foo/bar/spam",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/football/bar/spamfolder",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foodisgood/bar/spamisbad",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foodnfunvilleage/bar/toomuchspam",
			expectedCode: 200,
		},
	}

	pattern := "/api/:foo/bar/:spam"
	router := lenrouter.New(100, 5, lenrouter.Endpoint{
		Pattern: pattern,
		Handler: handlerAssert(t, pattern),
	})

	for _, routerTest := range routerTests {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, routerTest.reqPath, nil)
		router.ServeHTTP(w, r)
		if w.Code != routerTest.expectedCode {
			t.Errorf("expected OK status, got: %d, details: %+v", w.Code, routerTest)
		}
	}
	//lenrouter.Print(router)
}

func TestDifferentLengthsTwice(t *testing.T) {

	routerTests := []struct {
		reqPath      string
		expectedCode int
	}{
		{
			reqPath:      "/api/:foo/bar/:spam",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foo/bar/spam",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/football/bar/spamfolder",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foodisgood/bar/spamisbad",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foodnfunvilleage/bar/toomuchspam",
			expectedCode: 200,
		},
	}

	pattern := "/api/:foo/bar/:spam"
	router := lenrouter.New(100, 5, lenrouter.Endpoint{
		Pattern: pattern,
		Handler: handlerAssert(t, pattern),
	})

	for _, routerTest := range routerTests {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, routerTest.reqPath, nil)
		router.ServeHTTP(w, r)
		if w.Code != routerTest.expectedCode {
			t.Errorf("expected OK status, got: %d, details: %+v", w.Code, routerTest)
		}
	}
	for _, routerTest := range routerTests {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, routerTest.reqPath, nil)
		router.ServeHTTP(w, r)
		if w.Code != routerTest.expectedCode {
			t.Errorf("expected OK status, got: %d, details: %+v", w.Code, routerTest)
		}
	}
	//lenrouter.Print(router)
}

func TestTwoPatternsTwice(t *testing.T) {

	routerTests := []struct {
		reqPath      string
		expectedCode int
	}{
		{
			reqPath:      "/api/:foo/bar/:spam",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foo/bar/spam",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/football/bar/spamfolder",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foodisgood/bar/spamisbad",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foodnfunvilleage/bar/toomuchspam",
			expectedCode: 200,
		},

		{
			reqPath:      "/api/:foo/car/:spam",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foo/car/spam",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/football/car/spamfolder",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foodisgood/car/spamisbad",
			expectedCode: 200,
		},
		{
			reqPath:      "/api/foodnfunvilleage/car/toomuchspam",
			expectedCode: 200,
		},
	}

	router := lenrouter.New(100, 5,
		lenrouter.Endpoint{
			Pattern: "/api/:foo/bar/:spam",
			Handler: handlerAssert(t, "/api/:foo/bar/:spam"),
		},
		lenrouter.Endpoint{
			Pattern: "/api/:foo/car/:spam",
			Handler: handlerAssert(t, "/api/:foo/car/:spam"),
		},
	)

	for _, routerTest := range routerTests {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, routerTest.reqPath, nil)
		router.ServeHTTP(w, r)
		if w.Code != routerTest.expectedCode {
			t.Errorf("expected OK status, got: %d, details: %+v", w.Code, routerTest)
		}
	}
	for _, routerTest := range routerTests {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, routerTest.reqPath, nil)
		router.ServeHTTP(w, r)
		if w.Code != routerTest.expectedCode {
			t.Errorf("expected OK status, got: %d, details: %+v", w.Code, routerTest)
		}
	}
	//lenrouter.Print(router)
}
