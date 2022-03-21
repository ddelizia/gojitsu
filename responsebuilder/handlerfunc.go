package responsebuilder

import (
	"github.com/ddelizia/gojitsu"
	"net/http"
)

type handlerFunc struct {
	HandlerFunc http.HandlerFunc
}

func HandlerFunc(h http.HandlerFunc) gojitsu.ResponseBuilder {
	return &handlerFunc{
		HandlerFunc: h,
	}
}

func (rb *handlerFunc) Handle() http.HandlerFunc {
	return rb.HandlerFunc
}
