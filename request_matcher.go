package gojitsu

import (
	"github.com/gorilla/mux"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RequestMatcher

type RequestMatcher interface {
	Setup(m *mux.Router) *mux.Route
}
