package main

import (
	"fmt"
	"net/http"
)

type userKey string

const userCtx userKey = "user"

// createUserHandler godoc
//
//	@Summary		Create user
//	@Description	Creates a new user by delegating to the identity service
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateUserRequest	true	"Create user payload"
//	@Success		201		{object}	CreateUserResponse
//	@Failure		400		{object}	map[string]string	"Validation error"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/authentication/register-user [post]
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	app.log.Infow("create user request received",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	var reqBody CreateUserRequest

	if err := readJSON(w, r, &reqBody); err != nil {
		app.log.Warnw("failed to read create user request body",
			"err", err,
		)
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(reqBody); err != nil {
		app.log.Warnw("create user request validation failed",
			"email", reqBody.Email,
			"err", err,
		)
		app.badRequestResponse(w, r, err)
		return
	}

	app.log.Infow("calling identity service RegisterUser",
		"email", reqBody.Email,
	)

	resp, err := app.identityClient.Client.RegisterUser(
		r.Context(),
		reqBody.toProto(),
	)
	if err != nil {
		app.log.Errorw("identity service RegisterUser failed",
			"email", reqBody.Email,
			"err", err,
		)
		app.internalServerError(w, r, err)
		return
	}

	app.log.Infow("user created successfully",
		"user_id", resp.User.GetUserId(),
		"email", reqBody.Email,
	)

	writeJSON(w, http.StatusCreated, resp)
}

// loginUserHandler godoc
//
//	@Summary		Login user
//	@Description	Authenticates a user and returns a JWT access token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginUserRequest	true	"Login credentials"
//	@Success		201		{object}	LoginUserResponse
//	@Failure		400		{object}	map[string]string	"Validation error"
//	@Failure		401		{object}	map[string]string	"Invalid credentials"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/authentication/login [post]
func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginUserRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.log.Errorw("failed to read login JSON", "error", err)
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.log.Errorw("login payload validation failed", "error", err)
		app.badRequestResponse(w, r, err)
		return
	}

	app.log.Infow("login request received", "email", payload.Email)

	resp, err := app.identityClient.Client.Login(
		r.Context(),
		payload.toProto(),
	)
	if err != nil {
		app.log.Errorw("identity service login failed",
			"email", payload.Email,
			"err", err,
		)

		app.handleGRPCError(w, r, err)
		return
	}

	app.log.Infow("jwt issued successfully",
		"user_email", payload.Email,
	)

	if err := app.jsonResponse(w, http.StatusCreated, resp); err != nil {
		app.log.Errorw("failed to write login response", "err", err)
		app.internalServerError(w, r, err)
	}
}

// getUserHandler godoc
//
//	@Summary		Get current user
//	@Description	Returns the authenticated user's profile
//	@Tags			Users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ValidateTokenResponse
//	@Failure		401	{object}	map[string]string	"Unauthorized"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/user/me [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromCtx(r)
	if !ok {
		app.log.Errorw("user missing from context")
		app.unauthorizedErrorResponse(w, r, fmt.Errorf("unauthorized"))
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getUserFromCtx(r *http.Request) (*AuthenticatedUser, bool) {
	user, ok := r.Context().Value(userCtx).(*AuthenticatedUser)
	return user, ok
}

// verifyUserEmailHandler godoc
//
//	@Summary		Verify user email
//	@Description	Verifies a user's email address using a 6-digit verification token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		VerifyUserEmailRequest	true	"Email verification payload"
//	@Success		200		{object}	VerifyUserEmailResponse
//	@Failure		400		{object}	map[string]string	"Validation error"
//	@Failure		401		{object}	map[string]string	"Invalid or expired verification token"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/authentication/verify-email [post]
func (app *application) verifyUserEmailHandler(w http.ResponseWriter, r *http.Request) {
	var payload VerifyUserEmailRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.log.Errorw("failed to read verifyEmailToken JSON", "error", err)
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.log.Errorw("verifyEmailToken payload validation failed", "error", err)
		app.badRequestResponse(w, r, err)
		return
	}

	app.log.Infow("verifyEmailToken request received", "email", payload.Email, "token",payload.Token)

	resp, err := app.identityClient.Client.VerifyUserEmail(r.Context(), payload.toProto())
	if err != nil {
		app.log.Errorw("identity service verify email token",
			"email", payload.Email,
			"err", err,
		)

		app.handleGRPCError(w, r, err)
		return
	}

	app.log.Infow("user verified successfully",
		"userID", resp.UserId,
		"isVerified", resp.IsVerified,
	)

	writeJSON(w, http.StatusOK, resp)
}
