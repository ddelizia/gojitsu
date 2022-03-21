package gojitsu_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/ddelizia/gojitsu"
	"github.com/ddelizia/gojitsu/requestmatcher"
	"github.com/ddelizia/gojitsu/responsebuilder"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestServer(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	logrus.SetLevel(logrus.TraceLevel)

	currentPort := 9000

	c := http.DefaultClient

	buildConfig := func() *gojitsu.ServerConfig {
		currentPort++
		return &gojitsu.ServerConfig{
			Port:              currentPort,
			WriteTimeoutMills: 10000,
			ReadTimeoutMills:  10000,
			Host:              "127.0.0.1",
		}
	}

	tests := []struct {
		name   string
		config *gojitsu.ServerConfig
		route  *gojitsu.Route

		requestHeaders http.Header
		requestPath    string

		expectMethod          []string
		expectBody            string
		expectStatus          int
		expectResponseHeaders http.Header

		hasNoMatch bool
	}{
		{
			name:   "It should GET POST PATCH PUT a simple body",
			config: buildConfig(),
			route: &gojitsu.Route{
				Id: "Id",
				RequestMatcher: requestmatcher.Simple(
					[]string{
						http.MethodGet,
						http.MethodPost,
						http.MethodPatch,
						http.MethodPut,
					},
					gojitsu.String("/path/{code}"),
					[]string{"reqkey", "reqvalue0", "reqkey", "reqvalue1", "reqhello", "reqworld"},
				),
				ResponseBuilder: responsebuilder.String(200, gojitsu.String("This is body"), []string{"key", "value0", "key", "value1", "hello", "world"}),
			},

			requestPath: "/path/somepath",
			requestHeaders: map[string][]string{
				"reqkey":   {"reqvalue0", "reqvalue1"},
				"reqhello": {"reqworld"},
			},

			expectMethod: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut},
			expectBody:   "This is body",
			expectResponseHeaders: map[string][]string{
				"key":   {"value0", "value1"},
				"hello": {"world"},
			},
			expectStatus: 200,
		},
		{
			name:   "It should GET POST PATCH PUT a echo body",
			config: buildConfig(),
			route: &gojitsu.Route{
				Id: "Id",
				RequestMatcher: requestmatcher.Simple(
					[]string{
						http.MethodGet,
					},
					gojitsu.String("/path/{code}"),
					[]string{"reqkey", "reqvalue0", "reqkey", "reqvalue1", "reqhello", "reqworld"},
				),
				ResponseBuilder: responsebuilder.Echo(),
			},

			requestPath: "/path/somepath?string=data",
			requestHeaders: map[string][]string{
				"reqkey":   {"reqvalue0", "reqvalue1"},
				"reqhello": {"reqworld"},
			},

			expectMethod:          []string{http.MethodGet},
			expectBody:            "{\"Body\":\"\",\"Method\":\"GET\",\"Path\":\"/path/somepath\",\"Headers\":{\"Accept-Encoding\":[\"gzip\"],\"Reqhello\":[\"reqworld\"],\"Reqkey\":[\"reqvalue1\"],\"User-Agent\":[\"Go-http-client/1.1\"]},\"Query\":{\"string\":[\"data\"]},\"RequestUri\":\"/path/somepath?string=data\",\"PathParams\":{\"code\":\"somepath\"}}",
			expectResponseHeaders: map[string][]string{},
			expectStatus:          200,
		},
	}

	for _, tt := range tests {
		func() {
			srv := gojitsu.Server(tt.config, tt.route)
			gojitsu.AsyncServeAll(srv)
			defer srv.Close()
			for _, method := range tt.expectMethod {
				t.Run(fmt.Sprintf("%s [%s]", tt.name, method), func(t *testing.T) {
					// Building request
					req, err := http.NewRequest(method, fmt.Sprint("http://", tt.config.Host, ":", tt.config.Port, tt.requestPath), nil)

					if tt.requestHeaders != nil {
						for key, values := range tt.requestHeaders {
							for _, v := range values {
								req.Header.Set(key, v)
							}
						}
					}

					g.Expect(err).To(gomega.BeNil(), "Error creating the request")

					got, err := c.Do(req)
					g.Expect(err).To(gomega.BeNil(), "Error executing the request")

					if tt.hasNoMatch {
						g.Expect(got.StatusCode).To(gomega.Equal(404))
						return
					}

					bodyBytes, err := io.ReadAll(got.Body)
					defer got.Body.Close()
					bodyString := string(bodyBytes)

					g.Expect(bodyString).To(gomega.Equal(tt.expectBody))
					g.Expect(got.StatusCode).To(gomega.Equal(tt.expectStatus))
					if tt.expectResponseHeaders != nil {
						for key, value := range tt.expectResponseHeaders {
							g.Expect(got.Header.Values(key)).To(gomega.Equal(value))
						}
					}
				})
			}
		}()
	}
}
