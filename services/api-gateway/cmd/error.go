package main

import (
	"errors"
	"net/http"
)

var errMissingOrganizationID = errors.New("organization_id is required")
var errMissingTripID = errors.New("trip_id is required")

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.log.Errorw("internal server error",
		"method", r.Method,
		"path", r.URL.Path,
		"err", err,
	)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.log.Warnw("bad request",
		"method", r.Method,
		"path", r.URL.Path,
		"err", err,
	)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.log.Warnw("resource not found",
		"method", r.Method,
		"path", r.URL.Path,
	)
	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.log.Warnw("unauthorized request",
		"method", r.Method,
		"path", r.URL.Path,
		"err", err,
	)
	writeJSONError(w, http.StatusUnauthorized, "user unauthorized")
}

func (app *application) invalidCredentialsErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.log.Warnw("invalid credentials",
		"method", r.Method,
		"path", r.URL.Path,
		"err", err,
	)
	writeJSONError(w, http.StatusUnauthorized, "invalid credentials")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.log.Warnw("rate limit exceeded",
		"method", r.Method,
		"path", r.URL.Path,
		"retry_after", retryAfter,
	)
	w.Header().Set("Retry-After", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
