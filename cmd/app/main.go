package main

import (
	"net/http"
	"path/filepath"

	"FileTP/internal/config"
	"FileTP/internal/handlers"
	"FileTP/internal/middlewares"
	"FileTP/internal/pkg/logging"
	"FileTP/internal/storage/sql"
)

func main() {
	log := logging.NewLogger(true)
	middleware := middlewares.NewMiddleware(log)

	cfg := config.NewConfig(filepath.Clean("./root/"))

	db, err := sql.NewFileDB("sqlite.db")
	if err != nil {
		log.Error(err.Error())
	}
	defer db.Close()

	handler := handlers.NewHandler(log, db, cfg)
	mux := http.NewServeMux()

	mux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	mux.Handle("/", middleware.MiddlewareLogging(http.HandlerFunc(handler.MainPage)))
	mux.Handle("/upload", middleware.MiddlewareLogging(http.HandlerFunc(handler.CreateFileHandler)))
	mux.Handle("/open", middleware.MiddlewareLogging(http.HandlerFunc(handler.OpenFileHandler)))

	log.Info("Server started -> http://localhost:1212")
	err = http.ListenAndServe(":1212", mux)
	if err != nil {
		panic(err)
	}
}
