package model

import (
	"gorm.io/gorm"
	"time"
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
	ID         uint           `gorm:"primarykey" json:"id"`
	References References     `gorm:"type:json" json:"references"`
	Desc       string         `json:"desc"`
	Type       string         `json:"type"`
	IP         string         `json:"ip"`
	Magic      string         `json:"magic"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Version    uint           `gorm:"default:1" json:"-"`
}
