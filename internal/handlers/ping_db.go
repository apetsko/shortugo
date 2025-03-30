package handlers

import "net/http"

func (h *URLHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
