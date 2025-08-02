package prometheus

import (
	"CBCTF/internal/model"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for handler.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	HttpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Histogram of HTTP request body sizes.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10), // 100B ~ 51KB
		},
		[]string{"path", "method"},
	)

	HttpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Histogram of HTTP response body sizes.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"path", "method"},
	)

	InFlightRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "in_flight_requests",
			Help: "Current number of in-flight requests being handled.",
		},
	)

	FlagSubmissionTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ctf_flag_submissions_total",
			Help: "Total number of flag submissions",
		},
		[]string{"contest_name", "challenge_name", "team_name", "status", "challenge_type"},
	)

	ContestActiveTeams = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ctf_contest_active_teams",
			Help: "Number of active teams in contest",
		},
		[]string{"contest_name"},
	)

	ContestActiveUsers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ctf_contest_active_users",
			Help: "Number of active users in contest",
		},
		[]string{"contest_id", "contest_name"},
	)

	VictimContainerTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ctf_victim_containers_total",
			Help: "Total number of victim containers running",
		},
		[]string{"contest_name", "challenge_name"},
	)

	UserRegistrationTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ctf_user_registrations_total",
			Help: "Total number of user registrations",
		},
		[]string{"oauth_provider"},
	)

	FileUploadTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ctf_file_uploads_total",
			Help: "Total number of file uploads",
		},
		[]string{"file_type"},
	)

	FileUploadSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ctf_file_upload_size_bytes",
			Help:    "Size of uploaded files",
			Buckets: prometheus.ExponentialBuckets(1024, 2, 15), // 1KB ~ 32MB
		},
		[]string{"file_type", "contest_id"},
	)

	WebSocketConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ctf_websocket_connections",
			Help: "Current number of WebSocket connections",
		},
	)

	EmailSentTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ctf_email_sent_total",
			Help: "Total number of email sent",
		},
		[]string{"status"},
	)

	CacheHitRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ctf_cache_hit_rate",
			Help: "Cache hit rate (0-1)",
		},
		[]string{"cache_type"},
	)

	RateLimitHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ctf_rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
		[]string{"endpoint", "ip"},
	)

	ErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ctf_errors_total",
			Help: "Total number of errors",
		},
		[]string{"error_type", "component"},
	)
)
