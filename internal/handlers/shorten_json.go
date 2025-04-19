package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
)

// ShortenJSON handles the shortening of a single URL.
// It accepts a JSON object with the original URL in the request body and returns the shortened version.
//
// Request:
//   - Method: POST
//   - URL: /api/shorten
//   - Headers: Content-Type: application/json
//   - Body: {"url": "http://example.com"}
//
// Response:
//   - 201 Created: The URL shortening request is successful.
//   - 400 Bad Request: Invalid request body or JSON format.
//   - 409 Conflict: The URL already exists.
//   - 500 Internal Server Error: User authentication failed or other server error.
func (h *URLHandler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
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
	defer func() {
		if err2 := r.Body.Close(); err2 != nil {
			h.Logger.Error("Failed to close request body", "error", err2.Error())
		}
	}()

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var record models.URLRecord

	// Unmarshal the JSON object into a URLRecord
	err = json.Unmarshal(body, &record)
	if err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the original URL
	if record.URL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the URL
	IDlen := 8
	record.ID = utils.GenerateID(record.URL, IDlen)
	record.UserID = userID

	var resp models.Result

	// Check if the URL already exists in the storage
	url, err := h.Storage.Get(r.Context(), record.ID)
	if err != nil {
		if !errors.Is(err, shared.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.Logger.Error("error get URL by ID", "error", err.Error())
			return
		}

		// Store the new URL record
		ctx := r.Context()
		if err = h.Storage.Put(ctx, record); err != nil {
			h.Logger.Error("Failed to store record", "error", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Prepare the response with the shortened URL
		resp.Result = h.BaseURL + "/" + record.ID
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			h.Logger.Error(err.Error())
		}
	}

	// If the URL already exists, respond with 409 Conflict
	if url != "" {
		resp.Result = h.BaseURL + "/" + record.ID
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			h.Logger.Error(err.Error())
		}
		return
	}
}
