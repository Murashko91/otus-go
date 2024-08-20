package internalhttp

import (
	"net/http"
	"time"

	"github.com/murashko91/otus-go/hw12_13_14_15_calendar/internal/app"
)

func loggingMiddleware(next http.HandlerFunc, log app.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next(w, r)
		log.Info(r.RemoteAddr, startTime.String(),
			r.Method, r.URL.Path, r.Proto, r.Response.StatusCode,
			time.Since(startTime), r.UserAgent())
	}
}
