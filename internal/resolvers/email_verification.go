package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/graph/model"
	"github.com/weeb-vip/auth/http/handlers/logger"
	"github.com/weeb-vip/auth/http/handlers/requestinfo"
	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/mail"
	"github.com/weeb-vip/auth/internal/services/validation_token"
	"net/url"
)

type UserCreatedEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// EmailVerification is the legacy resolver returning only a success boolean.
// Deprecated: prefer EmailVerificationWithUser, which also returns the verified
// user's id so callers can link identity (e.g. analytics) on the verify page.
func EmailVerification( // nolint
	ctx context.Context,
	credentialService credential.Credential,
	userProducer func(ctx context.Context, value []byte) error,
) (bool, error) {
	_, ok, err := verifyEmailAndActivate(ctx, credentialService, userProducer)
	return ok, err
}

// EmailVerificationWithUser activates the caller's credentials and returns the
// verified user's id alongside the success flag.
func EmailVerificationWithUser( // nolint
	ctx context.Context,
	credentialService credential.Credential,
	userProducer func(ctx context.Context, value []byte) error,
) (*model.EmailVerificationResult, error) {
	userID, ok, err := verifyEmailAndActivate(ctx, credentialService, userProducer)
	result := &model.EmailVerificationResult{Success: ok}
	if userID != "" {
		uid := userID
		result.UserID = &uid
	}
	return result, err
}

// verifyEmailAndActivate holds the shared verification logic. It returns the
// verified user's id (empty when unavailable), whether activation succeeded,
// and any error. Semantics match the original EmailVerification exactly,
// including returning success=true even if the downstream event fails to
// produce.
func verifyEmailAndActivate(
	ctx context.Context,
	credentialService credential.Credential,
	userProducer func(ctx context.Context, value []byte) error,
) (string, bool, error) {
	log := logger.FromContext(ctx)
	req := requestinfo.FromContext(ctx)

	userID := req.UserID
	if userID == nil {
		log.Error("User ID is missing")
		return "", false, nil
	}

	err := credentialService.ActivateCredentials(ctx, *userID)
	if err != nil {
		_, err := handleError(ctx, "false", err)
		return "", false, err
	}

	credentials, err := credentialService.GetCredentialsByIdentifier(ctx, *userID)
	if err != nil {
		_, err := handleError(ctx, "false", err)
		return "", false, err
	}

	userCreatedEvent := UserCreatedEvent{
		UserID: credentials.UserID,
		Email:  credentials.Username,
	}

	payloadBytes, err := json.Marshal(userCreatedEvent)

	err = userProducer(ctx, payloadBytes)

	if err != nil {
		log.Errorf("Failed to produce user activated event: %v", err)
		return credentials.UserID, true, err
	}

	return credentials.UserID, true, nil
}

func ResendVerificationEmail( // nolint
	ctx context.Context,
	credentialService credential.Credential,
	validatonToken validation_token.ValidationToken,
	mailService mail.MailService,
	cfg *config.Config,
	username string,
) (bool, error) {
	log := logger.FromContext(ctx)

	credential, err := credentialService.GetCredentials(ctx, username)
	if err != nil {
		res, err := handleError(ctx, "false", err)
		if res != nil {
			return false, err
		}

		return false, err
	}

	if credential == nil {
		log.Error("Credential not found")
		return false, nil
	}

	if credential.Active {
		log.Error("Credential is already active")
		return false, nil
	}

	token, err := validatonToken.GenerateToken(credential.ID)
	if err != nil {
		return false, err
	}

	resetURL, err := url.Parse(cfg.APPConfig.VerificationBaseURL)
	if err != nil {
		return false, fmt.Errorf("invalid password reset base URL: %w", err)
	}

	query := resetURL.Query()
	query.Set("token", token)
	query.Set("email", username)
	resetURL.RawQuery = query.Encode()

	err = mailService.SendMail(ctx, []string{username}, "Email Verification", "verification.mjml", map[string]string{
		"token_url": resetURL.String(),
		"name":      username,
	})

	return true, nil
}
