package tracing

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	sdktrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var shutdownFunc func(context.Context) error

func InitTracing(ctx context.Context, serviceName string) (context.Context, error) {
	cfg := GetTracingConfig()

	exporter, err := createOTLPExporter(ctx, cfg)
	if err != nil {
		return ctx, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(cfg.ServiceVersion),
		attribute.String("environment", getEnvironment()),
	)

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	shutdownFunc = tracerProvider.Shutdown

	return ctx, nil
}

func createOTLPExporter(ctx context.Context, cfg TracingConfig) (trace.SpanExporter, error) {
	var opts []otlptracegrpc.Option

	opts = append(opts, otlptracegrpc.WithEndpoint(cfg.Endpoint))

	if cfg.Insecure {
		opts = append(opts,
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		)
	}

	return otlptracegrpc.New(ctx, opts...)
}

func GetTracer(ctx context.Context) sdktrace.Tracer {
	return otel.Tracer("auth-service")
}

func GetEnvironmentAttribute() sdktrace.SpanStartOption {
	return sdktrace.WithAttributes(
		attribute.String("environment", getEnvironment()),
	)
}

func getEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "development"
	}
	return env
}

func Shutdown(ctx context.Context) error {
	if shutdownFunc != nil {
		return shutdownFunc(ctx)
	}
	return nil
}