//+build integration

package e2e

import (
	"fmt"
	"net/http"

	"github.com/fortytw2/hydrocarbon"
	"github.com/fortytw2/hydrocarbon/pg"
)

type testServer struct {
	addr string
	s    http.Server
	db   *pg.DB
	mm   *hydrocarbon.MockMailer
}

type testServerPool struct {
	servers chan *testServer
}

func newTestServerPool(instances int) *testServerPool {
	var addrs []string
	for index := 0; index < instances; index++ {
		addrs = append(addrs, fmt.Sprintf(":690%d", index))
	}

	var servers []testServer
	for _, addr := range addrs {
		fullAddr := fmt.Sprintf("http://localhost%s", addr)

		server, db, mockMailer, cancel := SetupE2EServer(t, fullAddr)
		defer cancel()

		servers = append(servers, testServer{
			addr: fullAddr,
			s:    server,
			db:   db,
			mm:   mockMailer,
		})

		go http.ListenAndServe(addr, server)
	}

	serverChan := make(chan *testServer, len(servers))
	for _, s := range servers {
		serverChan <- s
	}

	return &testServerPool{
		servers: serverChan,
	}
}

func (t *testServerPool) getServer() (*testServer, func()) {
	s := <-t.servers
	return s, func() {
		t.servers <- s
	}
}
