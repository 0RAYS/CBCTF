package prometheus

import "fmt"

func RecordFlagSubmission(contestID uint, challengeType string, solved bool) {
	status := "failed"
	if solved {
		status = "success"
	}
	FlagSubmissionsTotal.WithLabelValues(fmt.Sprintf("%d", contestID), challengeType, status).Inc()
}

func RecordBlood(contestID uint, order string) {
	BloodTotal.WithLabelValues(fmt.Sprintf("%d", contestID), order).Inc()
}

func RecordUserRegister(provider string) {
	UserRegistrationTotal.WithLabelValues(provider).Inc()
}

func RecordUserLogin(provider string) {
	UserLoginTotal.WithLabelValues(provider).Inc()
}

func RecordFileUpload(fileType string, size int64) {
	FileUploadTotal.WithLabelValues(fileType).Inc()
	FileUploadSize.WithLabelValues(fileType).Observe(float64(size))
}

func RecordEmailSent(success bool) {
	if success {
		EmailSentTotal.WithLabelValues("success").Inc()
	} else {
		EmailSentTotal.WithLabelValues("failed").Inc()
	}
}

func RecordRateLimitHit(endpoint string) {
	RateLimitHits.WithLabelValues(endpoint).Inc()
}

func RecordCheatDetection(reasonType string) {
	CheatDetectionsTotal.WithLabelValues(reasonType).Inc()
}

func RecordCronJob(jobName string, duration float64, success bool) {
	CronJobDuration.WithLabelValues(jobName).Observe(duration)
	if success {
		CronJobRunsTotal.WithLabelValues(jobName, "success").Inc()
	} else {
		CronJobRunsTotal.WithLabelValues(jobName, "failed").Inc()
	}
}

func RecordTaskEnqueued(taskType string) {
	TaskEnqueuedTotal.WithLabelValues(taskType).Inc()
}

func RecordTaskProcessed(taskType string, duration float64, success bool) {
	TaskProcessingDuration.WithLabelValues(taskType).Observe(duration)
	if success {
		TaskProcessedTotal.WithLabelValues(taskType, "success").Inc()
	} else {
		TaskProcessedTotal.WithLabelValues(taskType, "failed").Inc()
	}
}
