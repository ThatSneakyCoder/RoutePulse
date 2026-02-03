package main

import (
	"net/http"
)

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
//	@Router			/authentication/registerUser [post]
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody CreateUserRequest

	if err := readJSON(w, r, &reqBody); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(reqBody); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.identityClient.Client.RegisterUser(
		r.Context(),
		reqBody.toProto(),
	)
	if err != nil {
		app.log.Errorw("identity service RegisterUser failed", "err", err)
		app.internalServerError(w, r, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// func (app *application) createJwtTokenHandler(w http.ResponseWriter, r *http.Request) {
// 	var payload CreateUserTokenPayload
// 	if err := readJSON(w, r, &payload); err != nil {
// 		app.log.Errorw("failed to read login JSON", "error", err)
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := Validate.Struct(payload); err != nil {
// 		app.log.Errorw("login payload validation failed", "error", err)
// 		app.badRequestResponse(w, r, err)
// 		return
// 	}

// 	app.log.Infow("login request received", "email", payload.Email)

// 	resp, err := app.identityClient.Client.Login(
// 		r.Context(),
// 		payload.toProto(),
// 	)
// 	if err != nil {
// 		app.log.Errorw("identity service login failed",
// 			"email", payload.Email,
// 			"err", err,
// 		)

// 		// Map identity-service errors → HTTP errors
// 		switch {
// 		case isUnauthorized(err):
// 			app.unauthorizedErrorResponse(w, r, err)
// 		default:
// 			app.internalServerError(w, r, err)
// 		}
// 		return
// 	}

// 	app.log.Infow("jwt issued successfully",
// 		"user_email", payload.Email,
// 	)

// 	if err := app.jsonResponse(w, http.StatusCreated, resp); err != nil {
// 		app.log.Errorw("failed to write login response", "err", err)
// 		app.internalServerError(w, r, err)
// 	}
// }
