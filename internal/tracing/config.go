package tracing

import (
	"os"
)

type TracingConfig struct {
	Endpoint       string
	Insecure       bool
	ServiceVersion string
}

func GetTracingConfig() TracingConfig {
	config := TracingConfig{
		Endpoint:       os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		Insecure:       true,
		ServiceVersion: os.Getenv("SERVICE_VERSION"),
	}

	if config.Endpoint == "" {
		config.Endpoint = "localhost:4317"
	}

	if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "false" {
		config.Insecure = false
	}

	if config.ServiceVersion == "" {
		config.ServiceVersion = "1.0.0"
	}

	return config
}