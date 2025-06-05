package model

import (
	"CBCTF/internel/i18n"
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	Desc      string    `json:"desc"`
	Type      string    `json:"type"`
	IP        string    `json:"ip"`
	Magic     string    `json:"magic"`
	Reference Reference `gorm:"type:json" json:"reference"`
	Basic
}

type Reference struct {
	UserID    uint `json:"user_id"`
	TeamID    uint `json:"team_id"`
	ContestID uint `json:"contest_id"`
	UsageID   uint `json:"usage_id"`
}

func (r Reference) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Reference) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Reference value")
	}
	return json.Unmarshal(bytes, r)
}

func (e Event) GetModelName() string {
	return "Event"
}

func (e Event) GetID() uint {
	return e.ID
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
