package config

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/jinzhu/configor"
)

type Config struct {
	APPConfig          AppConfig
	DBConfig           DBConfig
	RefreshTokenConfig RefreshTokenConfig
	UserClient         UserClientConfig
	EmailConfig        EmailConfig
	KafkaConfig        KafkaConfig
	NatsConfig         NatsConfig
	ProducerType       string `default:"kafka" env:"PRODUCER_TYPE"`
}

type AppConfig struct {
	Name                      string `env:"CONFIG__APP_CONFIG__NAME" required:"true" default:"card-delivery-service"`
	Version                   string `env:"APP__VERSION" default:"local"`
	Port                      int    `env:"CONFIG__APP_CONFIG__PORT" default:"3001"`
	KeyRollingDurationInHours int    `env:"CONFIG__APP_CONFIG__KEY_ROLLING_DURATION_IN_HOURS" default:"1"`
	InternalGraphQLURL        string `env:"INTERNAL_GRAPHQL_URL" default:"http://localhost:5001/graphql"`
	JWTValiditySeconds        int    `env:"CONFIG__APP_CONFIG__JWT_VALIDITY_SECONDS" default:"900"` // 15 minutes.
	PasswordResetBaseURL      string `env:"CONFIG__APP_CONFIG__PASSWORD_RESET_BASE_URL" default:"http://localhost:3000/auth/password-reset"`
	VerificationBaseURL       string `env:"CONFIG__APP_CONFIG__VERIFICATION_BASE_URL" default:"http://localhost:3000/auth/verify"`
	CookieDomain              string `env:"CONFIG__APP_CONFIG__COOKIE_DOMAIN" default:".weeb.vip"`
}

type DBConfig struct {
	Host               string `env:"DBHOST" required:"true" default:"localhost"`
	Port               uint   `env:"DBPORT" required:"true" default:"5432"`
	User               string `env:"DBUSER" required:"true" default:"postgres"`
	Password           string `env:"DBPASSWORD" required:"true" default:"mysecretpassword"`
	DB                 string `env:"DBNAME" required:"true" default:"auth"`
	SSL                string `env:"DBSSL" default:"require"`
	MigrationTableName string `env:"DBMIGRATIONTABLE" default:"__migrations_auth"`
}

type RefreshTokenConfig struct {
	TokenTTL int `env:"CONFIG__REFRESH_TOKEN_CONFIG__TOKEN_TTL" default:"4380"` // 6 months in hours.
}

type UserClientConfig struct {
	URL string `env:"CONFIG__USER_CLIENT__URL" required:"true"`
}

type EmailConfig struct {
	FromEmail string `env:"CONFIG__EMAIL_FROM_EMAIL" required:"true"`
	FromName  string `env:"CONFIG__EMAIL_FROM_NAME" required:"true"`
	Username  string `env:"CONFIG__EMAIL_USERNAME" required:"false"`
	Password  string `env:"CONFIG__EMAIL_PASSWORD" required:"false"`
	Host      string `env:"CONFIG__EMAIL_HOST" required:"true"`
	Port      int    `env:"CONFIG__EMAIL_PORT" required:"true"`
	SSLType   string `env:"CONFIG__EMAIL_SSL_TYPE" required:"true"`
}

type KafkaConfig struct {
	ConsumerGroupName string `default:"auth-group" env:"KAFKA_CONSUMER_GROUP_NAME"`
	BootstrapServers  string `default:"localhost:9092" env:"KAFKA_BOOTSTRAP_SERVERS"`
	Offset            string `default:"earliest" env:"KAFKA_OFFSET"`
	ProducerTopic     string `default:"user_created" env:"KAFKA_PRODUCER_TOPIC"`
}

// NatsConfig mirrors the producer half of KafkaConfig.
//
// Which of the two is used is decided by ProducerType at startup rather than
// by a separate binary or command: this producer lives inside the GraphQL
// handler, so there is no command name to select a transport with the way the
// consumers have.
type NatsConfig struct {
	URL string `default:"nats://localhost:4222" env:"NATSURL"`

	// Empty StreamName: user-created is produced here, not by Debezium, so the
	// driver creates a stream from the subject rather than binding to the CDC
	// one.
	StreamName string `env:"NATSSTREAMNAME"`

	ProducerSubject string `default:"user_created" env:"NATSPRODUCERSUBJECT"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := configor.
		New(&configor.Config{AutoReload: false}).
		Load(&config, fmt.Sprintf("%s/config.%s.json", getConfigLocation(), getEnv()))

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getConfigLocation() string {
	_, filename, _, _ := runtime.Caller(0) // nolint

	return path.Join(path.Dir(filename), "../config")
}

func getEnv() string {
	prod := "prod"
	dev := "dev"
	docker := "docker"

	val := os.Getenv("APP_ENV")
	switch val {
	case prod:
		return prod
	case docker:
		return docker
	default:
		return dev
	}
}
