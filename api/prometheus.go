package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "shuttletracker_http_requests_count",
		}, []string{"code", "method", "path"},
	)
	latency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "shuttletracker_http_latency_seconds",
		}, []string{"code", "method", "path"},
	)
)

func prometheusMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(wrapper, r)
		requests.WithLabelValues(strconv.Itoa(wrapper.Status()), r.Method, r.URL.Path).Inc()
		latency.WithLabelValues(strconv.Itoa(wrapper.Status()), r.Method, r.URL.Path).Observe(float64(time.Since(start).Nanoseconds()) / 1000000000)
	})
}
