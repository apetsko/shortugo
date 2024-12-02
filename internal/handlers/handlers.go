package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/apetsko/shortugo/internal/repositories"
	"github.com/apetsko/shortugo/internal/utils"
)

type URLHandler struct {
	Storage repositories.Storage
}

func NewURLHandler(storage repositories.Storage) *URLHandler {
	return &URLHandler{Storage: storage}
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

	ID, err := h.Storage.Put(URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortenURL, err := utils.FullURL(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", shortenURL)
}

func (h *URLHandler) ExpandURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ID := ctx.Value("id").(string)
	if len(ID) == 0 {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	URL, err := h.Storage.Get(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Fprintf(w, "%s", URL)
}
