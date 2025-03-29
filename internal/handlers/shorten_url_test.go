package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/storages/inmem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLHandler_ShortenURL(t *testing.T) {
	logger, err := logging.New()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	u := "http://localhost:8080"
	handler := NewURLHandler(u, inmem.New(), logger, "fortytwo")
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
				ID:   "http://localhost:8080/QrPnX5IU",
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
