package model

import (
	"CBCTF/internal/i18n"
)

const (
	SkipEventType = "skip"

	LoginEventType      = "login"
	RegisterEventType   = "register"
	OauthLoginEventType = "oauth_login"

	CreateAdminEventType = "create_admin"
	UpdateAdminEventType = "update_admin"

	CreateUserEventType = "create_user"
	UpdateUserEventType = "update_user"
	DeleteUserEventType = "delete_user"

	CreateContestEventType = "create_contest"
	UpdateContestEventType = "update_contest"
	DeleteContestEventType = "delete_contest"

	CreateChallengeEventType     = "create_challenge"
	UpdateChallengeEventType     = "update_challenge"
	DeleteChallengeEventType     = "delete_challenge"
	UploadChallengeFileEventType = "upload_challenge_file"

	CreateContestChallengeEventType = "create_contest_challenge"
	UpdateContestChallengeEventType = "update_contest_challenge"
	DeleteContestChallengeEventType = "delete_contest_challenge"

	UpdateContestChallengeFlagEventType = "update_contest_challenge_flag"

	CreateNoticeEventType = "create_notice"
	UpdateNoticeEventType = "update_notice"
	DeleteNoticeEventType = "delete_notice"

	CreateOauthEventType = "create_oauth"
	UpdateOauthEventType = "update_oauth"
	DeleteOauthEventType = "delete_oauth"

	CreateSmtpEventType = "create_smtp"
	UpdateSmtpEventType = "update_smtp"
	DeleteSmtpEventType = "delete_smtp"

	CreateWebhookEventType = "create_webhook"
	UpdateWebhookEventType = "update_webhook"
	DeleteWebhookEventType = "delete_webhook"

	ActivateEmailEventType = "activate_email"
	VerifyEmailEventType   = "verify_email"

	UploadAvatarEventType = "upload_avatar"
	DeleteAvatarEventType = "delete_avatar"

	JoinTeamEventType   = "join_team"
	CreateTeamEventType = "create_team"
	UpdateTeamEventType = "update_team"
	DeleteTeamEventType = "delete_team"
	LeaveTeamEventType  = "leave_team"
	KickMemberEventType = "kick_member"

	InitChallengeEventType      = "init_challenge"
	ResetChallengeEventType     = "reset_challenge"
	DownloadAttachmentEventType = "download_attachment"
	PullImageEventType          = "pull_image"
	StartVictimEventType        = "start_victim"
	IncreaseVictimEventType     = "increase_victim"
	StopVictimEventType         = "stop_victim"
	DownloadTrafficEventType    = "download_traffic"
	SubmitFlagEventType         = "submit_flag"

	UploadWriteUpEventType   = "upload_writeup"
	DownloadWriteUpEventType = "download_writeup"
)

var EventTypes = []string{
	LoginEventType, RegisterEventType, OauthLoginEventType, CreateAdminEventType, UpdateAdminEventType,
	CreateUserEventType, UpdateUserEventType, DeleteUserEventType, CreateContestEventType, UpdateContestEventType,
	DeleteContestEventType, CreateChallengeEventType, UpdateChallengeEventType, DeleteChallengeEventType,
	UploadChallengeFileEventType, CreateContestChallengeEventType, UpdateContestChallengeEventType,
	DeleteContestChallengeEventType, UpdateContestChallengeFlagEventType, CreateNoticeEventType, UpdateNoticeEventType,
	DeleteNoticeEventType, CreateOauthEventType, UpdateOauthEventType, DeleteOauthEventType, CreateSmtpEventType,
	UpdateSmtpEventType, DeleteSmtpEventType, CreateWebhookEventType, UpdateWebhookEventType, DeleteWebhookEventType,
	ActivateEmailEventType, VerifyEmailEventType, UploadAvatarEventType, DeleteAvatarEventType, JoinTeamEventType,
	CreateTeamEventType, UpdateTeamEventType, DeleteTeamEventType, LeaveTeamEventType, KickMemberEventType,
	InitChallengeEventType, ResetChallengeEventType, DownloadAttachmentEventType, PullImageEventType,
	StartVictimEventType, IncreaseVictimEventType, StopVictimEventType, DownloadTrafficEventType, SubmitFlagEventType,
	UploadWriteUpEventType, DownloadWriteUpEventType,
}

type Event struct {
	WebhookHistories []WebhookHistory `json:"-"`
	IsAdmin          bool             `json:"is_admin"`
	Type             string           `json:"type"`
	Success          bool             `json:"success"`
	IP               string           `json:"ip"`
	Magic            string           `json:"magic"`
	Models           UintMap          `gorm:"type:json" json:"models"`
	BasicModel
}

func (e Event) GetModelName() string {
	return "Event"
}

func (e Event) GetVersion() uint {
	return e.Version
}

func (e Event) GetBasicModel() BasicModel {
	return e.BasicModel
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
