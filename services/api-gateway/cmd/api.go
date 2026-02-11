package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	docs "github.com/ThatSneakyCoder/RoutePulse/services/api-gateway/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/ThatSneakyCoder/RoutePulse/services/api-gateway/internal/infrastructure/grpc"
)

type application struct {
	config             config
	log                *zap.SugaredLogger
	metrics            metrics
	identityClient     *grpc.IdentityServiceClient
	analyticsClient    *grpc.AnalyticsServiceClient
	organizationClient *grpc.OrganizationServiceClient
}

type metrics struct {
	httpRequestsTotal *prometheus.CounterVec
}

type config struct {
	addr string
	env  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.metrics.prometheusMiddleware)

	r.Handle("/metrics", promhttp.Handler())

	r.Route("/v1", func(r chi.Router) {
		if app.config.env == "development" {
			r.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
				docs.SwaggerInfo.Host = r.Host
				docs.SwaggerInfo.Schemes = []string{"http"}
				docs.SwaggerInfo.BasePath = "/v1"

				httpSwagger.Handler(
					httpSwagger.URL("/v1/swagger/doc.json"),
				)(w, r)
			})
		}
		r.Get("/health", app.healthCheckHandler)

		// identity service routers
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/register-user", app.createUserHandler)
			r.Post("/verify-email", app.verifyUserEmailHandler)
			// implement resend-email-verification
			r.Post("/login", app.loginUserHandler)
			r.Post("/forgot-password", app.forgotPasswordHandler)
			r.Put("/reset-password", app.resetPasswordHandler)
			// implement logout
			// implement de-activate account
			// implement change password
			// implement login attempt tracking
			// implement MFA (TOPT)
			// implement update profile
			// implement change email endpoint
			// implement way to add avatar
			// implement openfga and fine grained access
		})

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Route("/user/me", func(r chi.Router) {
				r.Get("/", app.getUserHandler)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Route("/organizations", func(r chi.Router) {
				r.Post("/", app.createOrganizationHandler)
				r.Get("/", app.listUserOrganizationsHandler)
			})
		})

		// Analytics Service routers
		r.Route("/analytics", func(r chi.Router) {
			r.Get("/vehicles-in-transit", app.getVehiclesInTransit)
			r.Get("/trips-today", app.getTripsToday)
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = "0.0.1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		app.log.Infow("shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdown <- srv.Shutdown(ctx)
	}()

	app.log.Infow("http server started", "addr", app.config.addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdown; err != nil {
		return err
	}

	app.log.Info("http server shut down gracefully")
	return nil
}
