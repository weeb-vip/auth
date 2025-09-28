package middleware

import (
	"context"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/weeb-vip/auth/internal/logger"
	"github.com/weeb-vip/auth/internal/metrics"
	"github.com/weeb-vip/auth/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type GraphQLTracingExtension struct{}

func (e GraphQLTracingExtension) ExtensionName() string {
	return "GraphQLTracing"
}

func (e GraphQLTracingExtension) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (e GraphQLTracingExtension) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	rc := graphql.GetOperationContext(ctx)
	startTime := time.Now()

	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "GraphQL "+string(rc.Operation.Operation),
		trace.WithAttributes(
			attribute.String("graphql.operation.name", rc.OperationName),
			attribute.String("graphql.operation.type", string(rc.Operation.Operation)),
			attribute.String("graphql.document", rc.RawQuery),
		),
		trace.WithSpanKind(trace.SpanKindServer),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	log := logger.FromCtx(ctx)
	log.Info().
		Str("operation_name", rc.OperationName).
		Str("operation_type", string(rc.Operation.Operation)).
		Msg("GraphQL operation started")

	responseHandler := next(ctx)

	return func(ctx context.Context) *graphql.Response {
		response := responseHandler(ctx)
		duration := time.Since(startTime)

		result := metrics.Success
		if response.Errors != nil && len(response.Errors) > 0 {
			result = metrics.Error
			for _, err := range response.Errors {
				span.RecordError(err)
				log.Error().
					Err(err).
					Str("operation_name", rc.OperationName).
					Msg("GraphQL operation error")
			}
		}

		log.Info().
			Str("operation_name", rc.OperationName).
			Str("operation_type", string(rc.Operation.Operation)).
			Dur("duration", duration).
			Str("result", result).
			Msg("GraphQL operation completed")

		// Record metrics
		metrics.GetAppMetrics().ResolverMetric(
			float64(duration.Milliseconds()),
			rc.OperationName,
			result,
		)

		return response
	}
}

func (e GraphQLTracingExtension) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	startTime := time.Now()

	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "GraphQL Field: "+fc.Field.Name,
		trace.WithAttributes(
			attribute.String("graphql.field.name", fc.Field.Name),
			attribute.String("graphql.field.path", fc.Path().String()),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	result, err := next(ctx)
	duration := time.Since(startTime)

	if err != nil {
		span.RecordError(err)
		log := logger.FromCtx(ctx)
		log.Error().
			Err(err).
			Str("field_name", fc.Field.Name).
			Str("field_path", fc.Path().String()).
			Dur("duration", duration).
			Msg("GraphQL field error")
	}

	return result, err
}

func (e GraphQLTracingExtension) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	return next(ctx)
}