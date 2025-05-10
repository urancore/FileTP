package handlers

import (
	"FileTP/internal/config"
	"FileTP/internal/pkg/logging"
	"FileTP/internal/storage/sql"
)

type FTPHandler struct {
	Log* logging.Logger
	*sql.FileDB
	*config.Config
}

func NewHandler(log* logging.Logger, filedb *sql.FileDB, cfg *config.Config) *FTPHandler {
	return &FTPHandler{
		Log: log,
		FileDB: filedb,
		Config: cfg,
	}
}
