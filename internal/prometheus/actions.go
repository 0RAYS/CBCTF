package prometheus

import "CBCTF/internal/model"

func UpdateFlagSubmissionMetrics(contest model.Contest, contestChallenge model.ContestChallenge, team model.Team, solved bool) {
	if solved {
		FlagSubmissionTotal.WithLabelValues(contest.Name, contestChallenge.Name, team.Name, "success", contestChallenge.Type).Inc()
	} else {
		FlagSubmissionTotal.WithLabelValues(contest.Name, contestChallenge.Name, team.Name, "failed", contestChallenge.Type).Inc()
	}
}

func AddContestActiveTeamsMetrics(contest model.Contest, count int) {
	ContestActiveTeams.WithLabelValues(contest.Name).Add(float64(count))
}

func SubContestActiveTeamsMetrics(contest model.Contest, count int) {
	ContestActiveTeams.WithLabelValues(contest.Name).Sub(float64(count))

}

func AddContestActiveUsersMetrics(contest model.Contest, count int) {
	ContestActiveUsers.WithLabelValues(contest.Name).Add(float64(count))
}

func SubContestActiveUsersMetrics(contest model.Contest, count int) {
	ContestActiveUsers.WithLabelValues(contest.Name).Sub(float64(count))
}

func AddVictimContainerMetrics(contest model.Contest, contestChallenge model.ContestChallenge, count int) {
	VictimContainerTotal.WithLabelValues(contest.Name, contestChallenge.Name).Add(float64(count))
}

func SubVictimContainerMetrics(contest model.Contest, contestChallenge model.ContestChallenge, count int) {
	VictimContainerTotal.WithLabelValues(contest.Name, contestChallenge.Name).Sub(float64(count))
}

func UpdateUserRegisterMetrics(provider string) {
	UserRegistrationTotal.WithLabelValues(provider).Inc()
}

func UpdateUserLoginMetrics(provider string) {
	UserLoginTotal.WithLabelValues(provider).Inc()
}

func UpdateFileUploadMetrics(fileType string, size int64) {
	FileUploadTotal.WithLabelValues(fileType).Inc()
	FileUploadSize.WithLabelValues(fileType).Observe(float64(size))
}

func UpdateWebSocketMetrics(connections int) {
	WebSocketConnections.Set(float64(connections))
}

func IncEmailSentMetrics(success bool) {
	if success {
		EmailSentTotal.WithLabelValues("success").Inc()
	} else {
		EmailSentTotal.WithLabelValues("failed").Inc()
	}
}

func UpdateCacheMetrics(cacheType string, hits, misses int64) {
	hitRate := float64(hits) / float64(hits+misses)
	CacheHitRate.WithLabelValues(cacheType).Set(hitRate)
}

func UpdateRateLimitMetrics(endpoint, ip string) {
	RateLimitHits.WithLabelValues(endpoint, ip).Inc()
}
