package mail

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/internal/logger"
	"github.com/weeb-vip/auth/internal/services/mjml"

	"github.com/wneessen/go-mail"
)

type MailService interface {
	SendMail(ctx context.Context, to []string, subject string, template string, values map[string]string) error
}

type mailService struct {
	auth   smtp.Auth
	mjml   mjml.MJMLService
	config config.EmailConfig
	client *mail.Client
}

// NewMailService creates a new instance of MailService
func NewMailService(cfg config.EmailConfig, mjmlService mjml.MJMLService) MailService {
	var client *mail.Client
	var err error

	log := logger.Global()
	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("ssl_type", cfg.SSLType).
		Bool("has_auth", cfg.Username != "").
		Msg("Initializing mail service")

	// Configure client based on SSL type and authentication requirements
	if cfg.SSLType == "none" {
		// For MailHog or other development SMTP servers without SSL/TLS
		if cfg.Username == "" && cfg.Password == "" {
			// No authentication required (e.g., MailHog)
			log.Info().Msg("Creating mail client without authentication (SSLType=none)")
			client, err = mail.NewClient(cfg.Host,
				mail.WithPort(cfg.Port),
				mail.WithTLSPolicy(mail.NoTLS),
			)
		} else {
			// Plain authentication without TLS
			log.Info().Msg("Creating mail client with plain auth (SSLType=none)")
			client, err = mail.NewClient(cfg.Host,
				mail.WithPort(cfg.Port),
				mail.WithSMTPAuth(mail.SMTPAuthPlain),
				mail.WithUsername(cfg.Username),
				mail.WithPassword(cfg.Password),
				mail.WithTLSPolicy(mail.NoTLS),
			)
		}
	} else {
		// Standard configuration with TLS/SSL
		log.Info().Str("ssl_type", cfg.SSLType).Msg("Creating mail client with TLS")
		client, err = mail.NewClient(cfg.Host,
			mail.WithPort(cfg.Port),
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(cfg.Username),
			mail.WithPassword(cfg.Password),
		)
	}

	if err != nil {
		panic(fmt.Errorf("failed to create mail client: %w", err))
	}
	log.Info().Msg("Mail client created successfully")
	return &mailService{
		client: client,
		mjml:   mjmlService,
		config: cfg,
	}
}

func (s *mailService) SendMail(ctx context.Context, to []string, subject string, template string, values map[string]string) error {
	log := logger.FromCtx(ctx)
	log.Info().
		Strs("to", to).
		Str("subject", subject).
		Str("template", template).
		Msg("SendMail started")

	message := mail.NewMsg()
	if err := message.From(s.config.FromEmail); err != nil {
		log.Error().Err(err).Msg("Failed to set from email")
		return fmt.Errorf("failed to set from email: %w", err)
	}

	if err := message.To(to...); err != nil {
		log.Error().Err(err).Msg("Failed to set to email")
		return fmt.Errorf("failed to set to email: %w", err)
	}

	message.Subject(subject)

	// Generate the email body using MJML
	log.Info().Str("template", template).Msg("Generating HTML from MJML template")
	body, err := s.mjml.GenerateHTMLFromMJML(ctx, template, values)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate email body from MJML")
		return fmt.Errorf("failed to generate email body: %w", err)
	}
	log.Info().Int("body_length", len(*body)).Msg("MJML generated successfully")

	message.SetBodyString(mail.TypeTextHTML, *body)

	log.Info().
		Str("host", s.config.Host).
		Int("port", s.config.Port).
		Msg("Dialing SMTP server")
	if err := s.client.DialAndSend(message); err != nil {
		log.Error().Err(err).Msg("Failed to send email via SMTP")
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().Strs("to", to).Msg("Email sent successfully")
	return nil
}
