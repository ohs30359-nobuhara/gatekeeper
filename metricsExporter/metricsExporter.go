package metricsExporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strings"
)

type MetricsExporter struct {
	cacheExporter prometheus.CounterVec
	rateLimitExporter prometheus.CounterVec
}

func Init() *MetricsExporter {
	m := &MetricsExporter{
		cacheExporter: *promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "request_cache_rate_metrics",
			Help: "cache hit rate per path",
		},[]string {
			"url", "hit",
		}),
		rateLimitExporter: *promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limit_metrics",
			Help: "request block rate per path",
		},[]string {
			"url",
		}),
	}
	return m
}

func (m MetricsExporter) Cache(url string, hit string) {
	m.cacheExporter.WithLabelValues(strings.Split(url, "?")[0], hit).Inc()
}

func (m MetricsExporter) RateLimit(url string) {
	m.rateLimitExporter.WithLabelValues(strings.Split(url, "?")[0]).Inc()
}
