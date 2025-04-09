package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
)

// ShortenURL handles the shortening of a single URL.
// It accepts a plain text URL in the request body and returns the shortened version.
//
// Request:
//   - Method: POST
//   - URL: /api/shorten
//   - Headers: Content-Type: text/plain
//   - Body: http://example.com
//
// Response:
//   - 201 Created: The URL shortening request is successful.
//   - 400 Bad Request: Invalid request body or empty URL.
//   - 409 Conflict: The URL already exists.
//   - 500 Internal Server Error: User authentication failed or other server error.
func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
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

	// Ensure the request body is closed after reading
	defer r.Body.Close()

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Convert the request body to a string
	url := string(body)
	if url == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the URL
	IDlen := 8
	record := models.URLRecord{
		URL:    url,
		ID:     utils.GenerateID(url, IDlen),
		UserID: userID,
	}

	ctx := r.Context()

	// Check if the URL already exists in the storage
	shortenURL, err := h.Storage.Get(ctx, record.ID)
	if err != nil {
		if !errors.Is(err, shared.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.Logger.Error("error get URL by ID", "error", err.Error())
			return
		}

		// Store the new URL record
		if err = h.Storage.Put(ctx, record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare the response with the shortened URL
		newShortenURL := h.BaseURL + "/" + record.ID
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(newShortenURL)); err != nil {
			h.Logger.Error(err.Error())
		}
		return
	}

	// If the URL already exists, respond with 409 Conflict
	if shortenURL != "" {
		w.WriteHeader(http.StatusConflict)
		shortenURL = h.BaseURL + "/" + record.ID
		if _, err := w.Write([]byte(shortenURL)); err != nil {
			h.Logger.Error(err.Error())
		}
		return
	}
}
