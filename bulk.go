package gojitsu

import (
	"github.com/sirupsen/logrus"
	"sync"
)

// AsyncServeAll runs all servers in async way and waits for all of them to start before continuing the program.
func AsyncServeAll(servers ...IServer) {
	wg := new(sync.WaitGroup)
	wg.Add(len(servers))

	for _, server := range servers {
		go func(server IServer) {
			server.ServeAsync(wg)
		}(server)
	}

	wg.Wait()

	logrus.Debug("All servers running")
}

// CloseAll is used to cloas all the http connections in bulk.
//
// Example usage:
//	AsyncServeAll(server1, server2)
//	defer CloseAll(server1, server2)
func CloseAll(servers ...IServer) {
	for _, server := range servers {
		go func() {
			server.Close()
		}()
	}

	logrus.Debug("All servers have been shut down")
}
