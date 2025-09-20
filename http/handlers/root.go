package handlers

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/ThatCatDev/ep/v2/drivers"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/weeb-vip/auth/internal/services/mail"
	"github.com/weeb-vip/auth/internal/services/mjml"
	"go.uber.org/zap"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/sirupsen/logrus"

	epKafka "github.com/ThatCatDev/ep/v2/drivers/kafka"
	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/graph"
	"github.com/weeb-vip/auth/graph/generated"
	"github.com/weeb-vip/auth/http/handlers/logger"
	"github.com/weeb-vip/auth/http/handlers/metrics"
	"github.com/weeb-vip/auth/http/handlers/requestinfo"
	"github.com/weeb-vip/auth/http/handlers/responsecontext"
	"github.com/weeb-vip/auth/internal/jwt"
	logger2 "github.com/weeb-vip/auth/internal/logger"
	"github.com/weeb-vip/auth/internal/measurements"
	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/passwordreset"
	"github.com/weeb-vip/auth/internal/services/refresh_token"
	"github.com/weeb-vip/auth/internal/services/session"
	"github.com/weeb-vip/auth/internal/services/validation_token"
)

func BuildRootHandler(tokenizer jwt.Tokenizer) http.Handler { // nolint
	logrus.SetFormatter(&logrus.TextFormatter{})
	log := logger2.Get()

	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	kafkaConfig := &epKafka.KafkaConfig{
		ConsumerGroupName:        conf.KafkaConfig.ConsumerGroupName,
		BootstrapServers:         conf.KafkaConfig.BootstrapServers,
		SaslMechanism:            nil,
		SecurityProtocol:         nil,
		Username:                 nil,
		Password:                 nil,
		ConsumerSessionTimeoutMs: nil,
		ConsumerAutoOffsetReset:  &conf.KafkaConfig.Offset,
		ClientID:                 nil,
		Debug:                    nil,
	}

	driver := epKafka.NewKafkaDriver(kafkaConfig)
	defer func(driver drivers.Driver[*kafka.Message]) {
		err := driver.Close()
		if err != nil {
			log.Error("Error closing Kafka driver", zap.String("error", err.Error()))
		} else {
			log.Info("Kafka driver closed successfully")
		}
	}(driver)

	authenticationService := credential.NewCredentialService()
	passwordResetService := passwordreset.NewPasswordResetService()
	sessionService := session.NewSessionService()
	refreshTokenService := refresh_token.NewRefreshTokenService(conf.RefreshTokenConfig)
	validationTokenService := validation_token.NewValidationTokenService(tokenizer)
	mjmlService := mjml.NewMJMLService()
	mailService := mail.NewMailService(conf.EmailConfig, mjmlService)
	resolvers := &graph.Resolver{
		CredentialService:    authenticationService,
		PasswordResetService: passwordResetService,
		JwtTokenizer:         tokenizer,
		SessionService:       sessionService,
		Config:               *conf,
		RefreshTokenService:  refreshTokenService,
		ValidationToken:      validationTokenService,
		MailService:          mailService,
		UserProducer:         kafkaProducer(context.Background(), driver, conf.KafkaConfig.ProducerTopic),
	}
	cfg := generated.Config{Resolvers: resolvers}
	cfg.Directives.Authenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		req := requestinfo.FromContext(ctx)

		if req.UserID == nil {
			// unauthorized
			return nil, fmt.Errorf("Access denied")
		}

		return next(ctx)
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))
	srv.Use(apollotracing.Tracer{})

	client := measurements.New()

	// Create response context middleware
	responseContextHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := responsecontext.WithResponseWriter(r.Context(), w)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return requestinfo.Handler()(logger.Handler()(metrics.Handler(client)(responseContextHandler(srv))))
}

func kafkaProducer(ctx context.Context, driver drivers.Driver[*kafka.Message], topic string) func(ctx context.Context, message *kafka.Message) error {
	return func(ctx context.Context, message *kafka.Message) error {
		log := logger2.FromCtx(ctx)
		log.Info("Producing message to Kafka", zap.String("topic", topic), zap.String("key", string(message.Key)), zap.String("value", string(message.Value)))
		if err := driver.Produce(ctx, topic, message); err != nil {
			log.Error("Failed to produce message", zap.String("topic", topic), zap.Error(err))
			return err
		}
		return nil
	}
}
