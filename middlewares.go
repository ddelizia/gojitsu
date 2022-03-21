package gojitsu

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func timerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		next.ServeHTTP(w, r)

		end := time.Now()

		duration := end.Sub(start)

		logrus.WithField("duration", duration).Debug("Request duration in ns")
	})
}

func dataLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(map[string]interface{}{
			"request": r,
		}).Debug("Request info")

		next.ServeHTTP(w, r)
	})
}
