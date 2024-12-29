package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	mw "github.com/apetsko/shortugo/internal/middleware"
	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Logger interface {
	Info(message string, keysAndValues ...interface{})
	Error(message string, keysAndValues ...interface{})
	Fatal(message string, keysAndValues ...interface{})
}

type Storage interface {
	Put(id string, url string) error
	Get(id string) (url string, err error)
}

type URLHandler struct {
	baseURL string
	storage Storage
	logger  Logger
}

func NewURLHandler(base string, storage Storage, logger Logger) *URLHandler {
	return &URLHandler{
		baseURL: base,
		storage: storage,
		logger:  logger,
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

	err = h.storage.Put(ID, URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortenURL := fmt.Sprintf("%s/%s", h.baseURL, ID)

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

	err = h.storage.Put(ID, req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp models.Response
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

func SetupRouter(handler *URLHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(mw.WithLogging(handler.logger))
	r.Use(mw.GzipMiddleware)

	r.Post("/", handler.ShortenURL)
	r.Post("/api/shorten", handler.ShortenJSON)
	r.Get("/{id}", handler.ExpandURL)

	return r
}
