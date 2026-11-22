package prometheus

import "strconv"

func RecordFlagSubmission(contestID uint, challengeType string, solved bool) {
	status := "failed"
	if solved {
		status = "success"
	}
	FlagSubmissionsTotal.WithLabelValues(strconv.FormatUint(uint64(contestID), 10), challengeType, status).Inc()
}

func RecordBlood(contestID uint, order string) {
	BloodTotal.WithLabelValues(strconv.FormatUint(uint64(contestID), 10), order).Inc()
}

func RecordUserRegister(provider string) {
	UserRegistrationTotal.WithLabelValues(provider).Inc()
}

func RecordUserLogin(provider string) {
	UserLoginTotal.WithLabelValues(provider).Inc()
}

func RecordPasswordReset(success bool) {
	PasswordResetTotal.WithLabelValues(statusLabel(success)).Inc()
}

func RecordFileUpload(uploadType string, size int64) {
	FileUploadTotal.WithLabelValues(uploadType).Inc()
	FileUploadSize.WithLabelValues(uploadType).Observe(float64(size))
}

func RecordEmailSent(emailKind string, success bool) {
	EmailSentTotal.WithLabelValues(emailKind, statusLabel(success)).Inc()
}

func RecordRateLimitHit(endpoint string) {
	RateLimitHits.WithLabelValues(endpoint).Inc()
}

func RecordCheatDetection(contestID uint, reasonType string) {
	CheatDetectionsTotal.WithLabelValues(strconv.FormatUint(uint64(contestID), 10), reasonType).Inc()
}

// RecordVictimTransition 记录一次靶机状态跃迁。
func RecordVictimTransition(action string, success bool) {
	VictimTransitionsTotal.WithLabelValues(action, statusLabel(success)).Inc()
}

// RecordVictimLifecycle 记录靶机从创建到就绪（或失败）的耗时。
func RecordVictimLifecycle(duration float64, success bool) {
	VictimLifecycleDuration.WithLabelValues(statusLabel(success)).Observe(duration)
}

func RecordVictimReadinessFailure(reason string) {
	VictimReadinessFailuresTotal.WithLabelValues(reason).Inc()
}

func RecordK8sOperationError(operation string) {
	K8sOperationErrorsTotal.WithLabelValues(operation).Inc()
}

func RecordCronJob(jobName string, duration float64, success bool) {
	CronJobDuration.WithLabelValues(jobName).Observe(duration)
	CronJobRunsTotal.WithLabelValues(jobName, statusLabel(success)).Inc()
}

func RecordTaskEnqueued(taskType string, success bool) {
	TaskEnqueuedTotal.WithLabelValues(taskType, statusLabel(success)).Inc()
}

func RecordTaskProcessed(taskType string, duration float64, success bool) {
	status := statusLabel(success)
	TaskProcessingDuration.WithLabelValues(taskType, status).Observe(duration)
	TaskProcessedTotal.WithLabelValues(taskType, status).Inc()
}

func statusLabel(success bool) string {
	if success {
		return "success"
	}
	return "failed"
}
