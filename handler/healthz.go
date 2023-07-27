package handler

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hr := &model.HealthzResponse{Message: "OK"}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(hr)
	if err != nil {
		log.Println(err)
	}
}
