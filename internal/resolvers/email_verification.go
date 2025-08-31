package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/weeb-vip/auth/config"
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

func EmailVerification( // nolint
	ctx context.Context,
	credentialService credential.Credential,
	userProducer func(ctx context.Context, message *kafka.Message) error,
) (bool, error) {
	log := logger.FromContext(ctx)
	req := requestinfo.FromContext(ctx)

	userID := req.UserID
	if userID == nil {
		log.Error("User ID is missing")
		return false, nil
	}

	err := credentialService.ActivateCredentials(ctx, *userID)
	if err != nil {
		res, err := handleError(ctx, "false", err)
		if res != nil {
			return false, err
		}

		return false, err
	}

	credential, err := credentialService.GetCredentials(ctx, *userID)
	if err != nil {
		res, err := handleError(ctx, "false", err)
		if res != nil {
			return false, err
		}

		return false, err
	}

	userCreatedEvent := UserCreatedEvent{
		UserID: *userID,
		Email:  credential.Username,
	}

	payloadBytes, err := json.Marshal(userCreatedEvent)

	err = userProducer(ctx, &kafka.Message{
		Value: payloadBytes,
	})

	if err != nil {
		log.Errorf("Failed to produce user activated event: %v", err)
		return true, err
	}

	return true, nil
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
