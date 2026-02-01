package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func newMetrics() metrics {
	// 1. Define the metrics
	m := metrics{
		httpRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests processed, labeled by method, route, and status.",
			},
			[]string{"method", "status", "route"},
		),
	}

	// 2. Register the metrics
	prometheus.MustRegister(m.httpRequestsTotal)

	return m
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
