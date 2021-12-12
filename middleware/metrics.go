package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func CollectMetrics(r chi.Router, registry *prometheus.Registry) {
	requests := promauto.With(registry).NewCounterVec(prometheus.CounterOpts{
		Name: "app_http_requests_total",
		Help: "The total number of HTTP requests.",
	}, []string{"method", "path", "code"})

	requestDuration := promauto.With(registry).NewHistogramVec(prometheus.HistogramOpts{
		Name:    "app_http_requests_duration_seconds",
		Help:    "HTTP request durations.",
		Buckets: []float64{.005, .01, .05, .1, .5, 1},
	}, []string{"code"})

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, 1)
			before := time.Now()
			next.ServeHTTP(ww, r)
			duration := time.Since(before)
			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}
			code := strconv.Itoa(status)
			requests.WithLabelValues(r.Method, r.URL.Path, code).Inc()
			requestDuration.WithLabelValues(code).Observe(duration.Seconds())
		})
	})
}
