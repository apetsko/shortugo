package models

type URLRecord struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	UserID  string `json:"userid"`
	Deleted bool   `json:"deleted"`
}

type Result struct {
	Result string `json:"result"`
}

type BatchDeleteRequest struct {
	Ids    []string
	UserID string
}

type BatchRequest struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type BatchResponse struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
