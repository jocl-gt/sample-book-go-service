package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"readcommend/services"
	"readcommend/utils"
)

type AuthorHandler struct {
	AuthorService *services.AuthorService
}

func NewAuthorHandler(authorService *services.AuthorService) *AuthorHandler {
	return &AuthorHandler{AuthorService: authorService}
}

func (h *AuthorHandler) GetAuthors(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query()) > 0 {
		utils.SendErrorResponse(w, fmt.Errorf("no query parameters allowed"), http.StatusBadRequest)
		return
	}

	authors, err := h.AuthorService.GetAuthors()
	if err != nil {
		utils.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authors); err != nil {
		utils.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
