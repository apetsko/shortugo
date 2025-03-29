package handlers

import (
	"context"

	"github.com/apetsko/shortugo/internal/auth"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
)

type Storage interface {
	Put(ctx context.Context, r models.URLRecord) error
	PutBatch(ctx context.Context, rr []models.URLRecord) error
	Get(ctx context.Context, id string) (url string, err error)
	ListLinksByUserID(ctx context.Context, baseURL, userID string) (rr []models.URLRecord, err error)
	DeleteUserURLs(ctx context.Context, IDs []string, userID string) (err error)
	Ping() error
	Close() error
}

type URLHandler struct {
	auth     auth.Authenticator
	baseURL  string
	storage  Storage
	secret   string
	ToDelete chan models.BatchDeleteRequest
	Logger   *logging.Logger
}

func NewURLHandler(baseURL string, s Storage, l *logging.Logger, secret string) *URLHandler {
	return &URLHandler{
		auth:     new(auth.Auth),
		baseURL:  baseURL,
		storage:  s,
		Logger:   l,
		secret:   secret,
		ToDelete: make(chan models.BatchDeleteRequest),
	}
}
