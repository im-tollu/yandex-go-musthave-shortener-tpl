// Package apimodel provides structures for (un)marshalling JSON bodies
// of HTTP requests and responses.
package apimodel

// LongURLJson is a request model that a Client will use to send an original URL to be shortened
type LongURLJson struct {
	URL string `json:"url"`
}

// ShortURLJson is a response with shortened URL
type ShortURLJson struct {
	Result string `json:"result"`
}

type ShortURLForUserJSON struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type LongBatchURLJson struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"original_url"`
}

type ShortBatchURLJson struct {
	CorrelationID string `json:"correlation_id"`
	URL           string `json:"short_url"`
}
