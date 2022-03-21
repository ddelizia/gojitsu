package responsebuilder

import (
	"github.com/ddelizia/gojitsu"
	"net/http"
)

type handler struct {
	Handler http.Handler
}

func Handler(h http.Handler) gojitsu.ResponseBuilder {
	return &handler{
		Handler: h,
	}
}

func (rb *handler) Handle() http.HandlerFunc {
	return rb.Handler.ServeHTTP
}
