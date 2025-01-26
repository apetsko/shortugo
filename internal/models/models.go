package models

import "github.com/golang-jwt/jwt/v4"

type URLRecord struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	UserID string `json:"userID"`
}

type Result struct {
	Result string `json:"result"`
}

type BatchRequest struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type BatchResponse struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type AuthClaims struct {
	UserID string `json:"userid"`
	jwt.RegisteredClaims
}
