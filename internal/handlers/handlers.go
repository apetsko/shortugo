package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/apetsko/shortugo/internal/repositories"
	"github.com/apetsko/shortugo/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type URLHandler struct {
	baseURL string
	storage repositories.Storage
}

func NewURLHandler(base string, storage repositories.Storage) *URLHandler {
	return &URLHandler{
		baseURL: base,
		storage: storage,
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

	ID, err := h.storage.Put(URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortenURL, err := utils.FullURL(h.baseURL, ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenURL))
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
	w.Write([]byte(URL))
}

func SetupRouter(handler *URLHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/", handler.ShortenURL)
	r.Get("/{id}", handler.ExpandURL)

	return r
}
