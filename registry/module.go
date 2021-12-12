package registry

import (
	"context"
	"net/http"

	"emailaddress.horse/thousand/repository"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(fxNew),
)

type Params struct {
	fx.In

	Repository *repository.Repository
}

func fxNew(lc fx.Lifecycle, params Params) (*prometheus.Registry, error) {
	r := prometheus.NewRegistry()

	if err := r.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
		return nil, err
	}

	if err := r.Register(collectors.NewGoCollector()); err != nil {
		return nil, err
	}

	if err := r.Register(newPgxCollector(params.Repository)); err != nil {
		return nil, err
	}

	s := &http.Server{
		Addr:    ":9091",
		Handler: promhttp.HandlerFor(r, promhttp.HandlerOpts{}),
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() { _ = s.ListenAndServe() }()
			return nil
		},
		OnStop: s.Shutdown,
	})

	return r, nil
}
