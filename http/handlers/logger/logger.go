package logger

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ctxKey struct{}

func Handler() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler { //nolint
		return addLogger(h)
	}
}
func addLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := logrus.WithFields(logrus.Fields{
			"user-agent": r.UserAgent(),
			"ip-address": r.RemoteAddr,
		})
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), &ctxKey{}, entry)))
	})
}

func FromContext(ctx context.Context) *logrus.Entry {
	entry := ctx.Value(&ctxKey{})

	// The ctxKey only gives this type.
	return entry.(*logrus.Entry) // nolint:forcetypeassert
}
