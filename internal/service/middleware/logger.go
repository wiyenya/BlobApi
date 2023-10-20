package middlewares

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
)

func Logger(entry *logan.Entry, durationThreshold time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()

			// Set a header for each response
			w.Header().Set("Content-Type", "application/json")

			defer func() {
				dur := time.Since(t1)
				lEntry := entry.WithFields(logan.F{
					"path":     r.URL.Path,
					"duration": dur,
				})
				lEntry.Info("request finished")

				if dur > durationThreshold {
					lEntry.WithField("http_request", r).Warn("slow request")
				}
			}()

			entry.WithField("path", r.URL.Path).Info("request started")
			next.ServeHTTP(w, r)
		})
	}
}
