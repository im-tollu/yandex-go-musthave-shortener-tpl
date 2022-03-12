package shortener

import (
	"net/url"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
)

type URLService interface {
	// ShortenURL does all the necessary logic and stores produced short URL
	ShortenURL(u model.URLToShorten) (*model.ShortenedURL, error)

	// GetURLByID does a search of a shortened URL by ID
	GetURLByID(id int) (*model.ShortenedURL, error)

	// LookupURL does a search of a shortened URL by long URL;
	// This is used to get a duplicate item when shortening a URL
	// that is already shortened
	LookupURL(u url.URL) (*model.ShortenedURL, error)

	// GetUserURLs returns all URLs shortened by the userID
	GetUserURLs(userID int64) ([]model.ShortenedURL, error)

	// ScheduleDeletion adds url to deletion queue
	ScheduleDeletion(u model.URLToDelete)

	// AbsoluteURL resolves a short URL relative to base URL
	AbsoluteURL(u model.ShortenedURL) (*url.URL, error)
}
