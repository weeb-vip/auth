package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/weeb-vip/auth/internal/logger"
	"github.com/weeb-vip/auth/internal/metrics"
	"github.com/weeb-vip/auth/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func TracingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			ctx := r.Context()

			// Extract existing trace context from headers
			propagator := propagation.TraceContext{}
			ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

			// Start new span
			tracer := tracing.GetTracer(ctx)
			ctx, span := tracer.Start(ctx, "HTTP "+r.Method+" "+r.URL.Path,
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.String()),
					attribute.String("http.host", r.Host),
					attribute.String("user_agent.original", r.UserAgent()),
				),
				trace.WithSpanKind(trace.SpanKindServer),
				tracing.GetEnvironmentAttribute(),
			)
			defer span.End()

			// Inject trace context to response headers
			propagator.Inject(ctx, propagation.HeaderCarrier(w.Header()))

			// Wrap response to capture status code
			wrapped := &responseWrapper{ResponseWriter: w, statusCode: 200}
			r = r.WithContext(ctx)

			// Log request start
			log := logger.FromCtx(ctx)
			log.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Msg("HTTP request started")

			next.ServeHTTP(wrapped, r)

			duration := time.Since(startTime)
			statusCode := strconv.Itoa(wrapped.statusCode)

			// Set span attributes
			span.SetAttributes(attribute.Int("http.status_code", wrapped.statusCode))

			// Record error if status code indicates failure
			if wrapped.statusCode >= 400 {
				span.RecordError(nil, trace.WithAttributes(
					attribute.String("error.type", "http_error"),
					attribute.Int("http.status_code", wrapped.statusCode),
				))
			}

			// Log request completion
			log.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status_code", wrapped.statusCode).
				Dur("duration", duration).
				Msg("HTTP request completed")

			// Record metrics
			metrics.GetAppMetrics().HTTPRequestMetric(
				float64(duration.Milliseconds()),
				r.Method,
				r.URL.Path,
				statusCode,
			)
		})
	}
}