package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/apetsko/shortugo/internal/models"
)

func (h *URLHandler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := h.auth.UserIDFromCookie(r, h.secret)
	if err != nil {
		h.Logger.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ids []string

	err = json.Unmarshal(body, &ids)
	if err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	go func() {
		h.ToDelete <- models.BatchDeleteRequest{Ids: ids, UserID: userID}
	}()

	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(ids); err != nil {
		h.Logger.Error(err.Error())
	}
}
