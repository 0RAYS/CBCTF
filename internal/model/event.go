package model

import (
	"CBCTF/internal/i18n"
)

const (
	UserLoginEventType          = "user_login"
	UserRegisterEventType       = "user_register"
	UserUpdateEventType         = "user_update"
	UserUpdatePasswordEventType = "user_update_password"
	UserVerifyEmailEventType    = "user_verify"
	UserDeleteEventType         = "user_delete"
	JoinTeamEventType           = "join_team"
	CreateTeamEventType         = "create_team"
	UpdateTeamEventType         = "update_team"
	LeaveTeamEventType          = "leave_team"
	KickMemberEventType         = "kick_member"
	InitChallengeEventType      = "init_usage"
	ResetChallengeEventType     = "reset_usage"
	DownloadAttachmentEventType = "download_attachment"
	StartVictimEventType        = "start_victim"
	IncreaseVictimEventType     = "increase_victim"
	StopVictimEventType         = "stop_victim"
	SubmitFlagEventType         = "submit_flag"
	UploadWriteUpEventType      = "upload_writeup"
)

type Event struct {
	UserID             *uint             `gorm:"default:null" json:"user_id"`
	User               *User             `json:"-"`
	TeamID             *uint             `gorm:"default:null" json:"team_id"`
	Team               *Team             `json:"-"`
	ContestID          *uint             `gorm:"default:null" json:"contest_id"`
	Contest            *Contest          `json:"-"`
	ContestChallengeID *uint             `gorm:"default:null" json:"contest_challenge_id"`
	ContestChallenge   *ContestChallenge `json:"-"`
	Desc               string            `json:"desc"`
	Type               string            `json:"type"`
	IP                 string            `json:"ip"`
	Magic              string            `json:"magic"`
	BasicModel
}

func (e Event) GetModelName() string {
	return "Event"
}

func (e Event) GetVersion() uint {
	return e.Version
}

func (e Event) CreateErrorString() string {
	return i18n.CreateEventError
}

func (e Event) DeleteErrorString() string {
	return i18n.DeleteEventError
}

func (e Event) GetErrorString() string {
	return i18n.GetEventError
}

func (e Event) NotFoundErrorString() string {
	return i18n.EventNotFound
}

func (e Event) UpdateErrorString() string {
	return i18n.UpdateEventError
}

func (e Event) GetUniqueKey() []string {
	return []string{"id"}
}

func (e Event) GetForeignKeys() []string {
	return []string{"id", "user_id", "team_id", "contest_id", "contest_challenge_id"}
}
