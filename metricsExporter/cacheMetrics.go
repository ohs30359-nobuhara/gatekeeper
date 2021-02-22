package metricsExporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strings"
)

type CacheMetrics struct {
	exporter prometheus.CounterVec
}

func Init() *CacheMetrics {
	m := &CacheMetrics{
		exporter: *promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "request_cache_rate_metrics",
			Help: "cache hit rate per path",
		},[]string {
			"url", "hit",
		}),
	}
	return m
}

func (m CacheMetrics) Count(url string, hit string) {
	m.exporter.WithLabelValues(strings.Split(url, "?")[0], hit).Inc()
}

