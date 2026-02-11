package main

import (
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

	app.log.Infow("trips today fetched successfully",
		"count", resp.Count,
	)

	writeJSON(w, http.StatusOK, resp)
}
