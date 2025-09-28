package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/http/handlers"
	"github.com/weeb-vip/auth/http/middleware"
	"github.com/weeb-vip/auth/internal/jwt"
	"github.com/weeb-vip/auth/internal/keypair"
	"github.com/weeb-vip/auth/internal/logger"
	"github.com/weeb-vip/auth/internal/metrics"
	observabilityMiddleware "github.com/weeb-vip/auth/internal/middleware"
	"github.com/weeb-vip/auth/internal/publishkey"
	"github.com/weeb-vip/auth/internal/tracing"

	"github.com/99designs/gqlgen/graphql/playground"
)

const minKeyValidityDurationMinutes = 5

func StartServer() error { // nolint
	return StartServerWithContext(context.Background())
}

func StartServerWithContext(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Initialize observability
	initObservability(ctx, cfg)

	rotatingKey, err := getRotatingSigningKey(cfg)
	if err != nil {
		return err
	}

	router := chi.NewRouter()

	// Add observability middleware
	router.Use(observabilityMiddleware.TracingMiddleware())

	// Add gzip compression middleware
	router.Use(middleware.GzipMiddleware())
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081", "http://localhost:3000"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	// Add metrics endpoint
	router.Handle("/metrics", metrics.NewPrometheusInstance().Handler())

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", handlers.BuildRootHandlerWithContext(ctx, jwt.New(rotatingKey)))
	router.Handle("/readyz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200) // nolint
	}))
	router.Handle("/livez", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200) // nolint
	}))

	log := logger.FromCtx(ctx)
	log.Info().
		Int("port", cfg.APPConfig.Port).
		Msg("Starting auth server")

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.APPConfig.Port), router) // nolint
}

func initObservability(ctx context.Context, cfg *config.Config) {
	environment := getEnvironment()

	// Initialize logger
	logger.Logger(
		logger.WithServerName("auth-service"),
		logger.WithVersion("1.0.0"),
		logger.WithEnvironment(environment),
	)

	// Initialize tracing
	tracedCtx, err := tracing.InitTracing(ctx, "auth-service")
	if err != nil {
		log := logger.FromCtx(ctx)
		log.Error().Err(err).Msg("Failed to initialize tracing")
	} else {
		ctx = tracedCtx
	}

	// Initialize metrics
	metrics.InitMetrics("auth-service", environment, "1.0.0")

	log := logger.FromCtx(ctx)
	log.Info().Msg("Observability initialized successfully")
}

func getEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "development"
	}
	return env
}

func getRotatingSigningKey(cfg *config.Config) (keypair.RotatingSigningKey, error) {
	rotatingKey, err := keypair.NewSigningKeyRotator(
		publishkey.NewKeyPublisher(
			cfg.APPConfig.InternalGraphQLURL).
			PublishToKeyManagementService)
	if err != nil {
		return nil, err
	}

	requestedDuration := time.Hour * time.Duration(cfg.APPConfig.KeyRollingDurationInHours)
	rotatingKey.RotateInBackground(getMinimumDuration(requestedDuration, time.Minute*minKeyValidityDurationMinutes))

	return rotatingKey, nil
}

func getMinimumDuration(askedDuration time.Duration, minimumDuration time.Duration) time.Duration {
	if askedDuration < minimumDuration {
		return minimumDuration
	}

	return askedDuration
}
