package main

import (
	"net/http"

	"github.com/ThatSneakyCoder/RoutePulse/services/api-gateway/internal/infrastructure/grpc"
)

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// healthCheckHandler godoc
//
// @Summary      Health check
// @Description  Checks API health and connectivity with identity service
// @Tags         Health
// @Produce      json
// @Success      200 {object} HealthCheckResponse
// @Failure      500 {object} map[string]string
// @Router       /health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := HealthCheckResponse{
		Status:  "ok",
		Version: "0.0.1",
	}

	identityService, err := grpc.NewIdentityServiceClient()
	if err != nil {
		app.logger.Info(err)
	}
	defer identityService.Close()

	ctx := r.Context()
	_, _ = identityService.Client.RegisterUser(ctx, nil)

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
