package handlers

import "net/http"

// PingDB handles the request to check the database connection.
// It pings the database and returns a status indicating whether the connection is successful.
//
// Request:
//   - Method: GET
//   - URL: /api/ping
//
// Response:
//   - 200 OK: The database connection is successful.
//   - 500 Internal Server Error: The database connection failed.
func (h *URLHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	// Ping the database to check the connection
	if err := h.Storage.Ping(); err != nil {
		// If the ping fails, respond with 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If the ping is successful, respond with 200 OK
	w.WriteHeader(http.StatusOK)
}
