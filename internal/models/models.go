package models

type URLRecord struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Result struct {
	Result string `json:"result"`
}

type BatchRequest struct {
	ID          string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type BatchResponse struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}
