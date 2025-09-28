package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	Success = "success"
	Error   = "error"
)

type AppMetrics struct {
	prometheus  *PrometheusClient
	defaultTags map[string]string
}

var (
	appMetrics *AppMetrics
	metricsOnce sync.Once
)

func InitMetrics(serviceName, environment, version string) *AppMetrics {
	metricsOnce.Do(func() {
		prometheus := NewPrometheusInstance()
		initMetrics(prometheus)

		appMetrics = &AppMetrics{
			prometheus: prometheus,
			defaultTags: map[string]string{
				"service": serviceName,
				"env":     environment,
				"version": version,
			},
		}
	})
	return appMetrics
}

func GetAppMetrics() *AppMetrics {
	if appMetrics == nil {
		return InitMetrics("auth-service", "development", "1.0.0")
	}
	return appMetrics
}

func initMetrics(prometheusInstance *PrometheusClient) {
	// GraphQL resolver metrics
	prometheusInstance.CreateHistogramVec(
		"resolver_request_duration_histogram_milliseconds",
		"GraphQL resolver request duration in milliseconds",
		[]string{"service", "protocol", "resolver", "result", "env"},
		[]float64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 5000},
	)

	// HTTP request metrics
	prometheusInstance.CreateHistogramVec(
		"http_request_duration_histogram_milliseconds",
		"HTTP request duration in milliseconds",
		[]string{"service", "method", "path", "status_code", "env"},
		[]float64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 5000},
	)

	// Database metrics
	prometheusInstance.CreateHistogramVec(
		"database_query_duration_histogram_milliseconds",
		"Database query duration in milliseconds",
		[]string{"service", "table", "method", "result", "env"},
		[]float64{10, 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000},
	)

	// Authentication metrics
	prometheusInstance.CreateCounterVec(
		"auth_requests_total",
		"Total number of authentication requests",
		[]string{"service", "method", "result", "env"},
	)

	// Session metrics
	prometheusInstance.CreateCounterVec(
		"session_operations_total",
		"Total number of session operations",
		[]string{"service", "operation", "result", "env"},
	)
}

func (m *AppMetrics) ResolverMetric(duration float64, resolver string, result string) {
	labels := prometheus.Labels{
		"service":  m.defaultTags["service"],
		"protocol": "graphql",
		"resolver": resolver,
		"result":   result,
		"env":      m.defaultTags["env"],
	}
	m.prometheus.RecordHistogram("resolver_request_duration_histogram_milliseconds", duration, labels)
}

func (m *AppMetrics) HTTPRequestMetric(duration float64, method, path, statusCode string) {
	labels := prometheus.Labels{
		"service":     m.defaultTags["service"],
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"env":         m.defaultTags["env"],
	}
	m.prometheus.RecordHistogram("http_request_duration_histogram_milliseconds", duration, labels)
}

func (m *AppMetrics) DatabaseMetric(duration float64, table, method, result string) {
	labels := prometheus.Labels{
		"service": m.defaultTags["service"],
		"table":   table,
		"method":  method,
		"result":  result,
		"env":     m.defaultTags["env"],
	}
	m.prometheus.RecordHistogram("database_query_duration_histogram_milliseconds", duration, labels)
}

func (m *AppMetrics) AuthRequestMetric(method, result string) {
	labels := prometheus.Labels{
		"service": m.defaultTags["service"],
		"method":  method,
		"result":  result,
		"env":     m.defaultTags["env"],
	}
	m.prometheus.IncrementCounter("auth_requests_total", labels)
}

func (m *AppMetrics) SessionOperationMetric(operation, result string) {
	labels := prometheus.Labels{
		"service":   m.defaultTags["service"],
		"operation": operation,
		"result":    result,
		"env":       m.defaultTags["env"],
	}
	m.prometheus.IncrementCounter("session_operations_total", labels)
}