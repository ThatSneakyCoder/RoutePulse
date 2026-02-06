package main

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (app *application) handleGRPCError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	st, ok := status.FromError(err)
	if !ok {
		app.internalServerError(w, r, err)
		return
	}

	switch st.Code() {

	case codes.InvalidArgument:
		app.badRequestResponse(w, r, err)

	case codes.Unauthenticated:
		app.unauthorizedErrorResponse(w, r, err)

	case codes.PermissionDenied:
		app.unauthorizedErrorResponse(w, r, err)

	case codes.NotFound:
		app.notFoundResponse(w, r)

	case codes.AlreadyExists:
		app.badRequestResponse(w, r, err)

	case codes.ResourceExhausted:
		app.rateLimitExceededResponse(w, r, "60")

	case codes.DeadlineExceeded:
		app.internalServerError(w, r, err)

	case codes.Unavailable:
		app.internalServerError(w, r, err)

	default:
		app.internalServerError(w, r, err)
	}
}
