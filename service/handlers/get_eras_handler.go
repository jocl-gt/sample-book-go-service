package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"readcommend/services"
	"readcommend/utils"
)

type ErasHandler struct {
	ErasService *services.ErasService
}

func NewErasHandler(erasService *services.ErasService) *ErasHandler {
	return &ErasHandler{ErasService: erasService}
}

func (h *ErasHandler) GetEras(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) > 0 {
		utils.SendErrorResponse(w, fmt.Errorf("no query parameters allowed"), http.StatusBadRequest)
		return
	}

	eras, err := h.ErasService.GetEras()
	if err != nil {
		utils.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(eras); err != nil {
		utils.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
