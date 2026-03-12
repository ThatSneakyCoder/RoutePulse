package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

// getVehicleHandler godoc
//
//	@Summary		Get vehicle
//	@Description	Get vehicle details by vehicle ID
//	@Tags			Fleet
//	@Produce		json
//	@Security		BearerAuth
//	@Param			vehicleId	path	string	true	"Vehicle ID"
//	@Success		200	{object}	VehicleResponse
//	@Failure		404	{object}	map[string]string	"Vehicle not found"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/fleet/vehicles/{vehicleId} [get]
func (app *application) getVehicleHandler(w http.ResponseWriter, r *http.Request) {

	vehicleID := chi.URLParam(r, "vehicleId")

	app.log.Infow(
		"get vehicle request received",
		"vehicle_id", vehicleID,
	)

	payload := GetVehicleRequest{
		VehicleID: vehicleID,
	}

	resp, err := app.fleetClient.Client.GetVehicle(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {

		app.log.Errorw(
			"fleet service GetVehicle failed",
			"vehicle_id", vehicleID,
			"error", err,
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

	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {
		app.log.Errorw(
			"failed to write vehicle response",
			"error", err,
		)
		app.internalServerError(w, r, err)
	}
}

// createDriverHandler godoc
//
//	@Summary		Create driver
//	@Description	Create a new driver in the fleet service
//	@Tags			Fleet
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	CreateDriverRequest	true	"Create driver payload"
//	@Success		201	{object}	DriverResponse
//	@Failure		400	{object}	map[string]string	"Validation error"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/fleet/drivers [post]
func (app *application) createDriverHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateDriverRequest

	if err := readJSON(w, r, &payload); err != nil {

		app.log.Warnw(
			"failed to read create driver payload",
			"error", err,
		)

		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {

		app.log.Warnw(
			"create driver validation failed",
			"error", err,
		)

		app.badRequestResponse(w, r, err)
		return
	}

	app.log.Infow(
		"calling fleet service CreateDriver",
		"organization_id", payload.OrganizationID,
	)

	resp, err := app.fleetClient.Client.CreateDriver(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {

		app.log.Errorw(
			"fleet service CreateDriver failed",
			"error", err,
		)

		app.handleGRPCError(w, r, err)
		return
	}

	d := resp.GetDriver()

	result := DriverResponse{
		DriverID:       d.DriverId,
		OrganizationID: d.OrganizationId,
		FirstName:      d.FirstName,
		LastName:       d.LastName,
		VehicleID:      d.VehicleId,
		Status:         d.Status,
		CreatedAt:      d.CreatedAt,
	}

	if err := app.jsonResponse(w, http.StatusCreated, result); err != nil {

		app.log.Errorw(
			"failed to write driver response",
			"error", err,
		)

		app.internalServerError(w, r, err)
	}
}

// listDriversHandler godoc
//
//	@Summary		List drivers
//	@Description	Get all drivers belonging to an organization
//	@Tags			Fleet
//	@Produce		json
//	@Security		BearerAuth
//	@Param			organization_id	query	string	true	"Organization ID"
//	@Success		200	{object}	ListDriversResponse
//	@Failure		400	{object}	map[string]string	"Invalid request"
//	@Failure		500	{object}	map[string]string	"Internal server error"
//	@Router			/fleet/drivers [get]
func (app *application) listDriversHandler(w http.ResponseWriter, r *http.Request) {

	orgID := r.URL.Query().Get("organization_id")

	payload := ListDriversRequest{
		OrganizationID: orgID,
	}

	if err := Validate.Struct(payload); err != nil {

		app.log.Warnw(
			"list drivers validation failed",
			"error", err,
		)

		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.fleetClient.Client.ListDrivers(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {

		app.log.Errorw(
			"fleet service ListDrivers failed",
			"organization_id", orgID,
			"error", err,
		)

		app.handleGRPCError(w, r, err)
		return
	}

	result := ListDriversResponse{}

	for _, d := range resp.Drivers {

		result.Drivers = append(result.Drivers, DriverResponse{
			DriverID:       d.DriverId,
			OrganizationID: d.OrganizationId,
			FirstName:      d.FirstName,
			LastName:       d.LastName,
			VehicleID:      d.VehicleId,
			Status:         d.Status,
			CreatedAt:      d.CreatedAt,
		})
	}

	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {

		app.log.Errorw(
			"failed to write drivers response",
			"error", err,
		)

		app.internalServerError(w, r, err)
	}
}

// createTripHandler godoc
//
//	@Summary		Create trip
//	@Description	Create a new trip
//	@Tags			Fleet
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	CreateTripRequest	true	"Create trip payload"
//	@Success		201	{object}	TripResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/fleet/trips [post]
func (app *application) createTripHandler(w http.ResponseWriter, r *http.Request) {

	var payload CreateTripRequest

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.fleetClient.Client.CreateTrip(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {
		app.handleGRPCError(w, r, err)
		return
	}

	t := resp.GetTrip()

	result := TripResponse{
		TripID:         t.TripId,
		OrganizationID: t.OrganizationId,
		VehicleID:      t.VehicleId,
		DriverID:       t.DriverId,
		Status:         t.Status,
		CreatedAt:      t.CreatedAt,
	}

	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {

		app.log.Errorw(
			"failed to createTrip",
			"error", err,
		)

		app.internalServerError(w, r, err)
	}
}

// listTripsHandler godoc
//
//	@Summary		List trips
//	@Description	Get trips for organization
//	@Tags			Fleet
//	@Produce		json
//	@Security		BearerAuth
//	@Param			organization_id	query	string	true	"Organization ID"
//	@Success		200	{object}	ListTripsResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/fleet/trips [get]
func (app *application) listTripsHandler(w http.ResponseWriter, r *http.Request) {

	orgID := r.URL.Query().Get("organization_id")

	payload := ListTripsRequest{
		OrganizationID: orgID,
	}

	resp, err := app.fleetClient.Client.ListTrips(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {
		app.handleGRPCError(w, r, err)
		return
	}

	result := ListTripsResponse{}

	for _, t := range resp.Trips {

		result.Trips = append(result.Trips, TripResponse{
			TripID:         t.TripId,
			OrganizationID: t.OrganizationId,
			VehicleID:      t.VehicleId,
			DriverID:       t.DriverId,
			Status:         t.Status,
			CreatedAt:      t.CreatedAt,
		})
	}

	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {

		app.log.Errorw(
			"failed to list trips",
			"error", err,
		)

		app.internalServerError(w, r, err)
	}
}

// updateVehicleHandler godoc
//
//	@Summary		Update vehicle metadata
//	@Description	Update vehicle plate, type, or capacity
//	@Tags			Fleet
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	UpdateVehicleRequest	true	"Update vehicle payload"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/fleet/vehicles/{vehicleId} [put]
func (app *application) updateVehicleHandler(w http.ResponseWriter, r *http.Request) {

	var payload UpdateVehicleRequest

	if err := readJSON(w, r, &payload); err != nil {
		app.log.Errorw("failed reading update vehicle payload", "error", err)
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.log.Errorw("update vehicle validation failed", "error", err)
		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.fleetClient.Client.UpdateVehicle(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {
		app.log.Errorw("fleet service UpdateVehicle failed", "error", err)
		app.handleGRPCError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {

		app.log.Errorw(
			"failed to update vehicle",
			"error", err,
		)

		app.internalServerError(w, r, err)
	}
}

// updateVehicleStatusHandler godoc
//
//	@Summary		Update vehicle status
//	@Description	Activate or deactivate vehicle
//	@Tags			Fleet
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	UpdateVehicleStatusRequest	true	"Vehicle status payload"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/fleet/vehicles/{vehicleId}/status [patch]
func (app *application) updateVehicleStatusHandler(w http.ResponseWriter, r *http.Request) {

	var payload UpdateVehicleStatusRequest

	if err := readJSON(w, r, &payload); err != nil {
		app.log.Errorw("failed reading vehicle status payload", "error", err)
		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.fleetClient.Client.UpdateVehicleStatus(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {
		app.log.Errorw("fleet service UpdateVehicleStatus failed", "error", err)
		app.handleGRPCError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {

		app.log.Errorw(
			"failed to update vehicle status",
			"error", err,
		)

		app.internalServerError(w, r, err)
	}
}

// updateDriverHandler godoc
//
//	@Summary		Update driver profile
//	@Description	Update driver first or last name
//	@Tags			Fleet
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	UpdateDriverRequest	true	"Update driver payload"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/fleet/drivers/{driverId} [put]
func (app *application) updateDriverHandler(w http.ResponseWriter, r *http.Request) {

	var payload UpdateDriverRequest

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.fleetClient.Client.UpdateDriver(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {
		app.log.Errorw("fleet service UpdateDriver failed", "error", err)
		app.handleGRPCError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {

		app.log.Errorw(
			"failed to update driver",
			"error", err,
		)

		app.internalServerError(w, r, err)
	}
}

// updateDriverStatusHandler godoc
//
//	@Summary		Update driver status
//	@Description	Activate or deactivate driver
//	@Tags			Fleet
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body	UpdateDriverStatusRequest	true	"Driver status payload"
//	@Success		200	{object}	map[string]bool
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/fleet/drivers/{driverId}/status [patch]
func (app *application) updateDriverStatusHandler(w http.ResponseWriter, r *http.Request) {

	var payload UpdateDriverStatusRequest

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.fleetClient.Client.UpdateDriverStatus(
		r.Context(),
		payload.toProto(),
	)

	if err != nil {
		app.log.Errorw("fleet service UpdateDriverStatus failed", "error", err)
		app.handleGRPCError(w, r, err)
		return
	}

		if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {

		app.log.Errorw(
			"failed to update driver status",
			"error", err,
		)

		app.internalServerError(w, r, err)
	}
}

func (app *application) startTripHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) completeTripHandler(w http.ResponseWriter, r *http.Request) {

}