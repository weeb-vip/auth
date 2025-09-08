package mail

import (
	"context"
	"fmt"
	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/internal/services/mjml"

	"github.com/wneessen/go-mail"
	"net/smtp"
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

	// Configure client based on SSL type and authentication requirements
	if cfg.SSLType == "none" {
		// For MailHog or other development SMTP servers without SSL/TLS
		if cfg.Username == "" && cfg.Password == "" {
			// No authentication required (e.g., MailHog)
			client, err = mail.NewClient(cfg.Host, 
				mail.WithPort(cfg.Port),
				mail.WithTLSPolicy(mail.NoTLS),
			)
		} else {
			// Plain authentication without TLS
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
	return &mailService{
		client: client,
		mjml:   mjmlService,
		config: cfg,
	}
}

func (s *mailService) SendMail(ctx context.Context, to []string, subject string, template string, values map[string]string) error {
	message := mail.NewMsg()
	if err := message.From(s.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set from email: %w", err)
	}

	if err := message.To(to...); err != nil {
		return fmt.Errorf("failed to set to email: %w", err)
	}

	message.Subject(subject)

	// Generate the email body using MJML
	body, err := s.mjml.GenerateHTMLFromMJML(ctx, template, values)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	message.SetBodyString(mail.TypeTextHTML, *body)

	if err := s.client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
