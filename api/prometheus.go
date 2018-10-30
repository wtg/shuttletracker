package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "shuttletracker_in_flight_requests",
		Help: "A gauge of requests currently being served.",
	})

	// duration is partitioned by the HTTP method and handler. It uses custom
	// buckets based on the expected request duration.
	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shuttletracker_request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"code", "method"},
	)

	// responseSize has no labels, making it a zero-dimensional ObserverVec.
	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "shuttletracker_response_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 11),
		},
		[]string{"code", "method"},
	)
)

func init() {
	prometheus.MustRegister(inFlightGauge, duration, responseSize)
}

func prometheusMetrics(next http.Handler) http.Handler {

	return promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerDuration(duration,
			promhttp.InstrumentHandlerResponseSize(responseSize, next),
		),
	)
}
