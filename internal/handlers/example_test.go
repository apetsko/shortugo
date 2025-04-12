package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/inmem"
	"go.uber.org/zap/zapcore"
)

func ExampleURLHandler_ShortenJSON() {
	storage := inmem.New()
	logger, _ := logging.New(zapcore.DebugLevel)
	handler := NewURLHandler("http://short.url", storage, logger, "secret")

	// Create a new short URL
	body := models.URLRecord{
		URL: "https://example.com",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(context.Background(), "POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()
	handler.ShortenJSON(rr, req)

	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Response Body:", rr.Body.String())
	// Output:
	// Status Code: 201
	// Response Body: {"result":"http://short.url/EAaArVRs"}
}

func ExampleURLHandler_ExpandURL() {
	storage := inmem.New()
	logger, _ := logging.New(zapcore.DebugLevel)

	handler := NewURLHandler("http://short.url", storage, logger, "secret")

	// Create a new short URL
	body := models.URLRecord{
		URL: "https://example.com",
	}
	jsonBody, _ := json.Marshal(body)
	put, _ := http.NewRequestWithContext(context.Background(), "POST", "/api/shorten", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()
	handler.ShortenJSON(rr, put)

	// Retrieve the original URL
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/EAaArVRs", nil)
	rr = httptest.NewRecorder() // reset the recorder

	handler.ExpandURL(rr, req)

	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Response Body:", rr.Body.String())
	// Output:
	// Status Code: 307
	// Response Body: https://example.com
}
func ExampleURLHandler_ListUserURLs() {
	storage := inmem.New()
	logger, _ := logging.New(zapcore.DebugLevel)
	handler := NewURLHandler("http://short.url", storage, logger, "secret")

	// Create new short URLs
	urls := []models.URLRecord{
		{URL: "https://example12.org"},
		{URL: "https://example23.com"},
	}

	var sessionCookie *http.Cookie
	rr := httptest.NewRecorder()

	for _, record := range urls {
		jsonBody, _ := json.Marshal(record)
		req, _ := http.NewRequestWithContext(context.Background(), "POST", "/api/shorten", bytes.NewBuffer(jsonBody))

		if sessionCookie != nil {
			req.AddCookie(sessionCookie)
		}

		handler.ShortenJSON(rr, req)
		bb := rr.Result()
		defer func() {
			if err := bb.Body.Close(); err != nil {
				logger.Error("Failed to close request body", "error", err.Error())
			}
		}()

		for _, c := range bb.Cookies() {
			if c.Name == "shortugo" {
				sessionCookie = c
				break
			}
		}
	}

	rr = httptest.NewRecorder()
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/api/user/urls", nil)

	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}

	handler.ListUserURLs(rr, req)

	fmt.Println("Status Code:", rr.Code)
	fmt.Println("Response Body:", rr.Body.String())

	// Output:
	// Status Code: 200
	// Response Body: [{"short_url":"http://short.url/pe7kPaGF","original_url":"https://example12.org"},{"short_url":"http://short.url/1H35Al_l","original_url":"https://example23.com"}]
}
