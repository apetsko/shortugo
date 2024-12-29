package handlers

import (
	"io"
	"log"
	"net/http"
	"strings"

	mw "github.com/apetsko/shortugo/internal/middleware"
	"github.com/apetsko/shortugo/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Storage interface {
	Put(id string, url string) error
	Get(id string) (url string, err error)
}

type URLHandler struct {
	baseURL string
	storage Storage
}

func NewURLHandler(base string, storage Storage) *URLHandler {
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

	ID := utils.Generate(URL)

	err = h.storage.Put(ID, URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortenURL := utils.FullURL(h.baseURL, ID)

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
	_, err = w.Write([]byte(URL))
	if err != nil {
		log.Println(err.Error())
	}
}

func SetupRouter(handler *URLHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(mw.WithLogging)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Post("/", handler.ShortenURL)
	r.Get("/{id}", handler.ExpandURL)

	return r
}
