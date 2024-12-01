package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/apetsko/shortugo/internal/storage"
	"github.com/apetsko/shortugo/internal/utils"
)

var s storage.Storage = storage.NewInMem()

func ShortenURL(w http.ResponseWriter, r *http.Request) {
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

	ID, err := s.Put(URL)
	shortenURL, err := utils.FullURL(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", shortenURL)
}

func ExpandURL(w http.ResponseWriter, r *http.Request) {
	ID := strings.TrimPrefix(r.URL.Path, "/")
	URL, err := s.Get(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Location", URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	fmt.Fprintf(w, "%s", URL)
}
