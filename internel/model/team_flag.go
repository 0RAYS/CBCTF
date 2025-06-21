package model

import "CBCTF/internel/i18n"

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

func (r TeamFlag) GetModelName() string {
	return "TeamFlag"
}

func (r TeamFlag) GetVersion() uint {
	return r.Version
}

func (r TeamFlag) CreateErrorString() string {
	return i18n.CreateTeamFlagError
}

func (r TeamFlag) DeleteErrorString() string {
	return i18n.DeleteTeamFlagError
}

func (r TeamFlag) GetErrorString() string {
	return i18n.GetTeamFlagError
}

func (r TeamFlag) NotFoundErrorString() string {
	return i18n.TeamFlagNotFound
}

func (r TeamFlag) UpdateErrorString() string {
	return i18n.UpdateTeamFlagError
}

func (r TeamFlag) GetUniqueKey() []string {
	return []string{"id"}
}
