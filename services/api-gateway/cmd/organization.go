package main

import (
	"errors"
	"net/http"

	"github.com/ThatSneakyCoder/RoutePulse/shared/proto/identity"
	"github.com/go-chi/chi/v5"
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

// getOrganizationHandler godoc
//
//	@Summary		Get organization members
//	@Description	Returns members of an organization enriched with user profile data
//	@Tags			Organizations
//	@Produce		json
//	@Security		BearerAuth
//	@Param			orgId	path		string	true	"Organization ID"
//	@Success		200		{object}	GetOrganizationMembersListResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/organizations/{orgId}/members [get]
func (app *application) listOrganizationMembersHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := getUserFromCtx(r)
	if !ok {
		app.log.Errorw("unauthorized get organization members attempt")
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized"))
		return
	}

	orgID := chi.URLParam(r, "orgId")
	if orgID == "" {
		app.badRequestResponse(w, r, errors.New("organization id required"))
		return
	}

	app.log.Infow("get organization members request received",
		"user_id", user.ID,
		"organization_id", orgID,
	)

	// 1. fetch membership from organization service
	orgReq := GetOrganizationMembersListRequest{
		OrganizationID: orgID,
	}

	orgResp, err := app.organizationClient.Client.ListOrganizationMembers(
		r.Context(),
		orgReq.toProto(),
	)

	if err != nil {
		app.log.Errorw("failed to fetch organization members",
			"user_id", user.ID,
			"organization_id", orgID,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	// 2. collect user IDs
	userIDs := make([]string, 0, len(orgResp.Members))
	for _, m := range orgResp.Members {
		userIDs = append(userIDs, m.UserId)
	}

	// 3. fetch user profiles from identity service
	idsReq := UserIDsRequest{
		UserIDs: userIDs,
	}

	usersResp, err := app.identityClient.Client.GetUsersByIDs(
		r.Context(),
		idsReq.toProto(),
	)

	if err != nil {
		app.log.Errorw("failed to fetch users from identity service",
			"organization_id", orgID,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	// 4. build lookup map
	userMap := make(map[string]*identity.UserSummary)
	for _, u := range usersResp.Users {
		userMap[u.Id] = u
	}

	// 5. merge membership + identity data
	members := make([]OrganizationMemberResponse, 0, len(orgResp.Members))

	for _, m := range orgResp.Members {

		member := OrganizationMemberResponse{
			UserID:   m.UserId,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		}

		if u, ok := userMap[m.UserId]; ok {
			member.FirstName = u.FirstName
			member.LastName = u.LastName
			member.Email = u.Email
		}

		members = append(members, member)
	}

	resp := GetOrganizationMembersListResponse{
		Members: members,
	}

	app.log.Infow("organization members returned successfully",
		"organization_id", orgID,
		"count", len(resp.Members),
	)

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.log.Errorw("failed to write response",
			"err", err,
		)
		app.internalServerError(w, r, err)
	}
}

// getOrganizationHandler godoc
//
//	@Summary		Get organization
//	@Description	Get details of a specific organization by ID
//	@Tags			Organizations
//	@Produce		json
//	@Security		BearerAuth
//	@Param			orgId	path		string	true	"Organization ID"
//	@Success		200		{object}	OrganizationResponse
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized"
//	@Failure		404		{object}	map[string]string	"Organization not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/organizations/{orgId} [get]
func (app *application) getOrganizationHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := getUserFromCtx(r)
	if !ok {
		app.log.Errorw("unauthorized get organization attempt")
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized"))
		return
	}

	orgID := chi.URLParam(r, "orgId")
	if orgID == "" {
		app.badRequestResponse(w, r, errors.New("organization id required"))
		return
	}

	app.log.Infow("get organization request received",
		"user_id", user.ID,
		"organization_id", orgID,
	)

	req := GetOrganizationRequest{
		OrganizationID: orgID,
	}

	resp, err := app.organizationClient.Client.GetOrganization(
		r.Context(),
		req.toProto(),
	)

	if err != nil {
		app.log.Errorw("failed to fetch organization",
			"user_id", user.ID,
			"organization_id", orgID,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	result := OrganizationResponse{
		OrganizationID: resp.Organization.OrganizationId,
		Name:           resp.Organization.Name,
		OwnerUserID:    resp.Organization.OwnerUserId,
		IsActive:       resp.Organization.IsActive,
		CreatedAt:      resp.Organization.CreatedAt,
		UpdatedAt:      resp.Organization.UpdatedAt,
	}

	app.log.Infow("organization fetched successfully",
		"organization_id", orgID,
	)

	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {
		app.internalServerError(w, r, err)
	}
}

// inviteUserToOrganizationHandler godoc
//
//	@Summary		Invite user to organization
//	@Description	Adds an existing user to an organization with a role
//	@Tags			Organizations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			orgId	path		string	true	"Organization ID"
//	@Param			request	body		InviteUserRequest	true	"Invite user payload"
//	@Success		200		{object}	InviteUserResponse
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		409		{object}	map[string]string	"User already member"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/organizations/{orgId}/invite [post]
func (app *application) inviteUserToOrganizationHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := getUserFromCtx(r)
	if !ok {
		app.log.Errorw("unauthorized invite attempt")
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized"))
		return
	}

	orgID := chi.URLParam(r, "orgId")
	if orgID == "" {
		app.log.Warn("organization id missing in invite request")
		app.badRequestResponse(w, r, errors.New("organization id required"))
		return
	}

	var payload InviteUserRequest

	if err := readJSON(w, r, &payload); err != nil {
		app.log.Errorw("failed to parse invite request body", "err", err)
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.log.Warnw("invite payload validation failed", "err", err)
		app.badRequestResponse(w, r, err)
		return
	}

	app.log.Infow("invite user request received",
		"inviter_user_id", user.ID,
		"organization_id", orgID,
		"email", payload.Email,
		"role", payload.Role,
	)

	// ---- Identity lookup ----
	identityReq := GetUserByEmailRequest{
		Email: payload.Email,
	}

	userResp, err := app.identityClient.Client.GetUserByEmail(
		r.Context(),
		identityReq.toProto(),
	)

	if err != nil {
		app.log.Errorw("identity service lookup failed",
			"email", payload.Email,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	app.log.Infow("identity lookup success",
		"user_id", userResp.UserId,
	)

	// ---- Add organization member ----
	orgReq := AddOrganizationMemberRequest{
		OrganizationID: orgID,
		UserID:         userResp.UserId,
		Role:           payload.Role,
	}

	_, err = app.organizationClient.Client.AddMember(
		r.Context(),
		orgReq.toProto(),
	)

	if err != nil {
		app.log.Errorw("failed to add member",
			"organization_id", orgID,
			"user_id", userResp.UserId,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	app.log.Infow("user invited successfully",
		"organization_id", orgID,
		"user_id", userResp.UserId,
	)

	resp := InviteUserResponse{
		Success: true,
	}

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.log.Errorw("failed to invite member", "err", err)
		app.internalServerError(w, r, err)
	}
}

// removeMemberHandler godoc
//
//	@Summary		Remove organization member
//	@Description	Removes a user from an organization
//	@Tags			Organizations
//	@Produce		json
//	@Security		BearerAuth
//	@Param			orgId	path		string	true	"Organization ID"
//	@Param			userId	path		string	true	"User ID"
//	@Success		200		{object}	RemoveOrganizationMemberResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/organizations/{orgId}/members/{userId} [delete]
func (app *application) removeMemberHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := getUserFromCtx(r)
	if !ok {
		app.log.Errorw("unauthorized remove member attempt")
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized"))
		return
	}

	orgID := chi.URLParam(r, "orgId")
	memberID := chi.URLParam(r, "userId")

	if orgID == "" || memberID == "" {
		app.badRequestResponse(w, r, errors.New("organization id and user id required"))
		return
	}

	app.log.Infow("remove member request received",
		"requester", user.ID,
		"organization_id", orgID,
		"user_id", memberID,
	)

	req := RemoveOrganizationMemberRequest{
		OrganizationID: orgID,
		UserID:         memberID,
	}

	_, err := app.organizationClient.Client.RemoveMember(
		r.Context(),
		req.toProto(),
	)

	if err != nil {
		app.log.Errorw("failed to remove organization member",
			"organization_id", orgID,
			"user_id", memberID,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	resp := RemoveOrganizationMemberResponse{
		Success: true,
	}

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.log.Errorw("failed to remove member", "err", err)
		app.internalServerError(w, r, err)
	}
}

// updateMemberRoleHandler godoc
//
//	@Summary		Update organization member role
//	@Description	Updates the role of a member within an organization
//	@Tags			Organizations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			orgId	path		string	true	"Organization ID"
//	@Param			userId	path		string	true	"User ID"
//	@Param			request	body		UpdateMemberRoleBody	true	"Role update payload"
//	@Success		200		{object}	UpdateOrganizationMemberRoleResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/organizations/{orgId}/members/{userId}/role [put]
func (app *application) updateMemberRoleHandler(w http.ResponseWriter, r *http.Request) {

	user, ok := getUserFromCtx(r)
	if !ok {
		app.log.Errorw("unauthorized role update attempt")
		app.unauthorizedErrorResponse(w, r, errors.New("unauthorized"))
		return
	}

	orgID := chi.URLParam(r, "orgId")
	memberID := chi.URLParam(r, "userId")

	if orgID == "" || memberID == "" {
		app.badRequestResponse(w, r, errors.New("organization id and user id required"))
		return
	}

	var payload UpdateMemberRoleBody

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.log.Infow("update member role request received",
		"requester", user.ID,
		"organization_id", orgID,
		"user_id", memberID,
		"role", payload.Role,
	)

	req := UpdateOrganizationMemberRoleRequest{
		OrganizationID: orgID,
		UserID:         memberID,
		Role:           payload.Role,
	}

	_, err := app.organizationClient.Client.UpdateMemberRole(
		r.Context(),
		req.toProto(),
	)

	if err != nil {
		app.log.Errorw("failed to update member role",
			"organization_id", orgID,
			"user_id", memberID,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	resp := UpdateOrganizationMemberRoleResponse{
		Success: true,
	}

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.log.Errorw("failed to update member permissions", "err", err)
		app.internalServerError(w, r, err)
	}
}
