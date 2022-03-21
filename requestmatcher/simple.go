package requestmatcher

import (
	"github.com/ddelizia/gojitsu"
	"github.com/ddelizia/gojitsu/validation"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type simple struct {
	Methods []string
	Headers []string
	Path    *string
}

var _ gojitsu.RequestMatcher = &simple{}

func Simple(methods []string, path *string, headers []string) gojitsu.RequestMatcher {
	validation.Headers(headers)
	return &simple{
		Methods: methods,
		Headers: headers,
		Path:    path,
	}
}

func (rm *simple) Setup(m *mux.Router) *mux.Route {
	route := m.NewRoute()
	if rm.Headers != nil {
		logrus.WithField("headers", rm.Headers).Trace("Setting headers on route")
		route = route.Headers(rm.Headers...)
	}
	if rm.Path != nil {
		logrus.WithField("path", *rm.Path).Trace("Setting path on route")
		route = route.Path(*rm.Path)
	}
	if rm.Methods != nil {
		logrus.WithField("methods", rm.Methods).Trace("Setting methods on route")
		route = route.Methods(rm.Methods...)
	}

	return route
}
