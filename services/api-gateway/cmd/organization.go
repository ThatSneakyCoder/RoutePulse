package main

import (
	"errors"
	"net/http"
)

// createOrganizationHandler godoc
//
//	@Summary		Create organization
//	@Description	Creates a new organization owned by the authenticated user
//	@Tags			Organizations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateOrganizationRequest	true	"Create organization payload"
//	@Success		201		{object}	OrganizationResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/organizations [post]
func (app *application) createOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromCtx(r)
	if !ok {
		app.log.Errorw("unauthorized create organization attempt")
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized"))
		return
	}

	var payload CreateOrganizationRequest
	if err := readJSON(w, r, &payload); err != nil {
		app.log.Errorw("failed to read create organization payload", "err", err)
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.log.Warnw("create organization validation failed", "err", err)
		app.badRequestResponse(w, r, err)
		return
	}

	app.log.Infow("create organization request received",
		"user_id", user.ID,
		"org_name", payload.Name,
	)

	resp, err := app.organizationClient.Client.CreateOrganization(
		r.Context(),
		payload.toProto(user.ID),
	)

	if err != nil {
		app.log.Errorw("failed to create organization",
			"user_id", user.ID,
			"org_name", payload.Name,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	app.log.Infow("organization created successfully",
		"organization_id", resp.Organization.OrganizationId,
		"user_id", user.ID,
	)

	if err := app.jsonResponse(w, http.StatusCreated, resp); err != nil {
		app.log.Errorw("failed to create organization", "err", err)
		app.internalServerError(w, r, err)
	}
}

// listUserOrganizationsHandler godoc
//
//	@Summary		List user organizations
//	@Description	Returns organizations the authenticated user belongs to
//	@Tags			Organizations
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ListUserOrganizationsResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/organizations [get]
func (app *application) listUserOrganizationsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromCtx(r)
	if !ok {
		app.log.Errorw("unauthorized list organizations attempt")
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized"))
		return
	}

	app.log.Infow("list user organizations request received",
		"user_id", user.ID,
	)

	var payload ListUserOrganizationsRequest

	resp, err := app.organizationClient.Client.ListUserOrganizations(
		r.Context(),
		payload.toProto(user.ID),
	)
	if err != nil {
		app.log.Errorw("failed to list user organizations",
			"user_id", user.ID,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	app.log.Infow("user organizations fetched successfully",
		"user_id", user.ID,
		"count", len(resp.Organizations),
	)

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.log.Errorw("failed to get organizations", "err", err)
		app.internalServerError(w, r, err)
	}
}
