package main

import (
	"net/http"
)

// createVehicleHandler godoc
//
//	@Summary		Create vehicle
//	@Description	Creates a vehicle in the fleet service
//	@Tags			Fleet
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateVehicleRequest	true	"Create vehicle payload"
//	@Success		201		{object}	VehicleResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/fleet/vehicles [post]
func (app *application) createVehicleHandler(w http.ResponseWriter, r *http.Request) {

	app.log.Infow("create vehicle request received",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	var payload CreateVehicleRequest

	if err := readJSON(w, r, &payload); err != nil {
		app.log.Warnw("failed to read create vehicle payload",
			"err", err,
		)
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.log.Warnw("create vehicle validation failed",
			"err", err,
		)
		app.badRequestResponse(w, r, err)
		return
	}

	app.log.Infow("calling fleet service CreateVehicle",
		"organization_id", payload.OrganizationID,
		"plate_number", payload.PlateNumber,
	)

	resp, err := app.fleetClient.Client.CreateVehicle(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {

		app.log.Errorw("fleet service CreateVehicle failed",
			"err", err,
			"organization_id", payload.OrganizationID,
		)

		app.handleGRPCError(w, r, err)
		return
	}

	v := resp.GetVehicle()

	result := VehicleResponse{
		VehicleID:      v.VehicleId,
		OrganizationID: v.OrganizationId,
		PlateNumber:    v.PlateNumber,
		VehicleType:    v.VehicleType,
		Capacity:       v.Capacity,
		Status:         v.Status,
		CreatedAt:      v.CreatedAt,
	}

	app.log.Infow("vehicle created successfully",
		"vehicle_id", result.VehicleID,
	)

	if err := app.jsonResponse(w, http.StatusCreated, result); err != nil {
		app.log.Errorw("failed to write response",
			"err", err,
		)
		app.internalServerError(w, r, err)
	}
}

// listVehiclesHandler godoc
//
//	@Summary		List vehicles
//	@Description	Returns vehicles belonging to an organization
//	@Tags			Fleet
//	@Produce		json
//	@Param			organization_id	query		string	true	"Organization ID"
//	@Success		200	{object}	ListVehiclesResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		BearerAuth
//	@Router			/fleet/vehicles [get]
func (app *application) listVehiclesHandler(w http.ResponseWriter, r *http.Request) {

	orgID := r.URL.Query().Get("organization_id")

	if orgID == "" {
		app.badRequestResponse(w, r, errMissingOrganizationID)
		return
	}

	payload := ListVehiclesRequest{
		OrganizationID: orgID,
	}

	app.log.Infow("listing vehicles",
		"organization_id", orgID,
	)

	resp, err := app.fleetClient.Client.ListVehicles(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {

		app.log.Errorw("fleet service ListVehicles failed",
			"err", err,
			"organization_id", orgID,
		)

		app.handleGRPCError(w, r, err)
		return
	}

	result := ListVehiclesResponse{}

	for _, v := range resp.Vehicles {

		result.Vehicles = append(result.Vehicles, VehicleResponse{
			VehicleID:      v.VehicleId,
			OrganizationID: v.OrganizationId,
			PlateNumber:    v.PlateNumber,
			VehicleType:    v.VehicleType,
			Capacity:       v.Capacity,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt,
		})
	}

	app.log.Infow("vehicles fetched successfully",
		"organization_id", orgID,
		"count", len(result.Vehicles),
	)

	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {
		app.log.Errorw("failed to write response",
			"err", err,
		)
		app.internalServerError(w, r, err)
	}
}

func (app *application) getVehicleHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) createDriverHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) listDriversHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) createTripHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) listTripsHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) startTripHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) completeTripHandler(w http.ResponseWriter, r *http.Request) {

}
