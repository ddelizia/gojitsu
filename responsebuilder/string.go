package responsebuilder

import (
	"github.com/ddelizia/gojitsu"
	"github.com/ddelizia/gojitsu/validation"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type responseBuilderString struct {
	Body    *string
	Status  int
	Headers []string
}

func String(status int, body *string, headers []string) gojitsu.ResponseBuilder {
	validation.Headers(headers)
	return &responseBuilderString{
		Body:    body,
		Status:  status,
		Headers: headers,
	}
}

func (rb *responseBuilderString) Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithField("request-call", r).Trace("Request details")

		logrus.WithField("request-path-params", mux.Vars(r)).Trace("Path params")

		for i := 0; i < len(rb.Headers)/2; i++ {
			key := rb.Headers[i*2]
			value := rb.Headers[i*2+1]
			logrus.WithFields(map[string]interface{}{
				"key":   key,
				"value": value,
			}).Trace("Setting response header")
			w.Header().Add(key, value)
		}

		w.WriteHeader(rb.Status)
		if rb.Body != nil {
			logrus.WithField("response-body", *rb.Body).Trace("Setting response body on route")
			_, err := w.Write([]byte(*rb.Body))
			if err != nil {
				logrus.WithError(err).Panic("Not able to write body")
				return
			}
		}
	}
}
