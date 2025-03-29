package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/apetsko/shortugo/internal/storages/shared"
)

func (h *URLHandler) ExpandURL(w http.ResponseWriter, r *http.Request) {
	ID := strings.TrimPrefix(r.URL.Path, "/")

	ctx := r.Context()
	URL, err := h.storage.Get(ctx, ID)
	if err != nil {
		if errors.Is(err, shared.ErrGone) {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusGone)
			return
		}

		if errors.Is(err, shared.ErrNotFound) {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", URL)
	w.Header().Add("Content-Type", "text/html")

	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err = w.Write([]byte(URL))
	if err != nil {
		h.Logger.Error(err.Error())
	}
}
