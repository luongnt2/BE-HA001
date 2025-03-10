package httputil

import (
	"encoding/json"
	"net/http"
)

func ResponseWrapSuccessJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if data != nil {
		json.NewEncoder(w).Encode(struct {
			Success bool        `json:"success"`
			Data    interface{} `json:"data"`
		}{
			Success: true,
			Data:    data,
		})
	}
}

func ResponseWrapIError(w http.ResponseWriter, httpCode int, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	json.NewEncoder(w).Encode(struct {
		Success bool   `json:"success"`
		TraceID string `json:"trace_id"`
		Code    int    `json:"code"`
		Message string `json:"message,omitempty"`
	}{
		Success: false,
		Code:    code,
		Message: err.Error(),
	})
	http.Error(w, "", http.StatusInternalServerError)
}
