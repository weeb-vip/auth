package session_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/weeb-vip/auth/internal/services/session"
)

func TestNewSessionService(t *testing.T) {
	t.Parallel()
	t.Run("should return a new session service", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)

		sessionService := session.NewSessionService()

		a.NotNil(sessionService)
	})

	t.Run("should create a new session", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)

		sessionService := session.NewSessionService()

		_, err := sessionService.CreateSession(context.TODO(), "username")
		a.NoError(err)
	})
}
