// Package models defines the data structures used throughout the application.
// It includes models for URL records, batch operations, and user-specific URL data.
package models

// URLRecord represents a record of a shortened URL.
type URLRecord struct {
	ID      string `json:"id"`      // Unique identifier for the URL record.
	URL     string `json:"url"`     // Original URL.
	UserID  string `json:"userid"`  // ID of the user who created the URL.
	Deleted bool   `json:"deleted"` // Flag indicating if the URL is deleted.
}

// Result represents a generic result message.
type Result struct {
	Result string `json:"result"` // Result message.
}

// BatchDeleteRequest represents a request to delete multiple URLs.
type BatchDeleteRequest struct {
	UserID string   // ID of the user requesting the deletion.
	Ids    []string // List of URL IDs to be deleted.
}

// BatchRequest represents a request to shorten multiple URLs.
type BatchRequest struct {
	ID          string `json:"correlation_id"` // Correlation ID for the batch request.
	OriginalURL string `json:"original_url"`   // Original URL to be shortened.
}

// BatchResponse represents a response for a batch URL shortening request.
type BatchResponse struct {
	ID       string `json:"correlation_id"` // Correlation ID for the batch response.
	ShortURL string `json:"short_url"`      // Shortened URL.
}

// UserURL represents a user's URL with both short and original versions.
type UserURL struct {
	ShortURL    string `json:"short_url"`    // Shortened URL.
	OriginalURL string `json:"original_url"` // Original URL.
}

// Stats presents count of users and urls
type Stats struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}
