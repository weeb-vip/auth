package responsecontext

import (
	"context"
	"net/http"
)

type ctxKey struct{}

func FromContext(ctx context.Context) http.ResponseWriter {
	writer, found := ctx.Value(&ctxKey{}).(http.ResponseWriter)
	if !found {
		panic("response writer not set in context")
	}
	return writer
}

func WithResponseWriter(ctx context.Context, writer http.ResponseWriter) context.Context {
	return context.WithValue(ctx, &ctxKey{}, writer)
}