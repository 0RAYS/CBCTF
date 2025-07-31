package model

import "CBCTF/internal/i18n"

// TeamFlag
// BelongsTo Team
// BelongsTo ContestFlag
type TeamFlag struct {
	TeamID          uint          `json:"team_id"`
	Team            Team          `json:"-"`
	ContestFlagID   uint          `json:"contest_flag_id"`
	ContestFlag     ContestFlag   `json:"-"`
	ChallengeFlagID uint          `json:"challenge_flag_id"`
	ChallengeFlag   ChallengeFlag `json:"-"`
	Value           string        `json:"value"`
	Solved          bool          `json:"solved"`
	BasicModel
}

func (t TeamFlag) GetModelName() string {
	return "TeamFlag"
}

func (t TeamFlag) GetVersion() uint {
	return t.Version
}

func (t TeamFlag) CreateErrorString() string {
	return i18n.CreateTeamFlagError
}

func (t TeamFlag) DeleteErrorString() string {
	return i18n.DeleteTeamFlagError
}

func (t TeamFlag) GetErrorString() string {
	return i18n.GetTeamFlagError
}

func (t TeamFlag) NotFoundErrorString() string {
	return i18n.TeamFlagNotFound
}

func (t TeamFlag) UpdateErrorString() string {
	return i18n.UpdateTeamFlagError
}

func (t TeamFlag) GetUniqueKey() []string {
	return []string{"id"}
}
