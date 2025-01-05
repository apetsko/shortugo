package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apetsko/shortugo/internal/storage/inmem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLHandler_ShortenURL(t *testing.T) {
	handler := NewURLHandler(inmem.New())
	type want struct {
		ID   string
		code int
	}
	tests := []struct {
		name string
		URL  string
		want want
	}{
		{
			name: "positive test #1",
			URL:  "https://practicum.yandex.ru/",
			want: want{

				code: 201,
				ID:   "http://localhost:8080/EwHXdJfB",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.want.ID, strings.NewReader(test.URL))
			w := httptest.NewRecorder()
			handler.ShortenURL(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Equal(t, test.want.ID, string(resBody))
		})
	}
}

func TestURLHandler_ExpandURL(t *testing.T) {
	handler := NewURLHandler(inmem.New())
	type want struct {
		code     int
		Location string
	}
	tests := []struct {
		shortenURL string
		name       string
		want       want
	}{
		{
			name:       "positive test #1",
			shortenURL: "http://localhost:8080/EwHXdJfB",
			want: want{
				code:     307,
				Location: "https://practicum.yandex.ru/",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler.storage.Put(test.want.Location)
			request := httptest.NewRequest(http.MethodGet, test.shortenURL, nil)
			w := httptest.NewRecorder()
			handler.ExpandURL(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.Location, res.Header.Get("Location"))
		})
	}
}
