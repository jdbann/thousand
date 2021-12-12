package registry

import (
	"emailaddress.horse/thousand/repository"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func() *prometheus.Registry {
		r := prometheus.NewRegistry()

		r.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		r.MustRegister(collectors.NewGoCollector())

		return r
	}),
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

	return r, nil
}
