package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HTTP 基础指标
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cbctf_http_request_duration_seconds",
			Help:    "Histogram of response time for handler.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HttpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cbctf_http_response_size_bytes",
			Help:    "Histogram of HTTP response body sizes.",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"method", "path"},
	)

	InFlightRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cbctf_http_in_flight_requests",
			Help: "Current number of in-flight requests being handled.",
		},
	)

	// CTF 业务事件指标（低基数标签）
	FlagSubmissionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_flag_submissions_total",
			Help: "Total number of flag submissions",
		},
		[]string{"contest_id", "challenge_type", "status"},
	)

	BloodTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_blood_total",
			Help: "Total blood events (first/second/third)",
		},
		[]string{"contest_id", "blood_order"},
	)

	UserRegistrationTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_user_registrations_total",
			Help: "Total number of user registrations",
		},
		[]string{"oauth_provider"},
	)

	UserLoginTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_user_logins_total",
			Help: "Total number of user logins",
		},
		[]string{"oauth_provider"},
	)

	FileUploadTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_file_uploads_total",
			Help: "Total number of file uploads",
		},
		[]string{"file_type"},
	)

	FileUploadSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cbctf_file_upload_size_bytes",
			Help:    "Size of uploaded files in bytes",
			Buckets: prometheus.ExponentialBuckets(1024, 2, 15), // 1KB ~ 32MB
		},
		[]string{"file_type"},
	)

	EmailSentTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_email_sent_total",
			Help: "Total number of emails sent",
		},
		[]string{"status"},
	)

	RateLimitHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
		[]string{"endpoint"},
	)

	CheatDetectionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_cheat_detections_total",
			Help: "Total number of cheat detection events",
		},
		[]string{"reason_type"},
	)

	// Cron Job 监控
	CronJobDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cbctf_cron_job_duration_seconds",
			Help:    "Cron job execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"job_name"},
	)

	CronJobRunsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_cron_job_runs_total",
			Help: "Total number of cron job runs",
		},
		[]string{"job_name", "status"},
	)

	// 异步任务监控
	TaskEnqueuedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_task_enqueued_total",
			Help: "Total number of enqueued async tasks",
		},
		[]string{"task_type"},
	)

	TaskProcessedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cbctf_task_processed_total",
			Help: "Total number of processed async tasks",
		},
		[]string{"task_type", "status"},
	)

	TaskProcessingDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cbctf_task_processing_duration_seconds",
			Help:    "Async task processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"task_type"},
	)
)
