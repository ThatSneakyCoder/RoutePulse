package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"github.com/ThatSneakyCoder/RoutePulse/shared/rabbitmq"
)

type trackingStreamMessage struct {
	Type string                                             `json:"type"`
	Data rabbitmq.TrackingDriverLocationUpdatedEventPayload `json:"data"`
}

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
	app.serveTrackingWebSocket(w, r, "driver")
}

func (app *application) dispatchTrackingWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	app.serveTrackingWebSocket(w, r, "dispatch")
}

func (app *application) serveTrackingWebSocket(w http.ResponseWriter, r *http.Request, clientType string) {
	conn, err := trackingWSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		app.log.Errorw("tracking websocket upgrade failed", "client_type", clientType, "error", err)
		return
	}
	defer conn.Close()

	app.log.Infow(
		clientType+" tracking websocket connected",
		"client_type", clientType,
		"remote_addr", r.RemoteAddr,
	)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			app.log.Infow(
				clientType+" tracking websocket disconnected",
				"client_type", clientType,
				"remote_addr", r.RemoteAddr,
				"error", err,
			)
			break
		}

		var payload trackingStreamMessage
		if err := json.Unmarshal(message, &payload); err != nil {
			app.log.Warnw(
				clientType+" tracking websocket message decode failed",
				"client_type", clientType,
				"remote_addr", r.RemoteAddr,
				"message_type", messageType,
				"raw_message", string(message),
				"error", err,
			)
			continue
		}

		app.log.Infow(
			clientType+" tracking update received",
			"client_type", clientType,
			"remote_addr", r.RemoteAddr,
			"message_type", payload.Type,
			"trip_id", payload.Data.TripID,
			"driver_id", payload.Data.DriverID,
			"vehicle_id", payload.Data.VehicleID,
			"latitude", payload.Data.Latitude,
			"longitude", payload.Data.Longitude,
			"recorded_at", payload.Data.RecordedAt,
			"sequence", payload.Data.Sequence,
		)

		if err := app.eventPublisher.PublishTrackingLocationUpdated(r.Context(), payload.Data); err != nil {
			app.log.Errorw(
				"failed to publish tracking location update",
				"client_type", clientType,
				"trip_id", payload.Data.TripID,
				"driver_id", payload.Data.DriverID,
				"vehicle_id", payload.Data.VehicleID,
				"sequence", payload.Data.Sequence,
				"error", err,
			)
		}
	}
}
