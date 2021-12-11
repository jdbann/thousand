package browser

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-chi/chi/v5"
)

type integrationTS struct {
	mux    *chi.Mux
	server *httptest.Server
}

func newIntegrationTS() *integrationTS {
	mux := chi.NewMux()
	return &integrationTS{
		mux:    mux,
		server: httptest.NewUnstartedServer(mux),
	}
}

func (ts *integrationTS) URL() string {
	wait := time.NewTimer(500 * time.Millisecond)
	for {
		select {
		case <-wait.C:
			panic("error getting address for integration test server")
		default:
			if ts.server.URL != "" {
				return ts.server.URL
			}
		}
	}
}

func (ts *integrationTS) ListenAndServe() error {
	ts.server.Start()
	return http.ErrServerClosed
}

func (ts *integrationTS) Shutdown(ctx context.Context) error {
	ts.server.Close()
	return nil
}
