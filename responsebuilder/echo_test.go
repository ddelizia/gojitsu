package responsebuilder_test

import (
	"bytes"
	"encoding/json"
	"github.com/ddelizia/gojitsu"
	"github.com/ddelizia/gojitsu/responsebuilder"
	"github.com/gorilla/mux"
	"github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Echo_Handle(t *testing.T) {

	echo := responsebuilder.Echo()

	tests := []struct {
		name           string
		routerPath     string
		requestMethod  string
		requestPath    string
		requestBody    *string
		requestHeaders map[string][]string
		wantMethod     string
		wantPath       string
		wantHeaders    map[string][]string
		wantQuery      map[string][]string
		wantBody       *string
		wantRequestUri string
		wantPathParams map[string]string
	}{
		{
			name:           "It should match a simple GET",
			routerPath:     "/helloGET",
			requestMethod:  http.MethodGet,
			requestPath:    "/helloGET",
			requestHeaders: map[string][]string{},
			wantMethod:     http.MethodGet,
			wantPath:       "/helloGET",
			wantHeaders:    map[string][]string{},
			wantQuery:      map[string][]string{},
			wantRequestUri: "/helloGET",
			wantPathParams: map[string]string{},
		},
		{
			name:           "It should match a simple POST",
			routerPath:     "/helloPOST",
			requestMethod:  http.MethodPost,
			requestPath:    "/helloPOST",
			requestHeaders: map[string][]string{},
			wantMethod:     http.MethodPost,
			wantPath:       "/helloPOST",
			wantHeaders:    map[string][]string{},
			wantQuery:      map[string][]string{},
			wantRequestUri: "/helloPOST",
			wantPathParams: map[string]string{},
		},
		{
			name:           "It should match a Query String",
			routerPath:     "/helloQUERY",
			requestMethod:  http.MethodPost,
			requestPath:    "/helloQUERY?hello=world",
			requestHeaders: map[string][]string{},
			wantMethod:     http.MethodPost,
			wantPath:       "/helloQUERY",
			wantHeaders:    map[string][]string{},
			wantQuery: map[string][]string{
				"hello": {"world"},
			},
			wantRequestUri: "/helloQUERY?hello=world",
			wantPathParams: map[string]string{},
		},
		{
			name:           "It should match a Multiple Query String",
			routerPath:     "/helloQUERYMULTIPLE",
			requestMethod:  http.MethodGet,
			requestPath:    "/helloQUERYMULTIPLE?hello=world&hello=world2",
			requestHeaders: map[string][]string{},
			wantMethod:     http.MethodGet,
			wantPath:       "/helloQUERYMULTIPLE",
			wantHeaders:    map[string][]string{},
			wantQuery: map[string][]string{
				"hello": {"world", "world2"},
			},
			wantRequestUri: "/helloQUERYMULTIPLE?hello=world&hello=world2",
			wantPathParams: map[string]string{},
		},
		{
			name:          "It should match multiple headers",
			routerPath:    "/helloHEADERS",
			requestMethod: http.MethodGet,
			requestPath:   "/helloHEADERS",
			requestHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantMethod: http.MethodGet,
			wantPath:   "/helloHEADERS",
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantQuery:      map[string][]string{},
			wantRequestUri: "/helloHEADERS",
			wantPathParams: map[string]string{},
		},
		{
			name:          "It should match multiple headers",
			routerPath:    "/helloHEADERSMULTIPLE",
			requestMethod: http.MethodGet,
			requestPath:   "/helloHEADERSMULTIPLE",
			requestHeaders: map[string][]string{
				"Content-Type": {"application/json", "application/xml"},
			},
			wantMethod: http.MethodGet,
			wantPath:   "/helloHEADERSMULTIPLE",
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json", "application/xml"},
			},
			wantQuery:      map[string][]string{},
			wantRequestUri: "/helloHEADERSMULTIPLE",
			wantPathParams: map[string]string{},
		},
		{
			name:           "It should match a simple POST with body",
			routerPath:     "/helloPOSTBODY",
			requestMethod:  http.MethodPost,
			requestPath:    "/helloPOSTBODY",
			requestHeaders: map[string][]string{},
			requestBody:    gojitsu.String("This is a body"),
			wantMethod:     http.MethodPost,
			wantPath:       "/helloPOSTBODY",
			wantHeaders:    map[string][]string{},
			wantQuery:      map[string][]string{},
			wantBody:       gojitsu.String("This is a body"),
			wantRequestUri: "/helloPOSTBODY",
			wantPathParams: map[string]string{},
		},
		{
			name:          "It should match a complex http query",
			routerPath:    "/helloPOST/example",
			requestMethod: http.MethodPost,
			requestPath:   "/helloPOST/example?hello=world",
			requestHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			requestBody: gojitsu.String("This is a body"),
			wantMethod:  http.MethodPost,
			wantPath:    "/helloPOST/example",
			wantHeaders: map[string][]string{
				"Content-Type": {"application/json"},
			},
			wantQuery: map[string][]string{
				"hello": {"world"},
			},
			wantBody:       gojitsu.String("This is a body"),
			wantRequestUri: "/helloPOST/example?hello=world",
			wantPathParams: map[string]string{},
		},
		{
			name:           "It should match path params",
			routerPath:     "/helloPathParams/{pathParam}",
			requestMethod:  http.MethodGet,
			requestPath:    "/helloPathParams/hello",
			requestHeaders: map[string][]string{},
			wantMethod:     http.MethodGet,
			wantPath:       "/helloPathParams/hello",
			wantHeaders:    map[string][]string{},
			wantQuery:      map[string][]string{},
			wantRequestUri: "/helloPathParams/hello",
			wantPathParams: map[string]string{
				"pathParam": "hello",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gomega.NewGomegaWithT(t)

			var body io.Reader
			if tt.requestBody != nil {
				body = bytes.NewBuffer([]byte(*tt.requestBody))
			}
			req, err := http.NewRequest(tt.requestMethod, tt.requestPath, body)
			req.Header = tt.requestHeaders
			g.Expect(err).ToNot(gomega.HaveOccurred())

			var routerPath string
			if tt.routerPath != "" {
				routerPath = tt.routerPath
			} else {
				routerPath = tt.requestPath
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.NewRoute().Path(routerPath).Handler(echo.Handle())
			router.ServeHTTP(rr, req)

			got := &responsebuilder.EchoResponse{}
			err = json.NewDecoder(rr.Body).Decode(got)
			g.Expect(err).ToNot(gomega.HaveOccurred())

			g.Expect(rr.Code).To(gomega.Equal(200))
			g.Expect(rr.Header().Get("Content-Type")).To(gomega.Equal("application/json"))

			g.Expect(got.Query).To(gomega.Equal(tt.wantQuery))
			g.Expect(got.Headers).To(gomega.Equal(tt.wantHeaders))
			g.Expect(got.Path).To(gomega.Equal(tt.wantPath))
			g.Expect(got.Method).To(gomega.Equal(tt.wantMethod))
			g.Expect(got.Body).To(gomega.Equal(tt.wantBody))
			g.Expect(got.RequestUri).To(gomega.Equal(tt.wantRequestUri))
			g.Expect(got.PathParams).To(gomega.Equal(tt.wantPathParams))
		})

	}
}
