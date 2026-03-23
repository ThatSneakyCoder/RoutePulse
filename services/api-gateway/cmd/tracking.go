package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var trackingWSUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *application) getTripCurrentLocationHandler(w http.ResponseWriter, r *http.Request) {
	tripID := chi.URLParam(r, "tripId")
	if tripID == "" {
		app.badRequestResponse(w, r, errMissingTripID)
		return
	}

	payload := GetTripTrackingRequest{
		TripID: tripID,
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.trackingClient.Client.GetTripCurrentLocation(
		r.Context(),
		payload.toTrackingCurrentProto(),
	)
	if err != nil {
		app.handleGRPCError(w, r, err)
		return
	}

	location := resp.GetLocation()
	result := TripCurrentLocationResponse{
		TripID:     location.GetTripId(),
		DriverID:   location.GetDriverId(),
		VehicleID:  location.GetVehicleId(),
		Latitude:   location.GetLatitude(),
		Longitude:  location.GetLongitude(),
		RecordedAt: location.GetRecordedAt(),
		Connection: location.GetConnection(),
	}

	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getTripLocationHistoryHandler(w http.ResponseWriter, r *http.Request) {
	tripID := chi.URLParam(r, "tripId")
	if tripID == "" {
		app.badRequestResponse(w, r, errMissingTripID)
		return
	}

	limit := int32(100)
	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		limit = int32(parsed)
	}

	payload := GetTripLocationHistoryRequest{
		TripID: tripID,
		Limit:  limit,
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.trackingClient.Client.GetTripLocationHistory(
		r.Context(),
		payload.toTrackingHistoryProto(),
	)
	if err != nil {
		app.handleGRPCError(w, r, err)
		return
	}

	result := TripLocationHistoryResponse{
		TripID: resp.GetTripId(),
		Points: []TripLocationHistoryPoint{},
	}

	for _, point := range resp.GetPoints() {
		result.Points = append(result.Points, TripLocationHistoryPoint{
			Latitude:   point.GetLatitude(),
			Longitude:  point.GetLongitude(),
			RecordedAt: point.GetRecordedAt(),
		})
	}

	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getTripGeometryHandler(w http.ResponseWriter, r *http.Request) {
	tripID := chi.URLParam(r, "tripId")
	if tripID == "" {
		app.badRequestResponse(w, r, errMissingTripID)
		return
	}

	payload := GetTripTrackingRequest{
		TripID: tripID,
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	resp, err := app.trackingClient.Client.GetTripGeometry(
		r.Context(),
		payload.toTrackingGeometryProto(),
	)
	if err != nil {
		app.handleGRPCError(w, r, err)
		return
	}

	geometry := resp.GetGeometry()
	result := TripGeometryResponse{
		TripID:          geometry.GetTripId(),
		PlannedGeometry: []CoordinateResponse{},
		ActualGeometry:  []CoordinateResponse{},
	}

	for _, point := range geometry.GetPlannedGeometry() {
		result.PlannedGeometry = append(result.PlannedGeometry, CoordinateResponse{
			Latitude:  point.GetLatitude(),
			Longitude: point.GetLongitude(),
		})
	}

	for _, point := range geometry.GetActualGeometry() {
		result.ActualGeometry = append(result.ActualGeometry, CoordinateResponse{
			Latitude:  point.GetLatitude(),
			Longitude: point.GetLongitude(),
		})
	}

	if err := app.jsonResponse(w, http.StatusOK, result); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) driverTrackingWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := trackingWSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	app.log.Infow("driver tracking websocket connected", "remote_addr", r.RemoteAddr)

	_ = conn.WriteJSON(map[string]string{
		"type":    "system",
		"message": "driver tracking websocket connected; streaming not implemented yet",
	})
}

func (app *application) dispatchTrackingWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := trackingWSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	app.log.Infow("dispatch tracking websocket connected", "remote_addr", r.RemoteAddr)

	_ = conn.WriteJSON(map[string]string{
		"type":    "system",
		"message": "dispatch tracking websocket connected; streaming not implemented yet",
	})
}
