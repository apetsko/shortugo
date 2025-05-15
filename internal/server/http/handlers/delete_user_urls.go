package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/apetsko/shortugo/internal/models"
)

// DeleteUserURLs handles the deletion of multiple user URLs.
// It accepts a JSON array of short URL IDs in the request body and schedules them for deletion asynchronously.
//
// Request:
//   - Method: DELETE
//   - URL: /api/user/urls
//   - Headers: Content-Type: application/json
//   - Body: ["abc123", "xyz789"]
//
// Response:
//   - 202 Accepted: The deletion request is accepted for processing.
//   - 400 Bad Request: Invalid request body or JSON format.
//   - 401 Unauthorized: User authentication failed.
//
// The function retrieves the user ID from a cookie, validates the request body, and then
// sends the batch delete request to the processing channel asynchronously.
func (h *URLHandler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the cookie
	userID, err := h.Auth.CookieGetUserID(r, h.Secret)
	if err != nil {
		h.Logger.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ids []string

	// Unmarshal the JSON array of short URL IDs
	err = json.Unmarshal(body, &ids)
	if err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send the batch delete request asynchronously
	go func() {
		h.ToDelete <- models.BatchDeleteRequest{Ids: ids, UserID: userID}
	}()

	// Respond with 202 Accepted and return the requested IDs
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(ids); err != nil {
		h.Logger.Error(err.Error())
	}
}
