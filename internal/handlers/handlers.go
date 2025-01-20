package handlers

import (
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

var ErrNotFound = errors.New("not found")

type Storage interface {
	Put(id string, url string) error
	Get(id string) (url string, err error)
	Close() error
	Ping() error
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

	URL := string(body)
	if URL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	ID := utils.Generate(URL)

	shortenURL, err := h.storage.Get(ID)
	if err != nil {
		if !errors.Is(err, shared.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = h.storage.Put(ID, URL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		shortenURL = fmt.Sprintf("%s/%s", h.baseURL, ID)
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

	var req models.Request

	err = json.Unmarshal(body, &req)
	if err != nil {
		h.logger.Info("Error unmarshaling request body", "error", err.Error())
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "Empty URL", http.StatusBadRequest)
		return
	}

	ID := utils.Generate(req.URL)

	var resp models.Response

	if err = h.storage.Put(ID, req.URL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp.Result = fmt.Sprintf("%s/%s", h.baseURL, ID)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error(err.Error())
	}
}

func (h *URLHandler) ExpandURL(w http.ResponseWriter, r *http.Request) {
	ID := strings.TrimPrefix(r.URL.Path, "/")

	URL, err := h.storage.Get(ID)
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
	r.Get("/{id}", handler.ExpandURL)
	r.Get("/ping", handler.PingDB)

	return r
}
