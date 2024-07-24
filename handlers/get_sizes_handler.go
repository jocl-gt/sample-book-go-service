package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"readcommend/services"
	"readcommend/utils"
)

type SizeHandler struct {
	SizeService *services.SizeService
}

func NewSizeHandler(sizeService *services.SizeService) *SizeHandler {
	return &SizeHandler{SizeService: sizeService}
}

func (h *SizeHandler) GetSizes(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) > 0 {
		utils.SendErrorResponse(w, fmt.Errorf("no query parameters allowed"), http.StatusBadRequest)
		return
	}

	sizes, err := h.SizeService.GetSizes()
	if err != nil {
		utils.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sizes); err != nil {
		utils.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
