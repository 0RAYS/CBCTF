package model

type Submission struct {
	ChallengeID     uint          `gorm:"index:idx_submissions_team_challenge_solved,priority:2;index:idx_submissions_team_challenge_unsolved_active,priority:2,where:deleted_at IS NULL AND solved = false;index" json:"challenge_id"`
	Challenge       Challenge     `json:"-"`
	ContestID       uint          `gorm:"index" json:"contest_id"`
	Contest         Contest       `json:"-"`
	TeamID          uint          `gorm:"index:idx_submissions_team_challenge_solved,priority:1;index:idx_submissions_team_challenge_unsolved_active,priority:1,where:deleted_at IS NULL AND solved = false;index:idx_submissions_challenge_flag_team_solved_active,priority:2,where:deleted_at IS NULL AND solved = true AND team_id <> 0;index:idx_submissions_team_solved_created_active,priority:1,where:deleted_at IS NULL AND solved = true;index" json:"team_id"`
	Team            Team          `json:"-"`
	UserID          uint          `gorm:"index;index:idx_submissions_user_challenge_flag_solved_active,priority:1,where:deleted_at IS NULL AND solved = true" json:"user_id"`
	User            User          `json:"-"`
	ChallengeFlagID uint          `gorm:"index:idx_submissions_challenge_flag_solved,priority:1;index:idx_submissions_challenge_flag_team_solved_active,priority:1,where:deleted_at IS NULL AND solved = true AND team_id <> 0;index:idx_submissions_user_challenge_flag_solved_active,priority:2,where:deleted_at IS NULL AND solved = true;index" json:"challenge_flag_id"`
	ChallengeFlag   ChallengeFlag `json:"-"`
	Value           string        `json:"value"`
	Solved          bool          `gorm:"index:idx_submissions_team_challenge_solved,priority:3;index:idx_submissions_challenge_flag_solved,priority:2;index:idx_submissions_team_challenge_unsolved_active,priority:3,where:deleted_at IS NULL AND solved = false;index:idx_submissions_challenge_flag_team_solved_active,priority:3,where:deleted_at IS NULL AND solved = true AND team_id <> 0;index:idx_submissions_team_solved_created_active,priority:2,where:deleted_at IS NULL AND solved = true;index" json:"solved"`
	Score           float64       `gorm:"default:0" json:"score"`
	IP              string        `json:"ip"`
	BaseModel
}
