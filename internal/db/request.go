package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

type RequestRepo struct {
	BaseRepo[model.Request]
}

func InitRequestRepo(tx *gorm.DB) *RequestRepo {
	return &RequestRepo{
		DB: tx,
	}
}

func (r *RequestRepo) Create(requests ...model.Request) model.RetVal {
	if len(requests) == 0 {
		return model.SuccessRetVal()
	}
	for i := range requests {
		if requests[i].RecordID == "" {
			requests[i].RecordID = utils.UUID()
		}
	}
	if res := r.DB.Model(&model.Request{}).Omit("id").Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "record_id"}}, DoNothing: true}).CreateInBatches(requests, 200); res.Error != nil {
		log.Logger.Warningf("Failed to create requests: %s", res.Error)
		return model.FailRetVal(i18n.Model.CreateError, model.Name(model.Request{}), fmt.Sprintf("count=%d", len(requests)), res.Error)
	}
	return model.SuccessRetVal()
}

func (r *RequestRepo) CountIP() (int64, model.RetVal) {
	var count int64
	res := r.DB.Model(&model.Request{}).Distinct("ip").Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count Request: %s", res.Error)
		return 0, model.FailRetVal(i18n.Model.Request.GetError, "Request", "", res.Error)
	}
	return count, model.SuccessRetVal()
}

// ContestSharedIPUser 汇总某个共享 IP 上单个用户在比赛期间的请求情况。
// MinGapSeconds 与 MaxGapSeconds 是两个不同用户请求之间的最小与最大间隔。
type ContestSharedIPUser struct {
	IP            string
	UserID        uint
	FirstTime     time.Time
	LastTime      time.Time
	RequestCount  int64
	TotalRequests int64
	TotalUsers    int64
	FirstSeen     time.Time
	LastSeen      time.Time
	MinGapSeconds float64
	MaxGapSeconds float64
}

// ListSharedContestUserIPs 返回比赛中被多个用户共用的 IP 及每个用户的请求聚合。
// 最小与最大间隔按相邻的不同用户请求计算，用于衡量共用行为的紧密程度。
func (r *RequestRepo) ListSharedContestUserIPs(contestID uint, start, end time.Time) ([]ContestSharedIPUser, model.RetVal) {
	if contestID == 0 {
		return nil, model.SuccessRetVal()
	}

	var userIPL []ContestSharedIPUser
	res := r.DB.Raw(`
		WITH scoped AS (
			SELECT requests.ip, requests.user_id, requests.time
			FROM requests
			INNER JOIN user_contests ON user_contests.user_id = requests.user_id
			INNER JOIN users ON users.id = requests.user_id AND users.deleted_at IS NULL
			WHERE user_contests.contest_id = ?
				AND requests.deleted_at IS NULL
				AND requests.time >= ? AND requests.time <= ?
		),
		per_user AS (
			SELECT ip, user_id, MIN(time) AS first_time, MAX(time) AS last_time, COUNT(*) AS request_count
			FROM scoped
			GROUP BY ip, user_id
		),
		shared AS (
			SELECT ip FROM per_user GROUP BY ip HAVING COUNT(*) > 1
		),
		ordered AS (
			SELECT scoped.ip, scoped.user_id, scoped.time,
				LAG(scoped.time) OVER (PARTITION BY scoped.ip ORDER BY scoped.time) AS prev_time,
				LAG(scoped.user_id) OVER (PARTITION BY scoped.ip ORDER BY scoped.time) AS prev_user
			FROM scoped
			INNER JOIN shared ON shared.ip = scoped.ip
		),
		gaps AS (
			SELECT ip,
				MIN(EXTRACT(EPOCH FROM (time - prev_time))) FILTER (WHERE prev_user IS DISTINCT FROM user_id) AS min_gap,
				MAX(EXTRACT(EPOCH FROM (time - prev_time))) FILTER (WHERE prev_user IS DISTINCT FROM user_id) AS max_gap
			FROM ordered
			GROUP BY ip
		)
		SELECT
			pu.ip,
			pu.user_id,
			pu.first_time,
			pu.last_time,
			pu.request_count,
			SUM(pu.request_count) OVER (PARTITION BY pu.ip)::bigint AS total_requests,
			COUNT(*) OVER (PARTITION BY pu.ip) AS total_users,
			MIN(pu.first_time) OVER (PARTITION BY pu.ip) AS first_seen,
			MAX(pu.last_time) OVER (PARTITION BY pu.ip) AS last_seen,
			g.min_gap::float8 AS min_gap_seconds,
			g.max_gap::float8 AS max_gap_seconds
		FROM per_user pu
		INNER JOIN shared s ON s.ip = pu.ip
		INNER JOIN gaps g ON g.ip = pu.ip
		ORDER BY pu.ip ASC, pu.first_time ASC, pu.user_id ASC
	`, contestID, start, end).Scan(&userIPL)
	if res.Error != nil {
		target := fmt.Sprintf("contest_id=%d,end=%s,start=%s", contestID, end.Format(time.RFC3339), start.Format(time.RFC3339))
		log.Logger.Warningf("Failed to list shared request IPs: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Request.GetError, "Request", target, res.Error)
	}
	return userIPL, model.SuccessRetVal()
}

// ContestIPRequest 是共享 IP 上的一次网页请求。
type ContestIPRequest struct {
	IP        string
	UserID    uint
	Time      time.Time
	Method    string
	Path      string
	Status    int
	Latency   time.Duration
	UserAgent string
}

// ListContestIPRequests 返回给定 IP 上的网页请求明细，每个 IP 最多保留 limit 条最早记录。
func (r *RequestRepo) ListContestIPRequests(contestID uint, ips []string, start, end time.Time, limit int) ([]ContestIPRequest, model.RetVal) {
	if contestID == 0 || len(ips) == 0 {
		return nil, model.SuccessRetVal()
	}

	var requests []ContestIPRequest
	res := r.DB.Raw(`
		WITH scoped AS (
			SELECT requests.ip, requests.user_id, requests.time, requests.method, requests.path,
				requests.status, requests.latency, requests.user_agent,
				ROW_NUMBER() OVER (PARTITION BY requests.ip ORDER BY requests.time ASC, requests.id ASC) AS event_rank
			FROM requests
			INNER JOIN user_contests ON user_contests.user_id = requests.user_id
			INNER JOIN users ON users.id = requests.user_id AND users.deleted_at IS NULL
			WHERE user_contests.contest_id = ?
				AND requests.deleted_at IS NULL
				AND requests.time >= ? AND requests.time <= ?
				AND requests.ip IN ?
		)
		SELECT ip, user_id, time, method, path, status, latency, user_agent
		FROM scoped
		WHERE event_rank <= ?
		ORDER BY ip ASC, time ASC
	`, contestID, start, end, ips, limit).Scan(&requests)
	if res.Error != nil {
		target := fmt.Sprintf("contest_id=%d,ips=%d", contestID, len(ips))
		log.Logger.Warningf("Failed to list contest IP requests: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Request.GetError, "Request", target, res.Error)
	}
	return requests, model.SuccessRetVal()
}
