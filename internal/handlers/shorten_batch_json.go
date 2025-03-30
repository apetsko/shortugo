package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/apetsko/shortugo/internal/models"
	"github.com/apetsko/shortugo/internal/utils"
)

func (h *URLHandler) ShortenBatchJSON(w http.ResponseWriter, r *http.Request) {
	userID, err := h.auth.CookieGetUserID(r, h.secret)
	if err != nil {
		userID, err = h.auth.CookieSetUserID(w, h.secret)
		if err != nil {
			h.Logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.Logger.Error("error closing body", "error", err.Error())
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var reqs []models.BatchRequest

	if err = json.Unmarshal(body, &reqs); err != nil {
		h.Logger.Info("Error unmarshaling request body", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resps []models.BatchResponse
	var records []models.URLRecord

	for _, req := range reqs {
		var resp models.BatchResponse

		if req.OriginalURL == "" {
			errStr := fmt.Sprintf("%s: Empty URL", http.StatusText(http.StatusBadRequest))
			h.Logger.Error(errStr, "id", req.ID)
			resp.ID = req.ID
			resp.ShortURL = errStr
			resps = append(resps, resp)
			continue
		}

		IDlen := 8
		var record = models.URLRecord{
			URL:    req.OriginalURL,
			ID:     utils.GenerateID(req.OriginalURL, IDlen),
			UserID: userID,
		}

		records = append(records, record)

		shortURL := fmt.Sprintf("%s/%s", h.baseURL, record.ID)
		resp = models.BatchResponse{ID: req.ID, ShortURL: shortURL}
		resps = append(resps, resp)
	}

	ctx := r.Context()
	if err = h.storage.PutBatch(ctx, records); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resps); err != nil {
		h.Logger.Error(err.Error())
	}
}
