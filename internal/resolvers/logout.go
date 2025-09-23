package resolvers

import (
	"context"
	"fmt"

	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/http/handlers/responsecontext"
)

func Logout(ctx context.Context, config *config.Config) (bool, error) {
	// Get response writer from context to set cookies
	responseWriter := responsecontext.FromContext(ctx)

	// Clear access token cookie by setting it to empty with immediate expiration
	accessTokenCookieStr := fmt.Sprintf(
		"access_token=; Path=/; Domain=%s; Max-Age=0; HttpOnly; Secure; SameSite=None",
		config.APPConfig.CookieDomain,
	)

	// Clear refresh token cookie by setting it to empty with immediate expiration
	refreshTokenCookieStr := fmt.Sprintf(
		"refresh_token=; Path=/; Domain=%s; Max-Age=0; HttpOnly; Secure; SameSite=None",
		config.APPConfig.CookieDomain,
	)

	// Set cookies manually to bypass Go's domain normalization
	responseWriter.Header().Add("Set-Cookie", accessTokenCookieStr)
	responseWriter.Header().Add("Set-Cookie", refreshTokenCookieStr)

	return true, nil
}