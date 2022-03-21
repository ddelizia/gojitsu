package gojitsu

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"sync"
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . IServer

// IServer interface provides methods needed to manage a mock http server
type IServer interface {
	// Serve starts up the server in sync way
	Serve()

	// ServeAsync starts up the server in an async way, when ready if notidy WaitGroup to Done
	ServeAsync(wg *sync.WaitGroup)

	// Close shuts down the server. Make sure to call this method before existing the method
	Close()

	// GetPort returns the port on which the server is running.
	GetPort() int
}

// ServerConfig is used to configure the http mock server.
type ServerConfig struct {
	// Port on which server is running.
	Port int

	// WriteTimeoutMills is the time in milliseconds after which the http write operation times out.
	WriteTimeoutMills time.Duration `yaml:"write_timeout_mills"`

	// WriteTimeoutMills is the time in milliseconds after which the http read operation times out.
	ReadTimeoutMills time.Duration `yaml:"read_timeout_mills"`

	// Host on which the server listen for requests.
	// Use 0.0.0.0 or (empty) for accepting calls from any source.
	Host string
}

type server struct {
	Routes     map[string]*Route
	MuxRoutes  map[string]*mux.Route
	Config     *ServerConfig
	Router     *mux.Router
	HttpServer http.Server
}

func Server(config *ServerConfig, routes ...*Route) IServer {
	router := mux.NewRouter()
	router.Use(timerMiddleware)
	router.Use(dataLoggerMiddleware)
	routesMap := map[string]*Route{}
	muxRoutesMap := map[string]*mux.Route{}
	for _, r := range routes {
		muxRoute := r.RequestMatcher.Setup(router)
		muxRoute.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			r.ResponseBuilder.Handle()(writer, request)
		})
		routesMap[r.Id] = r
		muxRoutesMap[r.Id] = muxRoute
	}
	return server{
		Routes:    routesMap,
		MuxRoutes: muxRoutesMap,
		Config:    config,
		Router:    router,
	}
}

func (g server) GetPort() int {
	return g.Config.Port
}

func (g server) Serve() {
	g.ServeAsync(nil)
}

func (g server) ServeAsync(wg *sync.WaitGroup) {
	port := g.GetPort()
	srv := http.Server{
		Handler:      g.Router,
		WriteTimeout: g.Config.WriteTimeoutMills * time.Millisecond,
		ReadTimeout:  g.Config.ReadTimeoutMills * time.Millisecond,
	}

	g.HttpServer = srv
	l, err := net.Listen("tcp", fmt.Sprint(g.Config.Host, ":", port))
	if err != nil {
		logrus.WithError(err).Panic("Not able to create tcp")
	}

	if wg != nil {
		wg.Done()
	}
	logrus.WithField("port", port).Info("Server starting")

	logrus.Fatal(g.HttpServer.Serve(l))
}

func (g server) Close() {
	logrus.WithField("port", g.Config.Port).Debug("Shutdown server")
	err := g.HttpServer.Shutdown(context.Background())
	if err != nil {

		logrus.WithField("Id", "Id").Panic("Not able to shutdown server")
	}
	err = g.HttpServer.Close()
	if err != nil {
		logrus.WithField("Id", "Id").Panic("Not able to close server")
	}
}
