package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func (app *application) rateLimitMiddleware(entry limiterEntry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !entry.config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// RealIP middleware already runs before this
			ip := r.RemoteAddr

			allowed, retryAfter := entry.limiter.Allow(ip)
			if !allowed {
				w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.log.Warnw("missing authorization header",
				"path", r.URL.Path,
				"method", r.Method,
			)
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("you are not authorized"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.log.Warnw("invalid authorization header format",
				"path", r.URL.Path,
				"method", r.Method,
			)
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("you aren't authorized"))
			return
		}

		token := parts[1]

		ctx := r.Context()
		payload := &ValidateTokenRequest{
			AccessToken: token,
		}

		user, err := app.identityClient.Client.ValidateToken(ctx, payload.toProto())
		if err != nil {
			app.log.Errorw("token validation failed",
				"path", r.URL.Path,
				"method", r.Method,
				"err", err,
			)
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("you aren't authorized"))
			return
		}

		app.log.Debugw("request authenticated",
			"user_id", user.User.UserId,
			"user_id_verified", user.User.IsVerified,
			"path", r.URL.Path,
			"method", r.Method,
		)

		authenticatedUser := &AuthenticatedUser{
			ID:         user.User.GetUserId(),
			Email:      user.User.GetEmail(),
			FirstName:  user.User.GetFirstname(),
			LastName:   user.User.GetLastname(),
			IsActive:   user.User.GetIsActive(),
			IsVerified: user.User.GetIsVerified(),
			CreatedAt:  user.User.GetCreatedAt().AsTime(),
			UpdatedAt:  user.User.GetUpdatedAt().AsTime(),
		}

		ctx = context.WithValue(ctx, userCtx, authenticatedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
