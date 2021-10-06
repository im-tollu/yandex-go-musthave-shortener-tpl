package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/api/handler"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/internal/app"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/backup"
	"github.com/im-tollu/yandex-go-musthave-shortener-tpl/storage/inmem"
)

type URLShortener struct {
	*chi.Mux
	Storage storage.Storage
	BaseURL url.URL
}

type URLShortenerServer struct {
	http.Server
	Storage     storage.BulkStorage
	StorageFile string
}

func (s *URLShortenerServer) ListenAndServe() error {
	if errRestore := backup.Restore(s.StorageFile, s.Storage); errRestore != nil {
		return fmt.Errorf("cannot restore URLs from storage file: %w", errRestore)
	}
	log.Printf("URL storage restored from [%s].", s.StorageFile)
	return s.Server.ListenAndServe()
}

func (s *URLShortenerServer) Shutdown(ctx context.Context) error {
	if errBackup := backup.Backup(s.StorageFile, s.Storage); errBackup != nil {
		return fmt.Errorf("cannot backup URLs to storage file: %w", errBackup)
	}
	log.Printf("URL storage backed up to [%s].", s.StorageFile)
	if errShutdown := s.Server.Shutdown(ctx); errShutdown != nil {
		return fmt.Errorf("cannot shutdown the server: %w", errShutdown)
	}
	return nil
}

func NewServer(conf app.Config) *URLShortenerServer {
	storage := inmem.New()
	return &URLShortenerServer{
		Server: http.Server{
			Addr:    conf.ServerAddress,
			Handler: handler.New(storage, conf.BaseURL),
		},
		Storage:     storage,
		StorageFile: conf.StorageFile,
	}
}
