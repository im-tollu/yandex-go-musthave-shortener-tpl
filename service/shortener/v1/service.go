package v1

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/model"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
)

type Service struct {
	Storage    storage.ShortenerStorage
	BaseURL    url.URL
	deleteChan chan model.URLToDelete
}

func New(s storage.ShortenerStorage, u url.URL) (*Service, error) {
	if s == nil {
		return nil, errors.New("storage should not be nil")
	}
	deleteChan := make(chan model.URLToDelete, 1)

	service := Service{s, u, deleteChan}
	go service.processDeleteURLs()

	return &service, nil
}

func (s *Service) ShortenURL(u model.URLToShorten) (*model.ShortenedURL, error) {
	shortenedURL, err := s.Storage.SaveURL(u)
	if err != nil {
		return nil, fmt.Errorf("cannot shorten url: %w", err)
	}
	log.Printf("Shortened: %s", shortenedURL)

	return &shortenedURL, nil
}

func (s *Service) GetURLByID(id int) (*model.ShortenedURL, error) {
	return s.Storage.GetURLByID(id)
}

func (s *Service) LookupURL(u url.URL) (*model.ShortenedURL, error) {
	return s.Storage.LookupURL(u)
}

func (s *Service) GetUserURLs(userID int64) ([]model.ShortenedURL, error) {
	return s.Storage.ListByUserID(userID)
}

func (s *Service) ScheduleDeletion(u model.URLToDelete) {
	s.deleteChan <- u
}

func (s *Service) AbsoluteURL(u model.ShortenedURL) (*url.URL, error) {
	urlPath := fmt.Sprintf("%d", u.ID)

	shortURL, err := s.BaseURL.Parse(urlPath)
	if err != nil {
		return nil, fmt.Errorf("cannot make absolute URL for id [%d]", u.ID)
	}

	return shortURL, nil
}

func (s *Service) processDeleteURLs() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	batch := makeEmptyBatch()

	for {
		select {
		case <-ticker.C:
			if len(batch) > 0 {
				s.Storage.DeleteBatchURLs(batch)
				batch = makeEmptyBatch()
			}
		case u, ok := <-s.deleteChan:
			if !ok {
				s.Storage.DeleteBatchURLs(batch)
				return
			}

			batch = append(batch, u)
			if len(batch) == cap(batch) {
				s.Storage.DeleteBatchURLs(batch)
				batch = makeEmptyBatch()
			}
		}
	}
}

func makeEmptyBatch() []model.URLToDelete {
	return make([]model.URLToDelete, 0, 1000)
}
