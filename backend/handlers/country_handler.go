package handlers

import (
	"Pinspire/backend/service"
	"encoding/json"
	"net/http"
)

type CountryHandler struct {
	Service *service.CountryService
}

func (h *CountryHandler) SetCountry(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Service.SetUserCountry(r.Context(), request.Code); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
