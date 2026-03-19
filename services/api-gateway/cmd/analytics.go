package main

import (
	"fmt"
	"net/http"

	"google.golang.org/protobuf/types/known/emptypb"
)

// getVehiclesInTransit godoc
//
//	@Summary		Get vehicles in transit
//	@Description	Returns the number of vehicles currently in transit
//	@Tags			Analytics
//	@Produce		json
//	@Success		200	{object}	VehiclesInTransitResponse
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/analytics/vehicles-in-transit [get]
func (app *application) getVehiclesInTransit(w http.ResponseWriter, r *http.Request) {
	app.log.Infow("getVehiclesInTransit request received",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	resp, err := app.analyticsClient.Client.GetVehiclesInTransit(r.Context(), &emptypb.Empty{})
	if err != nil {
		app.log.Errorw("analytics service GetVehiclesInTransit failed",
			"err", err,
		)
		app.internalServerError(w, r, err)
		return
	}

	fmt.Println("resp:", resp)

	app.log.Infow("vehicles in transit fetched successfully",
		"count", resp.Count,
	)

	writeJSON(w, http.StatusOK, resp)
}

// getTripsToday godoc
//
//	@Summary		Get trips today
//	@Description	Returns the total number of trips completed today
//	@Tags			Analytics
//	@Produce		json
//	@Success		200	{object}	TripsTodayResponse
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/analytics/trips-today [get]
func (app *application) getTripsToday(w http.ResponseWriter, r *http.Request) {
	app.log.Infow("getTripsToday request received",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	resp, err := app.analyticsClient.Client.GetTripsToday(r.Context(), &emptypb.Empty{})
	if err != nil {
		app.log.Errorw("analytics service GetTripsToday failed",
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	fmt.Println("resp:", resp)

	app.log.Infow("trips today fetched successfully",
		"count", resp.Count,
	)

	writeJSON(w, http.StatusOK, resp)
}

// getTotalMembers godoc
//
//	@Summary		Get total members
//	@Description	Returns total number of members in org
//	@Tags			Analytics
//	@Produce		json
//	@Success		200	{object}	TotalMembersResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/analytics/total-members [get]
func (app *application) getTotalMembers(w http.ResponseWriter, r *http.Request) {
	app.log.Infow("getTotalMembers request received",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	// extract user id from context
	user, ok := getUserFromCtx(r)
	if !ok {
		app.log.Errorw("failed to extract user id from context")
		app.internalServerError(w, r, fmt.Errorf("userID not found"))
		return
	}

	// use your own payload
	req := GetTotalMembersRequest{
		OwnerUserID: user.ID,
	}

	// convert to grpc internally
	grpcReq := toGetTotalMembersGRPC(req)

	grpcResp, err := app.organizationClient.Client.GetTotalMembers(
		r.Context(),
		grpcReq,
	)
	if err != nil {
		app.log.Errorw("GetTotalMembers failed",
			"user_id", user.ID,
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	resp := TotalMembersResponse{
		Count: grpcResp.Count,
	}

	app.log.Infow("getTotalMembers success",
		"user_id", user.ID,
		"count", resp.Count,
	)

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.log.Errorw("failed to write response",
			"user_id", user.ID,
			"err", err,
		)
		app.internalServerError(w, r, err)
	}
}

// getActiveUsersToday godoc
//
//	@Summary		Get active users today
//	@Description	Returns number of active users today
//	@Tags			Analytics
//	@Produce		json
//	@Success		200	{object}	ActiveUsersTodayResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/analytics/active-users-today [get]
func (app *application) getActiveUsersToday(w http.ResponseWriter, r *http.Request) {
	app.log.Infow("getActiveUsersToday request received",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	resp, err := app.analyticsClient.Client.GetActiveUsersToday(r.Context(), &emptypb.Empty{})
	if err != nil {
		app.log.Errorw("GetActiveUsersToday failed",
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	app.log.Infow("getActiveUsersToday success",
		"count", resp.Count,
	)

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.log.Errorw("failed to write response",
			"err", err,
		)
		app.internalServerError(w, r, err)
	}
}

// getRecentActivity godoc
//
//	@Summary		Get recent activity
//	@Description	Returns latest analytics events
//	@Tags			Analytics
//	@Produce		json
//	@Success		200	{object}	RecentActivityResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/analytics/recent-activity [get]
func (app *application) getRecentActivity(w http.ResponseWriter, r *http.Request) {
	app.log.Infow("getRecentActivity request received",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	resp, err := app.analyticsClient.Client.GetRecentActivity(r.Context(), &emptypb.Empty{})
	if err != nil {
		app.log.Errorw("GetRecentActivity failed",
			"err", err,
		)
		app.handleGRPCError(w, r, err)
		return
	}

	app.log.Infow("getRecentActivity success",
		"events_count", len(resp.Events),
	)

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.log.Errorw("failed to write response",
			"err", err,
		)
		app.internalServerError(w, r, err)
	}
}
