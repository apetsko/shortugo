package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/apetsko/shortugo/internal/auth"
	"github.com/apetsko/shortugo/internal/logging"
	mw "github.com/apetsko/shortugo/internal/middleware"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/storages/shared"
	"github.com/apetsko/shortugo/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Storage interface {
	Put(ctx context.Context, r models.URLRecord) error
	PutBatch(ctx context.Context, rr []models.URLRecord) error
	Get(ctx context.Context, id string) (url string, err error)
	GetLinksByUserID(ctx context.Context, baseURL, userID string) (rr []models.URLRecord, err error)
	DeleteUserURLs(ctx context.Context, IDs []string, userID string) (err error)
	Ping() error
	Close() error
}

type URLHandler struct {
	baseURL  string
	storage  Storage
	secret   string
	ToDelete chan models.BatchDeleteRequest
	logger   *logging.ZapLogger
}

func NewURLHandler(baseURL string, s Storage, l *logging.ZapLogger, secret string) *URLHandler {
	return &URLHandler{
		baseURL:  baseURL,
		storage:  s,
		logger:   l,
		secret:   secret,
		ToDelete: make(chan models.BatchDeleteRequest),
	}
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.CookieUserID(r, h.secret)
	if err != nil {
		userID, err = auth.SetUserIDCookie(w, h.secret)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

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

	IDlen := 8
	record := models.URLRecord{
		URL: url,
		ID:  utils.GenerateID(url, IDlen),
	}

	record.UserID = userID

	ctx := r.Context()

	shortenURL, err := h.storage.Get(ctx, record.ID)
	if err != nil {
		if !errors.Is(err, shared.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.logger.Error("error get URL by ID", "error", err.Error())
			return
		}

		if err = h.storage.Put(ctx, record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newShortenURL := fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(newShortenURL)); err != nil {
			h.logger.Error(err.Error())
		}
		return
	}

	if shortenURL != "" {
		w.WriteHeader(http.StatusConflict)
		shortenURL = fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		if _, err := w.Write([]byte(shortenURL)); err != nil {
			h.logger.Error(err.Error())
		}
		return
	}
}

func (h *URLHandler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.CookieUserID(r, h.secret)
	if err != nil {
		userID, err = auth.SetUserIDCookie(w, h.secret)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

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
	IDlen := 8
	record.ID = utils.GenerateID(record.URL, IDlen)

	record.UserID = userID

	var resp models.Result

	url, err := h.storage.Get(r.Context(), record.ID)
	if err != nil {
		if !errors.Is(err, shared.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.logger.Error("error get URL by ID", "error", err.Error())
			return
		}

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

	if url != "" {
		resp.Result = fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			h.logger.Error(err.Error())
		}
		return
	}
}

func (h *URLHandler) ShortenBatchJSON(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.CookieUserID(r, h.secret)
	if err != nil {
		userID, err = auth.SetUserIDCookie(w, h.secret)
		if err != nil {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("error closing body", "error", err.Error())
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var reqs []models.BatchRequest

	if err = json.Unmarshal(body, &reqs); err != nil {
		h.logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resps []models.BatchResponse
	var records []models.URLRecord

	for _, req := range reqs {
		var resp models.BatchResponse

		if req.OriginalURL == "" {
			errStr := fmt.Sprintf("%s: Empty URL", http.StatusText(http.StatusBadRequest))
			h.logger.Error(errStr, "id", req.ID)
			resp.ID = req.ID
			resp.ShortURL = errStr
			resps = append(resps, resp)
			continue
		}

		var record models.URLRecord
		record.URL = req.OriginalURL
		IDlen := 8
		record.ID = utils.GenerateID(record.URL, IDlen)
		record.UserID = userID

		records = append(records, record)

		shortURL := fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		resp = models.BatchResponse{ID: req.ID, ShortURL: shortURL}
		resps = append(resps, resp)
	}
	fmt.Println(resps)
	ctx := r.Context()
	if err = h.storage.PutBatch(ctx, records); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resps); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *URLHandler) AllUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromCookie(r, h.secret)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	records, err := h.storage.GetLinksByUserID(ctx, userID, h.baseURL)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if errors.Is(err, shared.ErrNotFound) {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var userURLs = make([]models.UserURL, 0, len(records))
	for _, record := range records {
		userURLs = append(userURLs, models.UserURL{
			ShortURL:    record.ID,
			OriginalURL: record.URL,
		})
	}

	resp, err := json.Marshal(userURLs)
	if err != nil {
		h.logger.Error(err.Error())
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *URLHandler) DeleteUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromCookie(r, h.secret)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var ids []string

	err = json.Unmarshal(body, &ids)
	if err != nil {
		h.logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	go func() {
		h.ToDelete <- models.BatchDeleteRequest{Ids: ids, UserID: userID}
	}()

	w.WriteHeader(http.StatusAccepted)
	if _, err := fmt.Fprintf(w, "%v", ids); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *URLHandler) ExpandURL(w http.ResponseWriter, r *http.Request) {
	ID := strings.TrimPrefix(r.URL.Path, "/")

	ctx := r.Context()
	URL, err := h.storage.Get(ctx, ID)
	if err != nil {
		if errors.Is(err, shared.ErrGone) {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusGone)
			return
		}

		if errors.Is(err, shared.ErrNotFound) {
			h.logger.Error(err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", URL)
	w.Header().Add("Content-Type", "text/html")

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
	r.Use(mw.LoggingMiddleware(handler.logger))
	r.Use(mw.GzipMiddleware(handler.logger))

	r.Post("/", handler.ShortenURL)
	r.Post("/api/shorten", handler.ShortenJSON)
	r.Post("/api/shorten/batch", handler.ShortenBatchJSON)
	r.Get("/api/user/urls", handler.AllUserURLs)
	r.Delete("/api/user/urls", handler.DeleteUserURLs)
	r.Get("/{id}", handler.ExpandURL)
	r.Get("/ping", handler.PingDB)

	return r
}
