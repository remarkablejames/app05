package handlers

import (
	"app05/internal/infrastructure/services"
	"app05/pkg/utils"
	"net/http"
)

type ServerHandler struct {
	serverService *services.ServerService
}

func NewServerStatusHandler(serverService *services.ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

// HealthResponse defines the structure of the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Uptime    string `json:"uptime"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// HealthCheck handles the GET /health route
func (h *ServerHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Get the health status
	statusResponse := h.serverService.GetHealthStatus()
	// Return the health response as a JSON response
	utils.SendJSON(w, statusResponse)
}
