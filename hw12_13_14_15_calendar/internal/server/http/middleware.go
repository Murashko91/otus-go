package internalhttp

import (
	"net/http"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (handler *statusWriter) WriteHeader(code int) {
	handler.status = code
	handler.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler, log app.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		httpWriter := statusWriter{w, 200}
		next.ServeHTTP(&httpWriter, r)
		log.Info(r.RemoteAddr, startTime.String(),
			r.Method, r.URL.Path, r.Proto, httpWriter.status,
			time.Since(startTime), r.UserAgent())
	})
}
