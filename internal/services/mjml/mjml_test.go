package mjml_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weeb-vip/auth/internal/services/mjml"
	"testing"
)

func TestMJML(t *testing.T) {
	t.Parallel()
	t.Run("should generate html from mjml", func(t *testing.T) {
		t.Parallel()
		a := assert.New(t)
		mjmlService := mjml.NewMJMLService()

		result, err := mjmlService.GenerateHTMLFromMJML(context.Background(), "verification.mjml", map[string]string{
			"token_url": "http://localhost:3000/verify?token=token",
			"name":      "John Doe",
		})

		a.NoError(err)
		a.NotEmpty(result)
	})
}
