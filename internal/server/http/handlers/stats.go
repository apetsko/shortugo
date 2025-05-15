package handlers

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

// Stats handles internal statistics requests for the URL shortening service.
// Access is restricted to clients within a trusted subnet (TrustedSubnet).
// Returns the total number of shortened URLs and unique users in JSON format.
//
//   - Method: GET
//   - Endpoint: /api/internal/stats
//   - Headers: X-Real-IP: <client IP address>
//   - Success: 200 OK with JSON body {"Urls": <int>, "Users": <int>}
//   - Errors:
//     403 Forbidden – if TrustedSubnet is not configured or IP is outside the allowed range
//     500 Internal Server Error – if JSON encoding fails
func (h *URLHandler) Stats(w http.ResponseWriter, r *http.Request) {
	if h.TrustedSubnet == nil {
		h.Logger.Error("Forbidden. TrustedSubnet is required")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	ipStr := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	ip := net.ParseIP(ipStr)
	if ip == nil {
		h.Logger.Error("Forbidden: Invalid IP address")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !h.TrustedSubnet.Contains(ip) {
		h.Logger.Errorf("Forbidden: IP %s not in trusted subnet", ip)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	stats, err := h.Storage.Stats(r.Context())
	if err != nil {
		h.Logger.Error("Failed retrieve stats: " + err.Error())
	}

	// Marshal the list of user URLs into JSON
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err = encoder.Encode(stats); err != nil {
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
