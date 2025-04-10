package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/utils"
)

// ShortenBatchJSON handles the batch shortening of URLs.
// It accepts a JSON array of URLs in the request body and returns their shortened versions.
//
// Request:
//   - Method: POST
//   - URL: /api/shorten/batch
//   - Headers: Content-Type: application/json
//   - Body: [{"id": "1", "original_url": "http://example.com"}, ...]
//
// Response:
//   - 201 Created: The batch shortening request is successful.
//   - 400 Bad Request: Invalid request body or JSON format.
//   - 500 Internal Server Error: User authentication failed or other server error.
func (h *URLHandler) ShortenBatchJSON(w http.ResponseWriter, r *http.Request) {
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
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			h.Logger.Error("error closing body", "error", err.Error())
		}
	}(r.Body)

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var reqs []models.BatchRequest

	// Unmarshal the JSON array of batch requests
	if err = json.Unmarshal(body, &reqs); err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resps []models.BatchResponse
	var records []models.URLRecord

	// Process each batch request
	for _, req := range reqs {
		var resp models.BatchResponse

		// Validate the original URL
		if req.OriginalURL == "" {
			errStr := http.StatusText(http.StatusBadRequest) + ": Empty URL"
			h.Logger.Error(errStr, "id", req.ID)
			resp.ID = req.ID
			resp.ShortURL = errStr
			resps = append(resps, resp)
			continue
		}

		// Generate a unique ID for the URL
		IDlen := 8
		var record = models.URLRecord{
			URL:    req.OriginalURL,
			ID:     utils.GenerateID(req.OriginalURL, IDlen),
			UserID: userID,
		}

		records = append(records, record)

		// Create the shortened URL
		shortURL := h.BaseURL + "/" + record.ID
		resp = models.BatchResponse{ID: req.ID, ShortURL: shortURL}
		resps = append(resps, resp)
	}

	// Store the batch of URL records
	ctx := r.Context()
	if err = h.Storage.PutBatch(ctx, records); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resps); err != nil {
		h.Logger.Error(err.Error())
	}
}
