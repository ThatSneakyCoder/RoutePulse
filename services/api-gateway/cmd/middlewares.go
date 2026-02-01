package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (m *metrics) prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we need encapsulation
		recorded := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// propagate the request
		next.ServeHTTP(recorded, r)

		// fetch basic request details
		route := chi.RouteContext(r.Context()).RoutePattern()
		method := r.Method
		status := strconv.Itoa(recorded.status)

		m.httpRequestsTotal.WithLabelValues(method, status, route).Inc()
	})
}
