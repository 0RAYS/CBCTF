package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type TeamFlagRepo struct {
	BaseRepo[model.TeamFlag]
}

// CreateIfAbsent ON CONFLICT avoids aborting an enclosing PostgreSQL transaction when two
// requests initialize the same flag concurrently. The stored value wins.
func (t *TeamFlagRepo) CreateIfAbsent(flag model.TeamFlag) (model.TeamFlag, model.RetVal) {
	res := t.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&flag)
	if res.Error != nil {
		target := fmt.Sprintf("contest_flag_id=%d,team_id=%d", flag.ContestFlagID, flag.TeamID)
		return model.TeamFlag{}, model.FailRetVal(i18n.Model.CreateError, model.Name(model.TeamFlag{}), target, res.Error)
	}
	return t.Get(GetOptions{Conditions: map[string]any{"team_id": flag.TeamID, "contest_flag_id": flag.ContestFlagID}})
}

type UpdateTeamFlagRepo struct {
	Value  *string
	Solved *bool
}

func (u UpdateTeamFlagRepo) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Value != nil {
		options["value"] = *u.Value
	}
	if u.Solved != nil {
		options["solved"] = *u.Solved
	}
	return options
}

func InitTeamFlagRepo(tx *gorm.DB) *TeamFlagRepo {
	return &TeamFlagRepo{
		DB: tx,
	}
}

func (t *TeamFlagRepo) DeleteByTeamID(teamIDL ...uint) model.RetVal {
	return t.DeleteByFieldID("team_id", teamIDL...)
}

func (t *TeamFlagRepo) DeleteByContestFlagID(contestFlagIDL ...uint) model.RetVal {
	return t.DeleteByFieldID("contest_flag_id", contestFlagIDL...)
}

// TeamFlagWithChallenge 是队伍解题记录与关联比赛 flag 的查询结果。
// ContestFlag 字段
type TeamFlagWithChallenge struct {
	ContestFlagValue         string
	ContestChallengeName     string
	ContestChallengeCategory string
	ChallengeRandID          string
	ChallengeName            string
	model.TeamFlag
	ContestFlagScore        float64
	ContestFlagCurrentScore float64
	ContestFlagMinScore     float64
	ContestFlagScoreType    uint
	ChallengeID             uint
	ContestChallengeHidden  bool
}

type SolvedContestFlagRow struct {
	model.ContestFlag
	TeamID uint
}

type ContestWrongFlagSubmission struct {
	CreatedAt          time.Time
	MatchedSolvedAt    time.Time
	Value              string
	IP                 string
	SubmissionID       uint
	TeamID             uint
	UserID             uint
	ContestChallengeID uint
	ChallengeName      string
	MatchedTeamID      uint
	MatchedUserID      uint
}

func (t *TeamFlagRepo) GetTeamFlagsWithChallenge(teamIDL ...uint) ([]TeamFlagWithChallenge, model.RetVal) {
	if len(teamIDL) == 0 {
		return nil, model.SuccessRetVal()
	}
	var results []TeamFlagWithChallenge
	res := t.DB.Table("team_flags").
		Select(`team_flags.*,
			contest_flags.value AS contest_flag_value,
			contest_flags.score AS contest_flag_score,
			contest_flags.current_score AS contest_flag_current_score,
			contest_flags.min_score AS contest_flag_min_score,
			contest_flags.score_type AS contest_flag_score_type,
			contest_challenges.name AS contest_challenge_name,
			contest_challenges.category AS contest_challenge_category,
			contest_challenges.hidden AS contest_challenge_hidden,
			contest_challenges.challenge_id AS challenge_id,
			contest_challenges.rand_id AS challenge_rand_id,
			contest_challenges.name AS challenge_name`).
		Joins("INNER JOIN contest_flags ON team_flags.contest_flag_id = contest_flags.id AND contest_flags.deleted_at IS NULL").
		Joins("INNER JOIN contest_challenges ON contest_flags.contest_challenge_id = contest_challenges.id AND contest_challenges.deleted_at IS NULL").
		Where("team_flags.team_id IN ? AND team_flags.deleted_at IS NULL", teamIDL).
		Scan(&results)
	target := fmt.Sprintf("team_id=%v", teamIDL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get TeamFlags: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.TeamFlag.GetError, "TeamFlag", target, res.Error)
	}
	return results, model.SuccessRetVal()
}

func (t *TeamFlagRepo) ListContestWrongFlagSubmissions(contestID uint, start, end time.Time) ([]ContestWrongFlagSubmission, model.RetVal) {
	if contestID == 0 {
		return nil, model.SuccessRetVal()
	}

	var results []ContestWrongFlagSubmission
	res := t.DB.Table("submissions").
		Select(`submissions.id AS submission_id,
			submissions.team_id,
			submissions.user_id,
			submissions.contest_challenge_id,
			contest_challenges.name AS challenge_name,
			submissions.value,
			submissions.ip,
			submissions.created_at,
			team_flags.team_id AS matched_team_id,
			owner_solve.user_id AS matched_user_id,
			owner_solve.solved_at AS matched_solved_at`).
		Joins("INNER JOIN team_flags ON team_flags.value = submissions.value AND team_flags.deleted_at IS NULL").
		Joins("INNER JOIN contest_flags ON contest_flags.id = team_flags.contest_flag_id AND contest_flags.deleted_at IS NULL").
		Joins("INNER JOIN teams AS submission_teams ON submission_teams.id = submissions.team_id AND submission_teams.deleted_at IS NULL").
		Joins("INNER JOIN teams AS matched_teams ON matched_teams.id = team_flags.team_id AND matched_teams.deleted_at IS NULL").
		Joins("INNER JOIN contest_challenges ON contest_challenges.id = submissions.contest_challenge_id AND contest_challenges.deleted_at IS NULL").
		// 合法持有者的首次解题决定了冒用发生得有多快。
		Joins(`LEFT JOIN LATERAL (
			SELECT solved_submissions.user_id, solved_submissions.created_at AS solved_at
			FROM submissions AS solved_submissions
			WHERE solved_submissions.contest_flag_id = team_flags.contest_flag_id
				AND solved_submissions.team_id = team_flags.team_id
				AND solved_submissions.solved = true
				AND solved_submissions.deleted_at IS NULL
			ORDER BY solved_submissions.created_at ASC, solved_submissions.id ASC
			LIMIT 1
		) AS owner_solve ON TRUE`).
		Where("submissions.contest_id = ? AND submissions.solved = false AND submissions.deleted_at IS NULL", contestID).
		Where("submissions.created_at >= ? AND submissions.created_at <= ?", start, end).
		Where("contest_flags.contest_challenge_id = submissions.contest_challenge_id").
		Where("submissions.team_id <> team_flags.team_id").
		Where(`NOT EXISTS (
			SELECT 1
			FROM team_flags AS own_team_flags
			INNER JOIN contest_flags AS own_contest_flags
				ON own_contest_flags.id = own_team_flags.contest_flag_id
				AND own_contest_flags.deleted_at IS NULL
			WHERE own_team_flags.team_id = submissions.team_id
				AND own_team_flags.value = submissions.value
				AND own_team_flags.deleted_at IS NULL
				AND own_contest_flags.contest_challenge_id = submissions.contest_challenge_id
		)`).
		Where("submission_teams.contest_id = ? AND matched_teams.contest_id = ?", contestID, contestID).
		Group(`submissions.id,
			submissions.team_id,
			submissions.user_id,
			submissions.contest_challenge_id,
			contest_challenges.name,
			submissions.value,
			submissions.ip,
			submissions.created_at,
			team_flags.team_id,
			owner_solve.user_id,
			owner_solve.solved_at`).
		Order("submissions.created_at ASC, submissions.id ASC, team_flags.team_id ASC").
		Scan(&results)
	target := fmt.Sprintf("contest_id=%d", contestID)
	if res.Error != nil {
		log.Logger.Warningf("Failed to list wrong flag submissions: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.TeamFlag.GetError, "TeamFlag", target, res.Error)
	}
	return results, model.SuccessRetVal()
}

func (t *TeamFlagRepo) GetSolvedContestFlags(teamIDL ...uint) ([]SolvedContestFlagRow, model.RetVal) {
	if len(teamIDL) == 0 {
		return nil, model.SuccessRetVal()
	}
	var results []SolvedContestFlagRow
	res := t.DB.Table("team_flags").
		Select("team_flags.team_id, contest_flags.*").
		Joins("INNER JOIN contest_flags ON contest_flags.id = team_flags.contest_flag_id AND contest_flags.deleted_at IS NULL").
		Where("team_flags.team_id IN ? AND team_flags.solved = true AND team_flags.deleted_at IS NULL", teamIDL).
		Order("team_flags.team_id ASC, contest_flags.id ASC").
		Scan(&results)
	target := fmt.Sprintf("team_id=%v", teamIDL)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get solved team contest flags: %s (target=%s)", res.Error, target)
		return nil, model.FailRetVal(i18n.Model.TeamFlag.GetError, "TeamFlag", target, res.Error)
	}
	return results, model.SuccessRetVal()
}

// CountSolvedForChallenge 统计队伍在比赛题目上已解出的生效 flag 数量。
// 归档（active = false）的 flag 不计入，但历史得分仍保留在 team_flags 中。
func (t *TeamFlagRepo) CountSolvedForChallenge(teamID uint, contestChallengeID uint) (int64, model.RetVal) {
	var count int64
	res := t.DB.Model(&model.TeamFlag{}).
		Joins("INNER JOIN contest_flags ON contest_flags.id = team_flags.contest_flag_id AND contest_flags.deleted_at IS NULL").
		Where("team_flags.team_id = ? AND contest_flags.contest_challenge_id = ? AND team_flags.solved = true AND contest_flags.active = true", teamID, contestChallengeID).
		Count(&count)
	target := fmt.Sprintf("contest_challenge_id=%d,team_id=%d", contestChallengeID, teamID)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count solved team flags for challenge: %s (target=%s)", res.Error, target)
		return 0, model.FailRetVal(i18n.Model.TeamFlag.GetError, "TeamFlag", target, res.Error)
	}
	return count, model.SuccessRetVal()
}

// CountSolvedTeamsForContestChallenge 统计已解出该比赛题目的队伍数。
// 只统计生效（active = true）flag 的解题记录，归档代次不计入。
func (t *TeamFlagRepo) CountSolvedTeamsForContestChallenge(contestChallengeID uint) (int64, model.RetVal) {
	var count int64
	res := t.DB.Model(&model.TeamFlag{}).
		Joins("INNER JOIN contest_flags ON contest_flags.id = team_flags.contest_flag_id AND contest_flags.deleted_at IS NULL").
		Where("contest_flags.contest_challenge_id = ? AND team_flags.solved = true AND contest_flags.active = true", contestChallengeID).
		Distinct("team_flags.team_id").
		Count(&count)
	target := fmt.Sprintf("contest_challenge_id=%d", contestChallengeID)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count solved teams for contest challenge: %s (target=%s)", res.Error, target)
		return 0, model.FailRetVal(i18n.Model.TeamFlag.GetError, "TeamFlag", target, res.Error)
	}
	return count, model.SuccessRetVal()
}

func (t *TeamFlagRepo) CountGenerated(teamID uint, contestFlagIDL ...uint) (int64, model.RetVal) {
	if len(contestFlagIDL) == 0 {
		return 0, model.SuccessRetVal()
	}
	var count int64
	res := t.DB.Model(&model.TeamFlag{}).
		Where("team_id = ? AND contest_flag_id IN ?", teamID, contestFlagIDL).
		Count(&count)
	target := fmt.Sprintf("contest_flag_id=%v,team_id=%d", contestFlagIDL, teamID)
	if res.Error != nil {
		log.Logger.Warningf("Failed to count generated team flags: %s (target=%s)", res.Error, target)
		return 0, model.FailRetVal(i18n.Model.TeamFlag.GetError, "TeamFlag", target, res.Error)
	}
	return count, model.SuccessRetVal()
}

func (t *TeamFlagRepo) GetByContestChallengeAndValue(teamID uint, contestChallengeID uint, value string) (model.TeamFlag, model.RetVal) {
	var teamFlag model.TeamFlag
	res := t.DB.Model(&model.TeamFlag{}).
		Joins("INNER JOIN contest_flags ON contest_flags.id = team_flags.contest_flag_id AND contest_flags.deleted_at IS NULL").
		Where("team_flags.team_id = ? AND contest_flags.contest_challenge_id = ? AND team_flags.value = ?", teamID, contestChallengeID, value).
		Limit(1).
		Find(&teamFlag)
	target := fmt.Sprintf("contest_challenge_id=%d,team_id=%d,value=%s", contestChallengeID, teamID, value)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get team flag by contest challenge and value: %s (target=%s)", res.Error, target)
		return model.TeamFlag{}, model.FailRetVal(i18n.Model.TeamFlag.GetError, "TeamFlag", target, res.Error)
	}
	if res.RowsAffected == 0 {
		return model.TeamFlag{}, model.FailRetVal(i18n.Model.TeamFlag.NotFound, "TeamFlag", target, nil)
	}
	return teamFlag, model.SuccessRetVal()
}
