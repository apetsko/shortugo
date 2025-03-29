package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/apetsko/shortugo/internal/auth"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
)

func (h *URLHandler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromCookie(r, h.secret)
	if err != nil {
		userID, err = auth.SetCookie(w, h.secret)
		if err != nil {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var record models.URLRecord

	err = json.Unmarshal(body, &record)
	if err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if record.URL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}
	IDlen := 8
	record.ID = utils.GenerateID(record.URL, IDlen)

	record.UserID = userID

	var resp models.Result

	url, err := h.storage.Get(r.Context(), record.ID)
	if err != nil {
		if !errors.Is(err, shared.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.Logger.Error("error get URL by ID", "error", err.Error())
			return
		}

		ctx := r.Context()
		if err = h.storage.Put(ctx, record); err != nil {
			h.Logger.Error("Failed to store record", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp.Result = fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			h.Logger.Error(err.Error())
		}
	}

	if url != "" {
		resp.Result = fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			h.Logger.Error(err.Error())
		}
		return
	}
}
