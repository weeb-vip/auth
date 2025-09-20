package responsecontext

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithResponseWriter(t *testing.T) {
	t.Run("should add response writer to context", func(t *testing.T) {
		ctx := context.Background()
		recorder := httptest.NewRecorder()

		newCtx := WithResponseWriter(ctx, recorder)

		writer := FromContext(newCtx)
		if writer != recorder {
			t.Errorf("Expected response writer to be %v, got %v", recorder, writer)
		}
	})
}

func TestFromContext(t *testing.T) {
	t.Run("should retrieve response writer from context", func(t *testing.T) {
		ctx := context.Background()
		recorder := httptest.NewRecorder()

		ctx = WithResponseWriter(ctx, recorder)
		writer := FromContext(ctx)

		if writer != recorder {
			t.Errorf("Expected response writer to be %v, got %v", recorder, writer)
		}
	})

	t.Run("should panic when response writer not in context", func(t *testing.T) {
		ctx := context.Background()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when response writer not in context")
			}
		}()

		FromContext(ctx)
	})
}

func TestResponseContextIntegration(t *testing.T) {
	t.Run("should work in HTTP handler middleware pattern", func(t *testing.T) {
		var capturedWriter http.ResponseWriter

		// Simulate middleware that adds response writer to context
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := WithResponseWriter(r.Context(), w)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}

		// Simulate handler that uses response writer from context
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedWriter = FromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		// Set up test
		wrappedHandler := middleware(handler)
		req := httptest.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		// Execute
		wrappedHandler.ServeHTTP(recorder, req)

		// Verify
		if capturedWriter != recorder {
			t.Errorf("Expected captured writer to be %v, got %v", recorder, capturedWriter)
		}
		if recorder.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", recorder.Code)
		}
	})
}