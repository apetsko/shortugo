package storages

import (
	"github.com/apetsko/shortugo/internal/handlers"
	"github.com/apetsko/shortugo/internal/logging"
	"github.com/apetsko/shortugo/internal/storages/infile"
	"github.com/apetsko/shortugo/internal/storages/inmem"
	"github.com/apetsko/shortugo/internal/storages/postgres"
)

func Init(databaseDSN, fileStoragePath string, logger *logging.ZapLogger) (handlers.Storage, error) {
	switch {
	case databaseDSN != "":
		s, err := postgres.New(databaseDSN)
		if err != nil {
			return nil, err
		}
		logger.Info("Using database storages")
		return s, nil
	case fileStoragePath != "":
		s, err := infile.New(fileStoragePath)
		if err != nil {
			return nil, err
		}
		logger.Info("Using file storages")
		return s, nil
	default:
		s := inmem.New()
		logger.Info("Using in-memory storages")
		return s, nil
	}
}
