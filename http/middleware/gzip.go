package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

// GzipMiddleware handles response compression only
func GzipMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseCompressed := false

			// Check if client accepts gzip encoding for response
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			// Wrap response writer with gzip compression
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")
			responseCompressed = true

			gzipWriter := gzip.NewWriter(w)
			defer gzipWriter.Close()

			gzipResponseWriter := &gzipResponseWriter{
				ResponseWriter: w,
				Writer:         gzipWriter,
			}

			if responseCompressed {
				log.Printf("Gzip response compression applied")
			}

			next.ServeHTTP(gzipResponseWriter, r)
		})
	}
}

// gzipResponseWriter wraps http.ResponseWriter to provide gzip compression
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	if flusher, ok := w.Writer.(*gzip.Writer); ok {
		flusher.Flush()
	}
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}