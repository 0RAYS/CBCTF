package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type SubmissionRepo struct {
	BaseRepo[model.Submission]
}

type BloodRankRow struct {
	ContestFlagID uint
	TeamID        uint
	BloodRank     int
}

type UpdateSubmissionOptions struct {
	Solved *bool
	Score  *float64
}

func (u UpdateSubmissionOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Solved != nil {
		options["solved"] = *u.Solved
	}
	if u.Score != nil {
		options["score"] = *u.Score
	}
	return options
}

func InitSubmissionRepo(tx *gorm.DB) *SubmissionRepo {
	return &SubmissionRepo{
		DB: tx,
	}
}

// LockAttemptScope 在当前事务内锁定队伍和赛事题目维度的提交尝试。
// 该锁用于让尝试次数检查和提交记录创建串行化，避免同队伍同题并发错误提交绕过 Attempt 限制。
func (s *SubmissionRepo) LockAttemptScope(teamID, contestChallengeID uint) model.RetVal {
	teamKey := fmt.Sprintf("team:%d", teamID)
	challengeKey := fmt.Sprintf("contest_challenge:%d", contestChallengeID)
	res := s.DB.Exec("SELECT pg_advisory_xact_lock(hashtext(?), hashtext(?))", teamKey, challengeKey)
	if res.Error != nil {
		target := fmt.Sprintf("contest_challenge_id=%d,team_id=%d", contestChallengeID, teamID)
		log.Logger.Warningf("Failed to lock submission attempt scope: %s (target=%s)", res.Error, target)
		return model.FailRetVal(i18n.Model.Submission.GetError, "Submission", target, res.Error)
	}
	return model.SuccessRetVal()
}

func (s *SubmissionRepo) DeleteByUserID(userIDL ...uint) model.RetVal {
	return s.DeleteByFieldID("user_id", userIDL...)
}

func (s *SubmissionRepo) DeleteByTeamID(teamIDL ...uint) model.RetVal {
	return s.DeleteByFieldID("team_id", teamIDL...)
}

func (s *SubmissionRepo) DeleteByContestChallengeID(contestChallengeIDL ...uint) model.RetVal {
	return s.DeleteByFieldID("contest_challenge_id", contestChallengeIDL...)
}

func (s *SubmissionRepo) DeleteByContestFlagID(contestFlagIDL ...uint) model.RetVal {
	return s.DeleteByFieldID("contest_flag_id", contestFlagIDL...)
}

func (s *SubmissionRepo) GetBloodRankMap(contestFlagIDL ...uint) (map[uint]map[uint]int, model.RetVal) {
	rankMap := make(map[uint]map[uint]int)
	if len(contestFlagIDL) == 0 {
		return rankMap, model.SuccessRetVal()
	}

	firstSolves := s.DB.Table("submissions").
		Select("contest_flag_id, team_id, MIN(created_at) AS first_solved_at, MIN(id) AS first_submission_id").
		Where("contest_flag_id IN ? AND solved = true AND team_id <> 0 AND deleted_at IS NULL", contestFlagIDL).
		Group("contest_flag_id, team_id")

	ranked := s.DB.Table("(?) AS first_solves", firstSolves).
		Select("contest_flag_id, team_id, ROW_NUMBER() OVER (PARTITION BY contest_flag_id ORDER BY first_solved_at ASC, first_submission_id ASC) AS blood_rank")

	var rows []BloodRankRow
	res := s.DB.Table("(?) AS ranked", ranked).
		Select("contest_flag_id, team_id, blood_rank").
		Where("blood_rank <= 3").
		Scan(&rows)
	if res.Error != nil {
		target := fmt.Sprintf("contest_flag_id=%v", contestFlagIDL)
		log.Logger.Warningf("Failed to get blood rank: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Submission.GetError, "Submission", target, res.Error)
	}

	for _, row := range rows {
		if rankMap[row.ContestFlagID] == nil {
			rankMap[row.ContestFlagID] = make(map[uint]int)
		}
		rankMap[row.ContestFlagID][row.TeamID] = row.BloodRank
	}
	return rankMap, model.SuccessRetVal()
}

type FlagSolverRow struct {
	SolvedAt time.Time
	UserName string
	TeamName string
	UserID   uint
	TeamID   uint
	Score    float64
}

func (s *SubmissionRepo) ListFlagSolvers(contestFlagID uint) ([]FlagSolverRow, model.RetVal) {
	var rows []FlagSolverRow
	res := s.DB.Table("submissions").
		Select("submissions.user_id, users.name AS user_name, submissions.team_id, teams.name AS team_name, submissions.score, submissions.created_at AS solved_at").
		Joins("INNER JOIN users ON submissions.user_id = users.id AND users.deleted_at IS NULL").
		Joins("INNER JOIN teams ON submissions.team_id = teams.id AND teams.deleted_at IS NULL").
		Where("submissions.contest_flag_id = ? AND submissions.solved = true AND submissions.deleted_at IS NULL", contestFlagID).
		Order("submissions.created_at ASC, submissions.id ASC").
		Scan(&rows)
	if res.Error != nil {
		target := fmt.Sprintf("contest_flag_id=%d", contestFlagID)
		log.Logger.Warningf("Failed to get flag solvers: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Submission.GetError, "Submission", target, res.Error)
	}
	return rows, model.SuccessRetVal()
}

// ContestIPSubmission 是某个 IP 上的一次 flag 提交及其题目。
type ContestIPSubmission struct {
	IP                 string
	SubmissionID       uint
	TeamID             uint
	UserID             uint
	ContestChallengeID uint
	ChallengeName      string
	SubmittedAt        time.Time
	Value              string
	Solved             bool
}

// ListContestIPSubmissions 返回给定 IP 上的 flag 提交，每个 IP 最多保留 limit 条最早记录。
func (s *SubmissionRepo) ListContestIPSubmissions(contestID uint, ips []string, start, end time.Time, limit int) ([]ContestIPSubmission, model.RetVal) {
	if contestID == 0 || len(ips) == 0 {
		return nil, model.SuccessRetVal()
	}

	var submissions []ContestIPSubmission
	res := s.DB.Raw(`
		WITH scoped AS (
			SELECT submissions.ip, submissions.id AS submission_id, submissions.team_id, submissions.user_id,
				submissions.contest_challenge_id, submissions.created_at AS submitted_at, submissions.value, submissions.solved,
				contest_challenges.name AS challenge_name,
				ROW_NUMBER() OVER (PARTITION BY submissions.ip ORDER BY submissions.created_at ASC, submissions.id ASC) AS event_rank
			FROM submissions
			INNER JOIN contest_challenges ON contest_challenges.id = submissions.contest_challenge_id AND contest_challenges.deleted_at IS NULL
			WHERE submissions.contest_id = ?
				AND submissions.deleted_at IS NULL
				AND submissions.created_at >= ? AND submissions.created_at <= ?
				AND submissions.ip IN ?
		)
		SELECT ip, submission_id, team_id, user_id, contest_challenge_id, challenge_name, submitted_at, value, solved
		FROM scoped
		WHERE event_rank <= ?
		ORDER BY ip ASC, submitted_at ASC
	`, contestID, start, end, ips, limit).Scan(&submissions)
	if res.Error != nil {
		target := fmt.Sprintf("contest_id=%d,ips=%d", contestID, len(ips))
		log.Logger.Warningf("Failed to list contest IP submissions: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Submission.GetError, "Submission", target, res.Error)
	}
	return submissions, model.SuccessRetVal()
}

func (s *SubmissionRepo) ListSolvedByTeamID(teamIDL ...uint) ([]model.Submission, model.RetVal) {
	if len(teamIDL) == 0 {
		return nil, model.SuccessRetVal()
	}
	submissions := make([]model.Submission, 0)
	res := s.DB.Model(&model.Submission{}).
		Where("team_id IN ? AND solved = true", teamIDL).
		Order("team_id ASC, created_at ASC, id ASC").
		Find(&submissions)
	if res.Error != nil {
		target := fmt.Sprintf("team_id=%v", teamIDL)
		log.Logger.Warningf("Failed to list solved submissions: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.Submission.GetError, "Submission", target, res.Error)
	}
	return submissions, model.SuccessRetVal()
}
