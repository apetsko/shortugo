package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/inmem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLHandler_ShortenJSON(t *testing.T) {
	logger, err := logging.New()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}

	u := "http://localhost:8080"
	handler := NewURLHandler(u, inmem.New(), logger, "fortytwo")
	type want struct {
		ID          string
		code        int
		ContentType string
	}
	tests := []struct {
		name   string
		record models.URLRecord
		want   want
	}{
		{
			name:   "positive test #1",
			record: models.URLRecord{URL: "https://practicum.yandex.ru/,", UserID: "55"},
			want: want{
				code:        201,
				ID:          "http://localhost:8080/3P9NwpqM",
				ContentType: "application/json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url, err := json.Marshal(test.record)
			require.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, test.want.ID, bytes.NewBuffer(url))
			w := httptest.NewRecorder()
			handler.ShortenJSON(w, request)
			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()

			b, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			var resp models.Result
			err = json.Unmarshal(b, &resp)
			require.NoError(t, err)

			assert.Equal(t, test.want.ContentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.ID, resp.Result)
		})
	}
}
