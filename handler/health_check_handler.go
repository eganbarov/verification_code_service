package handler

import (
	"encoding/json"
	"net/http"
)

type healthCheckResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"msg"`
}

type HealthCheckHandler struct{}

func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := healthCheckResponse{
		StatusCode: http.StatusOK,
		Message:    "Ok",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
