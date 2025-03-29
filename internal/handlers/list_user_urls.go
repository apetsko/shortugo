package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/apetsko/shortugo/internal/auth"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
)

func (h *URLHandler) ListUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromCookie(r, h.secret)
	if err != nil {
		userID, err = auth.SetCookie(w, h.secret)
		if err != nil {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	ctx := r.Context()
	records, err := h.storage.ListLinksByUserID(ctx, userID, h.baseURL)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if errors.Is(err, shared.ErrNotFound) {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var userURLs = make([]models.UserURL, 0, len(records))
	for _, record := range records {
		userURLs = append(userURLs, models.UserURL{
			ShortURL:    record.ID,
			OriginalURL: record.URL,
		})
	}

	resp, err := json.Marshal(userURLs)
	if err != nil {
		h.Logger.Error(err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Error(err.Error())
	}
}
