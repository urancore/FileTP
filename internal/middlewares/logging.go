package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"FileTP/internal/pkg/logging"
	"FileTP/internal/storage/sql"
)

type FTPMiddleware struct {
	Log *logging.Logger
	DB *sql.FileDB
}
func NewMiddleware(log *logging.Logger, db *sql.FileDB) *FTPMiddleware {
	return &FTPMiddleware{
		Log: log,
		DB: db,
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code before writing headers
func (r *responseRecorder) WriteHeader(status int) {
	if r.status == 0 { // Only set if not already set
		r.status = status
	}
	r.ResponseWriter.WriteHeader(status)
}

// Write captures writes to the response body and ensures status is set
func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 { // Set default status if not set
		r.status = http.StatusOK
	}
	return r.ResponseWriter.Write(b)
}


func (m *FTPMiddleware) MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w}


		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		status := recorder.status

		// Handle cases where no response was written
		if status == 0 {
			status = http.StatusOK
		}

		// Log request details using the provided logger
		msg := fmt.Sprintf("request handled - method=%s path=%s remote=%s status=%d duration=%s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			status,
			duration,)
		m.Log.Info(msg)
	})
}
