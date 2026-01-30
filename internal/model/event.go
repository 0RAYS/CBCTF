package model

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

	UpdateCheatEventType = "update_cheat"

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

	UploadPictureEventType = "upload_picture"
	DeletePictureEventType = "delete_picture"

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

	DownloadFileEventType = "download_file"
)

var EventTypes = []string{
	LoginEventType, RegisterEventType, OauthLoginEventType, CreateAdminEventType, UpdateAdminEventType,
	CreateUserEventType, UpdateUserEventType, DeleteUserEventType, CreateContestEventType, UpdateContestEventType,
	DeleteContestEventType, CreateChallengeEventType, UpdateChallengeEventType, DeleteChallengeEventType,
	UploadChallengeFileEventType, CreateContestChallengeEventType, UpdateContestChallengeEventType,
	DeleteContestChallengeEventType, UpdateContestChallengeFlagEventType, UpdateCheatEventType, CreateNoticeEventType,
	UpdateNoticeEventType, DeleteNoticeEventType, CreateOauthEventType, UpdateOauthEventType, DeleteOauthEventType,
	CreateSmtpEventType, UpdateSmtpEventType, DeleteSmtpEventType, CreateWebhookEventType, UpdateWebhookEventType,
	DeleteWebhookEventType, ActivateEmailEventType, VerifyEmailEventType, UploadPictureEventType, DeletePictureEventType,
	JoinTeamEventType, CreateTeamEventType, UpdateTeamEventType, DeleteTeamEventType, LeaveTeamEventType,
	KickMemberEventType, InitChallengeEventType, ResetChallengeEventType, DownloadAttachmentEventType, PullImageEventType,
	StartVictimEventType, IncreaseVictimEventType, StopVictimEventType, DownloadTrafficEventType, SubmitFlagEventType,
	UploadWriteUpEventType, DownloadWriteUpEventType, DownloadFileEventType,
}

type Event struct {
	IsAdmin bool    `json:"is_admin"`
	Type    string  `json:"type"`
	Success bool    `json:"success"`
	IP      string  `json:"ip"`
	Magic   string  `json:"magic"`
	Models  UintMap `gorm:"type:json" json:"models"`
	BaseModel
}

func (e Event) GetModelName() string {
	return "Event"
}

func (e Event) GetBaseModel() BaseModel {
	return e.BaseModel
}

func (e Event) GetUniqueKey() []string {
	return []string{"id"}
}

func (e Event) GetAllowedQueryFields() []string {
	return []string{}
}
