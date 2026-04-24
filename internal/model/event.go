package model

const (
	SkipEventType = "skip"

	LoginEventType      = "login"
	RegisterEventType   = "register"
	OauthLoginEventType = "oauth_login"

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

	UpdateCheatEventType      = "update_cheat"
	DeleteCheatEventType      = "delete_cheat"
	DeleteAllCheatEventType   = "delete_all_cheat"
	ManualCheckCheatEventType = "manual_check_cheat"

	CreateNoticeEventType = "create_notice"
	UpdateNoticeEventType = "update_notice"
	DeleteNoticeEventType = "delete_notice"

	CreateOauthEventType = "create_oauth"
	UpdateOauthEventType = "update_oauth"
	DeleteOauthEventType = "delete_oauth"

	CreateSmtpEventType = "create_smtp"
	UpdateSmtpEventType = "update_smtp"
	DeleteSmtpEventType = "delete_smtp"

	UpdateCronJobEventType = "update_cronjob"

	CreateWebhookEventType = "create_webhook"
	UpdateWebhookEventType = "update_webhook"
	DeleteWebhookEventType = "delete_webhook"

	UpdateBrandingEventType = "update_branding"
	UpdateSettingEventType  = "update_setting"
	RestartSystemEventType  = "restart_system"

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
	ExtendVictimEventType       = "extend_victim"
	StopVictimEventType         = "stop_victim"
	DownloadTrafficEventType    = "download_traffic"
	ReadFlagEventType           = "read_flag"
	SubmitFlagEventType         = "submit_flag"

	StartGeneratorEventType = "start_generator"
	StopGeneratorEventType  = "stop_generator"

	UploadWriteUpEventType   = "upload_writeup"
	DownloadWriteUpEventType = "download_writeup"

	DownloadFileEventType = "download_file"

	CreateRoleEventType = "create_role"
	UpdateRoleEventType = "update_role"
	DeleteRoleEventType = "delete_role"

	CreateGroupEventType = "create_group"
	UpdateGroupEventType = "update_group"
	DeleteGroupEventType = "delete_group"

	UpdatePermissionEventType = "update_permission"
	AssignPermissionEventType = "assign_permission"
	RevokePermissionEventType = "revoke_permission"

	AssignUserGroupEventType = "assign_user_group"
	RemoveUserGroupEventType = "remove_user_group"
)

var EventTypes = []string{
	LoginEventType, RegisterEventType, OauthLoginEventType,
	CreateUserEventType, UpdateUserEventType, DeleteUserEventType,
	CreateContestEventType, UpdateContestEventType, DeleteContestEventType,
	CreateChallengeEventType, UpdateChallengeEventType, DeleteChallengeEventType, UploadChallengeFileEventType,
	CreateContestChallengeEventType, UpdateContestChallengeEventType, DeleteContestChallengeEventType,
	UpdateContestChallengeFlagEventType,
	UpdateCheatEventType, DeleteCheatEventType, DeleteAllCheatEventType, ManualCheckCheatEventType,
	CreateNoticeEventType, UpdateNoticeEventType, DeleteNoticeEventType,
	CreateOauthEventType, UpdateOauthEventType, DeleteOauthEventType,
	CreateSmtpEventType, UpdateSmtpEventType, DeleteSmtpEventType,
	UpdateCronJobEventType,
	CreateWebhookEventType, UpdateWebhookEventType, DeleteWebhookEventType,
	UpdateBrandingEventType,
	UpdateSettingEventType, RestartSystemEventType,
	ActivateEmailEventType, VerifyEmailEventType,
	UploadPictureEventType, DeletePictureEventType,
	JoinTeamEventType, CreateTeamEventType, UpdateTeamEventType, DeleteTeamEventType,
	LeaveTeamEventType, KickMemberEventType,
	InitChallengeEventType, ResetChallengeEventType,
	DownloadAttachmentEventType, PullImageEventType, StartVictimEventType, ExtendVictimEventType, StopVictimEventType,
	DownloadTrafficEventType,
	ReadFlagEventType, SubmitFlagEventType,
	StartGeneratorEventType, StopGeneratorEventType,
	UploadWriteUpEventType, DownloadWriteUpEventType,
	DownloadFileEventType,
	CreateRoleEventType, UpdateRoleEventType, DeleteRoleEventType,
	CreateGroupEventType, UpdateGroupEventType, DeleteGroupEventType,
	UpdatePermissionEventType, AssignPermissionEventType, RevokePermissionEventType,
	AssignUserGroupEventType, RemoveUserGroupEventType,
}

type Event struct {
	Type    string  `json:"type"`
	Success bool    `json:"success"`
	IP      string  `json:"ip"`
	Magic   string  `json:"magic"`
	Models  UintMap `gorm:"type:jsonb" json:"models"`
	BaseModel
}
