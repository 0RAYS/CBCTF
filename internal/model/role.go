package model

const (
	AdminRoleName     = "admin"
	OrganizerRoleName = "organizer"
	UserRoleName      = "user"
)

var DefaultRoles = []Role{
	{Name: AdminRoleName, Description: "系统管理员, 拥有全部权限", Default: true},
	{Name: OrganizerRoleName, Description: "赛事主办方, 拥有赛事相关管理权限", Default: true},
	{Name: UserRoleName, Description: "参赛选手, 拥有参赛相关权限", Default: true},
}

var DefaultRolePermissionMap = map[string][]string{
	AdminRoleName: {
		PermSelfRead, PermSelfUpdate, PermSelfDelete, PermSelfActivate,
		PermUserContestRead, PermUserContestRank,
		PermUserTeamCreate, PermUserTeamJoin, PermUserTeamRead, PermUserTeamUpdate, PermUserTeamDelete,
		PermUserNoticeList,
		PermUserChallengeList, PermUserChallengeRead, PermUserChallengeInit, PermUserChallengeReset, PermUserChallengeSubmit,
		PermUserVictimControl,
		PermUserWriteupUpload, PermUserWriteupList,

		PermAdminIPSearch, PermAdminModelsSearch,
		PermAdminSystemStatus, PermAdminSystemRead, PermAdminSystemUpdate, PermAdminSystemRestart,
		PermAdminPermissionUpdate, PermAdminPermissionList,
		PermAdminRoleCreate, PermAdminRoleRead, PermAdminRoleUpdate, PermAdminRoleDelete, PermAdminRoleList, PermAdminRoleAssign, PermAdminRoleRevoke,
		PermAdminGroupCreate, PermAdminGroupRead, PermAdminGroupUpdate, PermAdminGroupDelete, PermAdminGroupList,
		PermAdminUserCreate, PermAdminUserRead, PermAdminUserUpdate, PermAdminUserDelete, PermAdminUserList, PermAdminUserAssign, PermAdminUserRevoke,
		PermAdminOauthCreate, PermAdminOauthRead, PermAdminOauthUpdate, PermAdminOauthDelete, PermAdminOauthList,
		PermAdminSMTPCreate, PermAdminSMTPRead, PermAdminSMTPUpdate, PermAdminSMTPDelete, PermAdminSMTPList,
		PermAdminCronJobList, PermAdminCronJobUpdate,
		PermAdminWebhookCreate, PermAdminWebhookRead, PermAdminWebhookUpdate, PermAdminWebhookDelete, PermAdminWebhookList,
		PermAdminChallengeCreate, PermAdminChallengeRead, PermAdminChallengeUpdate, PermAdminChallengeDelete, PermAdminChallengeList, PermAdminChallengeTest,
		PermAdminContestCreate, PermAdminContestRead, PermAdminContestUpdate, PermAdminContestDelete, PermAdminContestList, PermAdminContestRank,
		PermAdminTeamRead, PermAdminTeamUpdate, PermAdminTeamDelete, PermAdminTeamList,
		PermAdminTeamWriteupList, PermAdminTeamWriteupRead,
		PermAdminNoticeCreate, PermAdminNoticeUpdate, PermAdminNoticeDelete, PermAdminNoticeList,
		PermAdminCheatCreate, PermAdminCheatUpdate, PermAdminCheatDelete, PermAdminCheatList,
		PermAdminContestChallengeCreate, PermAdminContestChallengeRead, PermAdminContestChallengeUpdate, PermAdminContestChallengeDelete, PermAdminContestChallengeList,
		PermAdminContestChallengeFlagList, PermAdminContestChallengeFlagRead, PermAdminContestChallengeFlagUpdate,
		PermAdminImagePull,
		PermAdminVictimControl,
		PermAdminGeneratorControl,
		PermAdminFileList, PermAdminFileRead, PermAdminFileDelete,
		PermAdminLogRead,
	},
	OrganizerRoleName: {
		PermSelfRead, PermSelfUpdate, PermSelfDelete, PermSelfActivate,
		PermUserContestRead, PermUserContestRank,
		PermUserTeamCreate, PermUserTeamJoin, PermUserTeamRead, PermUserTeamUpdate, PermUserTeamDelete,
		PermUserNoticeList,
		PermUserChallengeList, PermUserChallengeRead, PermUserChallengeInit, PermUserChallengeReset, PermUserChallengeSubmit,
		PermUserVictimControl,
		PermUserWriteupUpload, PermUserWriteupList,

		PermAdminIPSearch,
		PermAdminSystemStatus,
		PermAdminChallengeCreate, PermAdminChallengeRead, PermAdminChallengeUpdate, PermAdminChallengeDelete, PermAdminChallengeList, PermAdminChallengeTest,
		PermAdminContestCreate, PermAdminContestRead, PermAdminContestUpdate, PermAdminContestDelete, PermAdminContestList, PermAdminContestRank,
		PermAdminTeamRead, PermAdminTeamUpdate, PermAdminTeamDelete, PermAdminTeamList,
		PermAdminTeamWriteupList, PermAdminTeamWriteupRead,
		PermAdminNoticeCreate, PermAdminNoticeUpdate, PermAdminNoticeDelete, PermAdminNoticeList,
		PermAdminCheatCreate, PermAdminCheatUpdate, PermAdminCheatDelete, PermAdminCheatList,
		PermAdminContestChallengeCreate, PermAdminContestChallengeRead, PermAdminContestChallengeUpdate, PermAdminContestChallengeDelete, PermAdminContestChallengeList,
		PermAdminContestChallengeFlagList, PermAdminContestChallengeFlagRead, PermAdminContestChallengeFlagUpdate,
		PermAdminImagePull,
		PermAdminVictimControl,
		PermAdminGeneratorControl,
	},
	UserRoleName: {
		PermSelfRead, PermSelfUpdate, PermSelfDelete, PermSelfActivate,
		PermUserContestRead, PermUserContestRank,
		PermUserTeamCreate, PermUserTeamJoin, PermUserTeamRead, PermUserTeamUpdate, PermUserTeamDelete,
		PermUserNoticeList,
		PermUserChallengeList, PermUserChallengeRead, PermUserChallengeInit, PermUserChallengeReset, PermUserChallengeSubmit,
		PermUserVictimControl,
		PermUserWriteupUpload, PermUserWriteupList,
	},
}

// Role 角色
// ManyToMany Permission
type Role struct {
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"-"`
	Name        string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Description string       `json:"description"`
	Default     bool         `json:"default"`
	BaseModel
}
