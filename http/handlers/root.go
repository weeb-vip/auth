package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/sirupsen/logrus"

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
	observabilityMiddleware "github.com/weeb-vip/auth/internal/middleware"
	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/services/mail"
	"github.com/weeb-vip/auth/internal/services/mjml"
	"github.com/weeb-vip/auth/internal/services/passwordreset"
	"github.com/weeb-vip/auth/internal/services/refresh_token"
	"github.com/weeb-vip/auth/internal/services/session"
	"github.com/weeb-vip/auth/internal/services/validation_token"
)

func BuildRootHandler(tokenizer jwt.Tokenizer) http.Handler { // nolint
	return BuildRootHandlerWithContext(context.Background(), tokenizer)
}

func BuildRootHandlerWithContext(ctx context.Context, tokenizer jwt.Tokenizer) http.Handler { // nolint
	logrus.SetFormatter(&logrus.TextFormatter{})
	log := logger2.FromCtx(ctx)

	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// The publisher outlives this function on purpose.
	//
	// This builds the handler once at startup (see server.go), so the driver
	// has to survive for the life of the process. The previous code deferred
	// driver.Close() here, which closed it the moment the builder returned --
	// before a single request had been served. Closing is left to process exit.
	publish, _, err := newPublisher(*conf)
	if err != nil {
		panic(err)
	}

	// Worth a line at startup: which broker this writes to is otherwise
	// invisible until something downstream notices events are missing.
	log.Info().Str("producer_type", conf.ProducerType).Msg("Event producer configured")

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
		UserProducer:         publish,
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
	srv.Use(observabilityMiddleware.GraphQLTracingExtension{})

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
