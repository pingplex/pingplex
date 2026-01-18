package users

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	MetricsNamespace = "coordinator"
	MetricsSubsystem = "users"
)

// Metrics handles Prometheus metrics collection for the users module.
type Metrics struct {
	// totalUsers is a Prometheus counter that tracks the total number
	// of users created in the system.
	totalUsers prometheus.Counter

	// activeUsers is a Prometheus gauge that tracks the current number
	// of active users.
	activeUsers prometheus.Gauge

	// createUserDuration is a Prometheus histogram that tracks the
	// duration of user creation operations.
	createUserDuration prometheus.Histogram

	// getUserDuration is a Prometheus histogram that tracks the
	// duration of user retrieval operations.
	getUserDuration prometheus.Histogram

	// updateUserDuration is a Prometheus histogram that tracks the
	// duration of user update operations.
	updateUserDuration prometheus.Histogram

	// deleteUserDuration is a Prometheus histogram that tracks the
	// duration of user deletion operations.
	deleteUserDuration prometheus.Histogram

	// cacheHits is a Prometheus counter that tracks the number of
	// cache hits for user data.
	cacheHits prometheus.Counter

	// cacheMisses is a Prometheus counter that tracks the number of
	// cache misses for user data.
	cacheMisses prometheus.Counter
}

// NewMetrics creates and initializes a new Metrics instance.
//
// This function serves as a constructor for the Metrics struct, initializing
// all the Prometheus metrics with their appropriate configuration.
//
// The metrics defined here are:
//   - users_total: A counter that tracks the total number of users
//   - users_active: A gauge that tracks the current number of active users
//   - users_create_duration: Histogram of user creation operation durations
//   - users_get_duration: Histogram of user retrieval operation durations
//   - users_update_duration: Histogram of user update operation durations
//   - users_delete_duration: Histogram of user deletion operation durations
//   - users_cache_hits: Counter of cache hits for user data
//   - users_cache_misses: Counter of cache misses for user data
//
// Returns:
//   - *Metrics: A pointer to the newly created Metrics instance
//
// Example:
//
//	metrics := users.NewMetrics()
//	// Use metrics in your service
//	service := users.New(config, repo, metrics, logger)
func NewMetrics() *Metrics {
	return &Metrics{
		totalUsers: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      "total",
			Help:      "Total number of users created",
		}),
		activeUsers: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      "active",
			Help:      "Current number of active users",
		}),
		createUserDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      "create_duration_seconds",
			Help:      "Duration of user creation operations in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
		getUserDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      "get_duration_seconds",
			Help:      "Duration of user retrieval operations in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
		updateUserDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      "update_duration_seconds",
			Help:      "Duration of user update operations in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
		deleteUserDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      "delete_duration_seconds",
			Help:      "Duration of user deletion operations in seconds",
			Buckets:   prometheus.DefBuckets,
		}),
		cacheHits: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      "cache_hits_total",
			Help:      "Total number of cache hits for user data",
		}),
		cacheMisses: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: MetricsNamespace,
			Subsystem: MetricsSubsystem,
			Name:      "cache_misses_total",
			Help:      "Total number of cache misses for user data",
		}),
	}
}

// IncTotalUsers increments the total users counter.
//
// This method should be called whenever a new user is created in the system.
func (m *Metrics) IncTotalUsers() {
	m.totalUsers.Inc()
}

// IncActiveUsers increments the active users gauge.
func (m *Metrics) IncActiveUsers() {
	m.activeUsers.Inc()
}

// DecActiveUsers decrements the active users gauge.
func (m *Metrics) DecActiveUsers() {
	m.activeUsers.Dec()
}

// ObserveCreateUserDuration records the duration of a user creation operation.
func (m *Metrics) ObserveCreateUserDuration(seconds float64) {
	m.createUserDuration.Observe(seconds)
}

// ObserveGetUserDuration records the duration of a user retrieval operation.
func (m *Metrics) ObserveGetUserDuration(seconds float64) {
	m.getUserDuration.Observe(seconds)
}

// ObserveUpdateUserDuration records the duration of a user update operation.
func (m *Metrics) ObserveUpdateUserDuration(seconds float64) {
	m.updateUserDuration.Observe(seconds)
}

// ObserveDeleteUserDuration records the duration of a user deletion operation.
func (m *Metrics) ObserveDeleteUserDuration(seconds float64) {
	m.deleteUserDuration.Observe(seconds)
}

// IncCacheHits increments the cache hits counter.
func (m *Metrics) IncCacheHits() {
	m.cacheHits.Inc()
}

// IncCacheMisses increments the cache misses counter.
func (m *Metrics) IncCacheMisses() {
	m.cacheMisses.Inc()
}
