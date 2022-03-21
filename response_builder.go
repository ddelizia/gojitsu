package gojitsu

import (
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ResponseBuilder

type ResponseBuilder interface {
	Handle() http.HandlerFunc
}
