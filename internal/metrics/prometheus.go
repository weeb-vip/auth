package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusClient struct {
	registry   *prometheus.Registry
	histograms map[string]*prometheus.HistogramVec
	counters   map[string]*prometheus.CounterVec
	gauges     map[string]*prometheus.GaugeVec
	mu         sync.RWMutex
}

var (
	prometheusInstance *PrometheusClient
	once               sync.Once
)

func NewPrometheusInstance() *PrometheusClient {
	once.Do(func() {
		prometheusInstance = &PrometheusClient{
			registry:   prometheus.NewRegistry(),
			histograms: make(map[string]*prometheus.HistogramVec),
			counters:   make(map[string]*prometheus.CounterVec),
			gauges:     make(map[string]*prometheus.GaugeVec),
		}
	})
	return prometheusInstance
}

func (p *PrometheusClient) CreateHistogramVec(name, help string, labels []string, buckets []float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.histograms[name]; exists {
		return
	}

	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: buckets,
		},
		labels,
	)

	p.histograms[name] = histogram
	p.registry.MustRegister(histogram)
}

func (p *PrometheusClient) CreateCounterVec(name, help string, labels []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.counters[name]; exists {
		return
	}

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labels,
	)

	p.counters[name] = counter
	p.registry.MustRegister(counter)
}

func (p *PrometheusClient) CreateGaugeVec(name, help string, labels []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.gauges[name]; exists {
		return
	}

	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labels,
	)

	p.gauges[name] = gauge
	p.registry.MustRegister(gauge)
}

func (p *PrometheusClient) RecordHistogram(name string, value float64, labels prometheus.Labels) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if histogram, exists := p.histograms[name]; exists {
		histogram.With(labels).Observe(value)
	}
}

func (p *PrometheusClient) IncrementCounter(name string, labels prometheus.Labels) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if counter, exists := p.counters[name]; exists {
		counter.With(labels).Inc()
	}
}

func (p *PrometheusClient) SetGauge(name string, value float64, labels prometheus.Labels) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if gauge, exists := p.gauges[name]; exists {
		gauge.With(labels).Set(value)
	}
}

func (p *PrometheusClient) Handler() http.Handler {
	return promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{})
}