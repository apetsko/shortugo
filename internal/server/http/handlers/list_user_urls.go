package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
)

// ListUserURLs handles the request to list all URLs associated with a user.
// It retrieves the user ID from the request's cookie or sets a new one if not found.
// Then, it fetches the list of URLs associated with the user ID from the storage and returns them in the response.
func (h *URLHandler) ListUserURLs(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the cookie
	userID, err := h.Auth.CookieGetUserID(r, h.Secret)
	if err != nil {
		// If the user ID is not found, set a new one
		userID, err = h.Auth.CookieSetUserID(w, h.Secret)
		if err != nil {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// Get the context from the request
	ctx := r.Context()
	// Fetch the list of URLs associated with the user ID from the storage
	records, err := h.Storage.ListLinksByUserID(ctx, h.BaseURL, userID)
	if err != nil {
		// Handle the case where no URLs are found for the user
		if errors.Is(err, shared.ErrNotFound) {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Log any other errors and respond with a generic internal server error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the list of user URLs for the response
	var userURLs = make([]models.UserURL, 0, len(records))
	for _, record := range records {
		userURLs = append(userURLs, models.UserURL{
			ShortURL:    record.ID,
			OriginalURL: record.URL,
		})
	}

	// Marshal the list of user URLs into JSON
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err = encoder.Encode(userURLs); err != nil {
		h.Logger.Error("Error marshaling user URLs:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = buf.WriteTo(w)
	if err != nil {
		h.Logger.Error(err.Error())
	}
}
