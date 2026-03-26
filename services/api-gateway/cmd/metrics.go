package main

import (
	"bufio"
	"net"
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

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}

	return r.ResponseWriter.Write(data)
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}

	return hijacker.Hijack()
}
