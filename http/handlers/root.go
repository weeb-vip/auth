package handlers

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/sirupsen/logrus"

	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/graph"
	"github.com/weeb-vip/auth/graph/generated"
	"github.com/weeb-vip/auth/http/handlers/logger"
	"github.com/weeb-vip/auth/http/handlers/metrics"
	"github.com/weeb-vip/auth/http/handlers/requestinfo"
	"github.com/weeb-vip/auth/internal/jwt"
	"github.com/weeb-vip/auth/internal/measurements"
	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/passwordreset"
	"github.com/weeb-vip/auth/internal/services/refresh_token"
	"github.com/weeb-vip/auth/internal/services/session"
	"github.com/weeb-vip/auth/internal/services/validation_token"
)

func BuildRootHandler(tokenizer jwt.Tokenizer) http.Handler { // nolint
	logrus.SetFormatter(&logrus.TextFormatter{})

	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	authenticationService := credential.NewCredentialService()
	passwordResetService := passwordreset.NewPasswordResetService()
	sessionService := session.NewSessionService()
	refreshTokenService := refresh_token.NewRefreshTokenService(conf.RefreshTokenConfig)
	validationTokenService := validation_token.NewValidationTokenService(tokenizer)
	resolvers := &graph.Resolver{
		CredentialService:    authenticationService,
		PasswordResetService: passwordResetService,
		JwtTokenizer:         tokenizer,
		SessionService:       sessionService,
		Config:               *conf,
		RefreshTokenService:  refreshTokenService,
		ValidationToken:      validationTokenService,
	}
	cfg := generated.Config{Resolvers: resolvers}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))
	srv.Use(apollotracing.Tracer{})

	client := measurements.New()

	return requestinfo.Handler()(logger.Handler()(metrics.Handler(client)(srv)))
}
