package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/apetsko/shortugo/internal/logging"
	mw "github.com/apetsko/shortugo/internal/middleware"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storage/shared"

	"github.com/apetsko/shortugo/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Storage interface {
	Put(ctx context.Context, r models.URLRecord) error
	PutBatch(ctx context.Context, rr []models.URLRecord) error
	Get(ctx context.Context, id string) (url string, err error)
	Ping() error
	Close() error
}

type URLHandler struct {
	baseURL string
	storage Storage
	logger  *logging.ZapLogger
}

func NewURLHandler(b string, s Storage, l *logging.ZapLogger) *URLHandler {
	return &URLHandler{
		baseURL: b,
		storage: s,
		logger:  l,
	}
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	url := string(body)
	if url == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	record := models.URLRecord{
		URL: url,
		ID:  utils.Generate(url),
	}

	ctx := r.Context()
	shortenURL, err := h.storage.Get(ctx, record.ID)
	if err != nil {
		if !errors.Is(err, shared.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = h.storage.Put(ctx, record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		shortenURL = fmt.Sprintf("%s/%s", h.baseURL, record.ID)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(shortenURL)); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *URLHandler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var record models.URLRecord

	err = json.Unmarshal(body, &record)
	if err != nil {
		h.logger.Info("Error unmarshaling request body", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if record.URL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	record.ID = utils.Generate(record.URL)

	var resp models.Result

	ctx := r.Context()
	if err = h.storage.Put(ctx, record); err != nil {
		h.logger.Error("Failed to store record", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Result = fmt.Sprintf("%s/%s", h.baseURL, record.ID)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *URLHandler) ShortenBatchJSON(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("error closing body", "error", err.Error())
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var reqs []models.BatchRequest

	err = json.Unmarshal(body, &reqs)
	if err != nil {
		h.logger.Info("Error unmarshaling request body", "error", err.Error())
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var resps []models.BatchResponse
	var records []models.URLRecord
	for _, req := range reqs {
		var resp models.BatchResponse

		if req.OriginalURL == "" {
			err := fmt.Errorf("%s: Empty URL", http.StatusText(http.StatusBadRequest))
			h.logger.Error(err.Error(), "id", req.ID)
			resp.ID = req.ID
			resp.ShortURL = err.Error()
			resps = append(resps, resp)
			continue
		}

		var record models.URLRecord
		record.URL = req.OriginalURL
		record.ID = utils.Generate(record.URL)

		records = append(records, record)

		shortURL := fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		resp = models.BatchResponse{ID: req.ID, ShortURL: shortURL}
		resps = append(resps, resp)

		ctx := r.Context()
		if err = h.storage.PutBatch(ctx, records); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resps); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *URLHandler) ExpandURL(w http.ResponseWriter, r *http.Request) {
	ID := strings.TrimPrefix(r.URL.Path, "/")

	ctx := r.Context()
	URL, err := h.storage.Get(ctx, ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", URL)
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err = w.Write([]byte(URL))
	if err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *URLHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	err := h.storage.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func SetupRouter(handler *URLHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(mw.WithLogging(handler.logger))
	r.Use(mw.GzipMiddleware(handler.logger))

	r.Post("/", handler.ShortenURL)
	r.Post("/api/shorten", handler.ShortenJSON)
	r.Post("/api/shorten/batch", handler.ShortenBatchJSON)
	r.Get("/{id}", handler.ExpandURL)
	r.Get("/ping", handler.PingDB)

	return r
}
