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
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/ThatSneakyCoder/RoutePulse/services/api-gateway/internal/infrastructure/grpc"
	"github.com/ThatSneakyCoder/RoutePulse/services/api-gateway/internal/infrastructure/ratelimiter"
	"github.com/ThatSneakyCoder/RoutePulse/shared/env"
)

type application struct {
	config             config
	log                *zap.SugaredLogger
	metrics            metrics
	identityClient     *grpc.IdentityServiceClient
	analyticsClient    *grpc.AnalyticsServiceClient
	organizationClient *grpc.OrganizationServiceClient
	fleetClient        *grpc.FleetServiceClient
	trackingClient     *grpc.TrackingServiceClient
	limiters           rateLimiters
}

type metrics struct {
	httpRequestsTotal *prometheus.CounterVec
}

type config struct {
	addr string
	env  string
}

type rateLimiters struct {
	global   limiterEntry
	register limiterEntry
	login    limiterEntry
}

type limiterEntry struct {
	limiter ratelimiter.Limiter
	config  rateLimiterConfig
}

type rateLimiterConfig struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}

func newLimiterEntry(cfg rateLimiterConfig) limiterEntry {
	return limiterEntry{
		config:  cfg,
		limiter: ratelimiter.NewFixedWindowLimiter(cfg.RequestsPerTimeFrame, cfg.TimeFrame),
	}
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{env.GetString("FRONTEND_URL", "http://localhost:5173")},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.rateLimitMiddleware(app.limiters.global))
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
			r.With(app.rateLimitMiddleware(app.limiters.register)).Post("/register-user", app.createUserHandler)
			r.Post("/verify-email", app.verifyUserEmailHandler)
			// implement resend-email-verification
			r.With(app.rateLimitMiddleware(app.limiters.login)).Post("/login", app.loginUserHandler)
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

				r.Get("/{orgId}", app.getOrganizationHandler)

				r.Get("/{orgId}/members", app.listOrganizationMembersHandler)
				r.Post("/{orgId}/invite", app.inviteUserToOrganizationHandler)

				r.Delete("/{orgId}/members/{userId}", app.removeMemberHandler)
				r.Put("/{orgId}/members/{userId}/role", app.updateMemberRoleHandler)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Route("/fleet", func(r chi.Router) {

				r.Route("/vehicles", func(r chi.Router) {

					r.Post("/", app.createVehicleHandler)
					r.Get("/", app.listVehiclesHandler)
					r.Get("/all", app.listAllVehiclesOfUser)
					r.Get("/{vehicleId}", app.getVehicleHandler)

					r.Put("/{vehicleId}", app.updateVehicleHandler)

					r.Patch("/{vehicleId}/status", app.updateVehicleStatusHandler)
				})

				r.Route("/drivers", func(r chi.Router) {

					r.Post("/", app.createDriverHandler)
					r.Get("/", app.listDriversHandler)

					r.Put("/{driverId}", app.updateDriverHandler)

					r.Patch("/{driverId}/status", app.updateDriverStatusHandler)
				})

				r.Route("/trips", func(r chi.Router) {

					r.Post("/", app.createTripHandler)
					r.Post("/preview-route", app.previewRouteHandler)
					r.Get("/", app.listTripsHandler)
					r.Get("/all", app.listAllTripsHandler)

					r.Post("/{tripId}/start", app.startTripHandler)
					r.Post("/{tripId}/complete", app.completeTripHandler)
				})
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Route("/tracking", func(r chi.Router) {
				r.Get("/trips/{tripId}/current", app.getTripCurrentLocationHandler)
				r.Get("/trips/{tripId}/history", app.getTripLocationHistoryHandler)
				r.Get("/trips/{tripId}/geometry", app.getTripGeometryHandler)
			})
		})

		r.Route("/ws", func(r chi.Router) {
			r.Get("/driver-tracking", app.driverTrackingWebSocketHandler)
			r.Get("/dispatch-tracking", app.dispatchTrackingWebSocketHandler)
		})

		// Analytics Service routers
		r.Group(func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)

			r.Route("/analytics", func(r chi.Router) {
				r.Get("/vehicles-in-transit", app.getVehiclesInTransit)
				r.Get("/trips-today", app.getTripsToday)

				r.Get("/total-members", app.getTotalMembers)
				r.Get("/active-users-today", app.getActiveUsersToday)
				r.Get("/recent-activity", app.getRecentActivity)
			})
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
