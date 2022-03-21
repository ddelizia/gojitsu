package responsebuilder

import (
	"encoding/json"
	"github.com/ddelizia/gojitsu"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type echo struct {
}

func Echo() gojitsu.ResponseBuilder {
	return &echo{}
}

type EchoResponse struct {
	Body       *string
	Method     string
	Path       string
	Headers    map[string][]string
	Query      map[string][]string
	RequestUri string
	PathParams map[string]string
}

func (rb *echo) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var theBody *string
		if r.Body != nil {

			b, err := io.ReadAll(r.Body)
			if err != nil {
				logrus.Fatalln(err)
			}
			theBody = gojitsu.String(string(b))
		}

		echo := &EchoResponse{
			Headers:    r.Header,
			Body:       theBody,
			Method:     r.Method,
			Path:       r.URL.Path,
			Query:      r.URL.Query(),
			RequestUri: r.URL.RequestURI(),
			PathParams: mux.Vars(r),
		}

		jData, err := json.Marshal(echo)
		if err != nil {
			logrus.WithError(err).Panic("Not able to marshal data")
		}

		_, err = w.Write(jData)
		if err != nil {
			logrus.WithError(err).Panic("Not able to write output")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}
}
