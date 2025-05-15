package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/apetsko/shortugo/internal/storages/shared"
)

// ExpandURL handles requests for expanding a shortened URL.
// It retrieves the original URL from the storage and redirects the client to it.
func (h *URLHandler) ExpandURL(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path (remove the leading "/")
	ID := strings.TrimPrefix(r.URL.Path, "/")

	// Get the context from the request
	ctx := r.Context()

	// Retrieve the original URL from the storage using the ID
	URL, err := h.Storage.Get(ctx, ID)
	if err != nil {
		// Handle the case where the URL is no longer available (gone)
		if errors.Is(err, shared.ErrGone) {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusGone)
			return
		}

		// Handle the case where the URL is not found
		if errors.Is(err, shared.ErrNotFound) {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Log any other errors and respond with a generic internal server error
		h.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the "Location" header for the redirect response
	w.Header().Set("Location", URL)
	// Add the "Content-Type" header for the response
	w.Header().Add("Content-Type", "text/html")

	// Respond with a temporary redirect (HTTP 307) and the URL in the body
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err = w.Write([]byte(URL))
	if err != nil {
		// Log any error that occurs while writing the response
		h.Logger.Error(err.Error())
	}
}
