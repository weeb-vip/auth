package logger

import (
	"context"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type Option func(*Config)

type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
}

var (
	globalLogger zerolog.Logger
	once         sync.Once
)

func WithServerName(name string) Option {
	return func(c *Config) {
		c.ServiceName = name
	}
}

func WithVersion(version string) Option {
	return func(c *Config) {
		c.ServiceVersion = version
	}
}

func WithEnvironment(env string) Option {
	return func(c *Config) {
		c.Environment = env
	}
}

func Logger(opts ...Option) {
	once.Do(func() {
		config := &Config{
			ServiceName:    "auth-service",
			ServiceVersion: "1.0.0",
			Environment:    "development",
		}

		for _, opt := range opts {
			opt(config)
		}

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		var logger zerolog.Logger
		if config.Environment == "development" {
			output := zerolog.ConsoleWriter{Out: os.Stdout}
			logger = zerolog.New(output)
		} else {
			logger = zerolog.New(os.Stdout)
		}

		globalLogger = logger.
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Str("service", config.ServiceName).
			Str("version", config.ServiceVersion).
			Str("environment", config.Environment).
			Logger()

		log.Logger = globalLogger
	})
}

func FromCtx(ctx context.Context) zerolog.Logger {
	return withTraceContext(ctx, globalLogger)
}

func Global() zerolog.Logger {
	return globalLogger
}

func withTraceContext(ctx context.Context, logger zerolog.Logger) zerolog.Logger {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return logger
	}

	spanContext := span.SpanContext()
	return logger.With().
		Str("trace_id", spanContext.TraceID().String()).
		Str("span_id", spanContext.SpanID().String()).
		Logger()
}
