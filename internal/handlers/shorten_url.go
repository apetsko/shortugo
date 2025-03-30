package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
)

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	userID, err := h.auth.CookieGetUserID(r, h.secret)
	if err != nil {
		userID, err = h.auth.CookieSetUserID(w, h.secret)
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

	url := string(body)
	if url == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	IDlen := 8
	record := models.URLRecord{
		URL:    url,
		ID:     utils.GenerateID(url, IDlen),
		UserID: userID,
	}

	ctx := r.Context()

	shortenURL, err := h.storage.Get(ctx, record.ID)
	if err != nil {
		if !errors.Is(err, shared.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.Logger.Error("error get URL by ID", "error", err.Error())
			return
		}

		if err = h.storage.Put(ctx, record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newShortenURL := fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(newShortenURL)); err != nil {
			h.Logger.Error(err.Error())
		}
		return
	}

	if shortenURL != "" {
		w.WriteHeader(http.StatusConflict)
		shortenURL = fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		if _, err := w.Write([]byte(shortenURL)); err != nil {
			h.Logger.Error(err.Error())
		}
		return
	}
}
