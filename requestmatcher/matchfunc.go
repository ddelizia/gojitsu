package requestmatcher

import (
	"github.com/ddelizia/gojitsu"
	"github.com/gorilla/mux"
	"net/http"
)

type MatcherFunc = func(req *http.Request, match *mux.RouteMatch) bool

type matchFunc struct {
	Matcher MatcherFunc
}

var _ gojitsu.RequestMatcher = &matchFunc{}

func MatchFunc(matcherFunc MatcherFunc) gojitsu.RequestMatcher {
	return &matchFunc{
		Matcher: matcherFunc,
	}
}

func (r *matchFunc) Setup(m *mux.Router) *mux.Route {
	route := m.NewRoute()
	route.MatcherFunc(r.Matcher)
	return route
}
