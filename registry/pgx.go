package registry

import (
	"emailaddress.horse/thousand/repository"
	"github.com/prometheus/client_golang/prometheus"
)

type pgxCollector struct {
	repo *repository.Repository

	maxConns          *prometheus.Desc
	totalConns        *prometheus.Desc
	acquiredConns     *prometheus.Desc
	constructingConns *prometheus.Desc
	idleConns         *prometheus.Desc

	acquireCount         *prometheus.Desc
	emptyAcquireCount    *prometheus.Desc
	canceledAcquireCount *prometheus.Desc

	acquireDuration *prometheus.Desc
}

func newPgxCollector(repo *repository.Repository) prometheus.Collector {
	return &pgxCollector{
		repo: repo,
		maxConns: prometheus.NewDesc(
			"repo_pgx_max_connections",
			"Maximum number of connections available to the pool.",
			nil, prometheus.Labels{"db_name": "repository"},
		),
		totalConns: prometheus.NewDesc(
			"repo_pgx_total_connections",
			"Number of constructing, acquired and idle connections in the pool.",
			nil, prometheus.Labels{"db_name": "repository"},
		),
		acquiredConns: prometheus.NewDesc(
			"repo_pgx_acquired_connections",
			"Number of currently acquired connections in the pool.",
			nil, prometheus.Labels{"db_name": "repository"},
		),
		constructingConns: prometheus.NewDesc(
			"repo_pgx_constructing_connections",
			"Number of connections with construction in progress in the pool.",
			nil, prometheus.Labels{"db_name": "repository"},
		),
		idleConns: prometheus.NewDesc(
			"repo_pgx_idle_connections",
			"Number of currently idle connections in the pool.",
			nil, prometheus.Labels{"db_name": "repository"},
		),
		acquireCount: prometheus.NewDesc(
			"repo_pgx_acquire_count_total",
			"The count of successful acquires from the pool.",
			nil, prometheus.Labels{"db_name": "repository"},
		),
		canceledAcquireCount: prometheus.NewDesc(
			"repo_pgx_canceled_acquire_count_total",
			"The count of acquires from the pool that were canceled by a context.",
			nil, prometheus.Labels{"db_name": "repository"},
		),
		emptyAcquireCount: prometheus.NewDesc(
			"repo_pgx_empty_acquire_count_total",
			"The count of successful acquires from the pool that waited for a resource to be released or constructed because the pool was empty.",
			nil, prometheus.Labels{"db_name": "repository"},
		),
		acquireDuration: prometheus.NewDesc(
			"repo_pgx_acquire_duration_seconds_total",
			"The total duration of all successful acquires from the pool.",
			nil, prometheus.Labels{"db_name": "repository"},
		),
	}
}

func (c *pgxCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxConns
	ch <- c.totalConns
	ch <- c.acquiredConns
	ch <- c.constructingConns
	ch <- c.idleConns

	ch <- c.acquireCount
	ch <- c.emptyAcquireCount
	ch <- c.canceledAcquireCount

	ch <- c.acquireDuration
}

func (c *pgxCollector) Collect(ch chan<- prometheus.Metric) {
	s := c.repo.Stat()

	ch <- prometheus.MustNewConstMetric(c.maxConns, prometheus.GaugeValue, float64(s.MaxConns()))
	ch <- prometheus.MustNewConstMetric(c.totalConns, prometheus.GaugeValue, float64(s.TotalConns()))
	ch <- prometheus.MustNewConstMetric(c.acquiredConns, prometheus.GaugeValue, float64(s.AcquiredConns()))
	ch <- prometheus.MustNewConstMetric(c.constructingConns, prometheus.GaugeValue, float64(s.ConstructingConns()))
	ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(s.IdleConns()))

	ch <- prometheus.MustNewConstMetric(c.acquireCount, prometheus.CounterValue, float64(s.AcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.emptyAcquireCount, prometheus.CounterValue, float64(s.EmptyAcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.canceledAcquireCount, prometheus.CounterValue, float64(s.CanceledAcquireCount()))

	ch <- prometheus.MustNewConstMetric(c.acquireDuration, prometheus.CounterValue, s.AcquireDuration().Seconds())
}
