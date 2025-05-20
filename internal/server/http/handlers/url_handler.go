package handlers

import (
	"context"
	"net"

	"github.com/apetsko/shortugo/internal/auth"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
)

// Storage defines the interface for URL storage operations.
type Storage interface {
	// Put stores a single URL record.
	Put(ctx context.Context, r models.URLRecord) error
	// PutBatch stores a batch of URL records.
	PutBatch(ctx context.Context, rr []models.URLRecord) error
	// Get retrieves a URL by its ID.
	Get(ctx context.Context, id string) (url string, err error)
	// ListLinksByUserID lists all URLs associated with a user ID.
	ListLinksByUserID(ctx context.Context, baseURL, userID string) (rr []models.URLRecord, err error)
	// DeleteUserURLs deletes URLs associated with a user ID.
	DeleteUserURLs(ctx context.Context, IDs []string, userID string) (err error)
	// Stats retrieves counts of url and users.
	Stats(ctx context.Context) (*models.Stats, error)
	// Ping checks the connection to the storage.
	Ping() error
	// Close closes the connection to the storage.
	Close() error
}

// URLHandler handles URL shortening and related operations.
type URLHandler struct {
	Auth          auth.Authenticator             // Authenticator for user authentication.
	Storage       Storage                        // Storage interface for URL operations.
	ToDelete      chan models.BatchDeleteRequest // Channel for batch delete requests.
	Logger        *logging.Logger                // Logger for logging operations.
	TrustedSubnet *net.IPNet                     // indicates trusted subnet
	Secret        string                         // Secret key for authentication.
	BaseURL       string                         // Base URL for shortened links.
}

// NewURLHandler creates a new URLHandler instance.
func NewURLHandler(baseURL string, s Storage, l *logging.Logger, secret, trustedSubnet string) *URLHandler {
	_, network, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		l.Error("Invalid trusted subnet: " + trustedSubnet)
		network = nil
	}
	return &URLHandler{
		Auth:          new(auth.Auth),                       // Initialize the authenticator.
		BaseURL:       baseURL,                              // Set the base URL.
		Storage:       s,                                    // Set the storage interface.
		Logger:        l,                                    // Set the logger.
		Secret:        secret,                               // Set the secret key.
		ToDelete:      make(chan models.BatchDeleteRequest), // Initialize the delete request channel.
		TrustedSubnet: network,                              // indicates trusted subnet
	}
}
