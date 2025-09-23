package resolvers

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/http/handlers/responsecontext"
)

func TestLogout(t *testing.T) {
	tests := []struct {
		name         string
		cookieDomain string
		expectedAccessTokenCookie  string
		expectedRefreshTokenCookie string
	}{
		{
			name:         "logout with standard domain",
			cookieDomain: ".weeb.vip",
			expectedAccessTokenCookie:  "access_token=; Path=/; Domain=.weeb.vip; Max-Age=0; HttpOnly; Secure; SameSite=None",
			expectedRefreshTokenCookie: "refresh_token=; Path=/; Domain=.weeb.vip; Max-Age=0; HttpOnly; Secure; SameSite=None",
		},
		{
			name:         "logout with localhost domain",
			cookieDomain: "localhost",
			expectedAccessTokenCookie:  "access_token=; Path=/; Domain=localhost; Max-Age=0; HttpOnly; Secure; SameSite=None",
			expectedRefreshTokenCookie: "refresh_token=; Path=/; Domain=localhost; Max-Age=0; HttpOnly; Secure; SameSite=None",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test config
			cfg := &config.Config{
				APPConfig: config.AppConfig{
					CookieDomain: tt.cookieDomain,
				},
			}

			// Create test response writer
			recorder := httptest.NewRecorder()

			// Create context with response writer
			ctx := responsecontext.WithResponseWriter(context.Background(), recorder)

			// Call logout function
			result, err := Logout(ctx, cfg)

			// Assertions
			assert.NoError(t, err)
			assert.True(t, result)

			// Check that both cookies are set correctly
			cookies := recorder.Header()["Set-Cookie"]
			assert.Len(t, cookies, 2, "Expected 2 Set-Cookie headers")

			// Check access token cookie
			accessTokenFound := false
			refreshTokenFound := false

			for _, cookie := range cookies {
				if strings.Contains(cookie, "access_token=;") {
					assert.Equal(t, tt.expectedAccessTokenCookie, cookie)
					accessTokenFound = true
				}
				if strings.Contains(cookie, "refresh_token=;") {
					assert.Equal(t, tt.expectedRefreshTokenCookie, cookie)
					refreshTokenFound = true
				}
			}

			assert.True(t, accessTokenFound, "access_token cookie should be set")
			assert.True(t, refreshTokenFound, "refresh_token cookie should be set")
		})
	}
}