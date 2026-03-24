package agents

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	requestsTotal *prometheus.CounterVec
	latencySec    *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return NewMetricsWithRegisterer(prometheus.DefaultRegisterer)
}

func NewMetricsWithRegisterer(registerer prometheus.Registerer) *Metrics {
	factory := promauto.With(registerer)

	return &Metrics{
		requestsTotal: factory.NewCounterVec(prometheus.CounterOpts{
			Namespace: "pingplex",
			Subsystem: "agents",
			Name:      "requests_total",
			Help:      "Total number of agents service requests partitioned by method and status",
		}, []string{"method", "status"}),
		latencySec: factory.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "pingplex",
			Subsystem: "agents",
			Name:      "request_duration_seconds",
			Help:      "Duration of agents service requests partitioned by method",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method"}),
	}
}

func (m *Metrics) Observe(method string, startedAt time.Time, err error) {
	status := "ok"
	if err != nil {
		status = "error"
	}

	m.requestsTotal.WithLabelValues(method, status).Inc()
	m.latencySec.WithLabelValues(method).Observe(time.Since(startedAt).Seconds())
}
