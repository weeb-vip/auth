package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

// GzipMiddleware handles both request decompression and response compression
func GzipMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCompressed := false
			responseCompressed := false

			// Handle compressed request bodies
			if r.Header.Get("Content-Encoding") == "gzip" {
				gzipReader, err := gzip.NewReader(r.Body)
				if err != nil {
					log.Printf("Failed to decompress gzip request: %v", err)
					http.Error(w, "Invalid gzip data", http.StatusBadRequest)
					return
				}
				defer gzipReader.Close()
				r.Body = gzipReader
				requestCompressed = true
			}

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

			if requestCompressed || responseCompressed {
				log.Printf("Gzip compression applied - request_compressed: %t, response_compressed: %t", requestCompressed, responseCompressed)
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