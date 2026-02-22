package model

const (
	PermSelfRead     = "self:read"
	PermSelfUpdate   = "self:update"
	PermSelfDelete   = "self:delete"
	PermSelfActivate = "self:activate"

	PermUserContestRead = "user:contest:read"
	PermUserContestRank = "user:contest:rank"

	PermUserTeamCreate = "user:team:create"
	PermUserTeamJoin   = "user:team:join"
	PermUserTeamRead   = "user:team:read"
	PermUserTeamUpdate = "user:team:update"
	PermUserTeamDelete = "user:team:delete"

	PermUserNoticeList = "user:notice:list"

	PermUserChallengeList   = "user:challenge:list"
	PermUserChallengeRead   = "user:challenge:read"
	PermUserChallengeInit   = "user:challenge:init"
	PermUserChallengeReset  = "user:challenge:reset"
	PermUserChallengeSubmit = "user:challenge:submit"

	PermUserVictimControl = "user:victim:control"

	PermUserWriteupUpload = "user:writeup:upload"
	PermUserWriteupList   = "user:writeup:list"

	PermAdminIPSearch     = "admin:ip:search"
	PermAdminModelsSearch = "admin:models:search"

	PermAdminSystemRead    = "admin:system:read"
	PermAdminSystemUpdate  = "admin:system:update"
	PermAdminSystemRestart = "admin:system:restart"

	PermAdminPermissionRead   = "admin:permission:read"
	PermAdminPermissionUpdate = "admin:permission:update"
	PermAdminPermissionList   = "admin:permission:list"

	PermAdminRoleCreate = "admin:role:create"
	PermAdminRoleRead   = "admin:role:read"
	PermAdminRoleUpdate = "admin:role:update"
	PermAdminRoleDelete = "admin:role:delete"
	PermAdminRoleList   = "admin:role:list"
	PermAdminRoleAssign = "admin:role:assign"
	PermAdminRoleRevoke = "admin:role:revoke"

	PermAdminGroupCreate = "admin:group:create"
	PermAdminGroupRead   = "admin:group:read"
	PermAdminGroupUpdate = "admin:group:update"
	PermAdminGroupDelete = "admin:group:delete"
	PermAdminGroupList   = "admin:group:list"

	PermAdminUserCreate = "admin:user:create"
	PermAdminUserRead   = "admin:user:read"
	PermAdminUserUpdate = "admin:user:update"
	PermAdminUserDelete = "admin:user:delete"
	PermAdminUserList   = "admin:user:list"
	PermAdminUserAssign = "admin:user:assign"
	PermAdminUserRevoke = "admin:user:revoke"

	PermAdminOauthCreate = "admin:oauth:create"
	PermAdminOauthRead   = "admin:oauth:read"
	PermAdminOauthUpdate = "admin:oauth:update"
	PermAdminOauthDelete = "admin:oauth:delete"
	PermAdminOauthList   = "admin:oauth:list"

	PermAdminSMTPCreate = "admin:smtp:create"
	PermAdminSMTPRead   = "admin:smtp:read"
	PermAdminSMTPUpdate = "admin:smtp:update"
	PermAdminSMTPDelete = "admin:smtp:delete"
	PermAdminSMTPList   = "admin:smtp:list"

	PermAdminWebhookCreate = "admin:webhook:create"
	PermAdminWebhookRead   = "admin:webhook:read"
	PermAdminWebhookUpdate = "admin:webhook:update"
	PermAdminWebhookDelete = "admin:webhook:delete"
	PermAdminWebhookList   = "admin:webhook:list"

	PermAdminChallengeCreate = "admin:challenge:create"
	PermAdminChallengeRead   = "admin:challenge:read"
	PermAdminChallengeUpdate = "admin:challenge:update"
	PermAdminChallengeDelete = "admin:challenge:delete"
	PermAdminChallengeList   = "admin:challenge:list"
	PermAdminChallengeTest   = "admin:challenge:test"

	PermAdminContestCreate = "admin:contest:create"
	PermAdminContestRead   = "admin:contest:read"
	PermAdminContestUpdate = "admin:contest:update"
	PermAdminContestDelete = "admin:contest:delete"
	PermAdminContestList   = "admin:contest:list"
	PermAdminContestRank   = "admin:contest:rank"

	PermAdminTeamRead   = "admin:team:read"
	PermAdminTeamUpdate = "admin:team:update"
	PermAdminTeamDelete = "admin:team:delete"
	PermAdminTeamList   = "admin:team:list"

	PermAdminTeamWriteupList = "admin:team_writeup:list"
	PermAdminTeamWriteupRead = "admin:team_writeup:read"

	PermAdminNoticeCreate = "admin:notice:create"
	PermAdminNoticeUpdate = "admin:notice:update"
	PermAdminNoticeDelete = "admin:notice:delete"
	PermAdminNoticeList   = "admin:notice:list"

	PermAdminCheatCreate = "admin:cheat:create"
	PermAdminCheatUpdate = "admin:cheat:update"
	PermAdminCheatDelete = "admin:cheat:delete"
	PermAdminCheatList   = "admin:cheat:list"

	PermAdminContestChallengeCreate = "admin:contest_challenge:create"
	PermAdminContestChallengeRead   = "admin:contest_challenge:read"
	PermAdminContestChallengeUpdate = "admin:contest_challenge:update"
	PermAdminContestChallengeDelete = "admin:contest_challenge:delete"
	PermAdminContestChallengeList   = "admin:contest_challenge:list"

	PermAdminContestChallengeFlagList   = "admin:contest_challenge_flag:list"
	PermAdminContestChallengeFlagRead   = "admin:contest_challenge_flag:read"
	PermAdminContestChallengeFlagUpdate = "admin:contest_challenge_flag:update"

	PermAdminImagePull = "admin:image:pull"

	PermAdminVictimControl = "admin:victim:control"

	PermAdminFileList   = "admin:file:list"
	PermAdminFileRead   = "admin:file:read"
	PermAdminFileDelete = "admin:file:delete"

	PermAdminLogRead = "admin:log:read"
)

var Permissions = []Permission{
	{Name: PermSelfRead, Resource: "self", Operation: "read", Description: "查看自身信息"},
	{Name: PermSelfUpdate, Resource: "self", Operation: "update", Description: "更新自身信息"},
	{Name: PermSelfDelete, Resource: "self", Operation: "delete", Description: "删除自身账号"},
	{Name: PermSelfActivate, Resource: "self", Operation: "activate", Description: "激活自身账号"},

	{Name: PermUserContestRead, Resource: "user:contest", Operation: "read", Description: "查看比赛详情"},
	{Name: PermUserContestRank, Resource: "user:contest", Operation: "rank", Description: "查看比赛排名"},

	{Name: PermUserTeamCreate, Resource: "user:team", Operation: "create", Description: "创建队伍"},
	{Name: PermUserTeamJoin, Resource: "user:team", Operation: "join", Description: "加入队伍"},
	{Name: PermUserTeamRead, Resource: "user:team", Operation: "read", Description: "查看队伍详情"},
	{Name: PermUserTeamUpdate, Resource: "user:team", Operation: "update", Description: "更新队伍"},
	{Name: PermUserTeamDelete, Resource: "user:team", Operation: "delete", Description: "删除队伍"},

	{Name: PermUserNoticeList, Resource: "user:notice", Operation: "list", Description: "查看公告列表"},

	{Name: PermUserChallengeList, Resource: "user:challenge", Operation: "list", Description: "查看题目列表"},
	{Name: PermUserChallengeRead, Resource: "user:challenge", Operation: "read", Description: "查看题目详情"},
	{Name: PermUserChallengeInit, Resource: "user:challenge", Operation: "init", Description: "初始化题目环境"},
	{Name: PermUserChallengeReset, Resource: "user:challenge", Operation: "reset", Description: "重置题目环境"},
	{Name: PermUserChallengeSubmit, Resource: "user:challenge", Operation: "submit", Description: "提交 Flag"},

	{Name: PermUserVictimControl, Resource: "user:victim", Operation: "control", Description: "控制靶机"},

	{Name: PermUserWriteupUpload, Resource: "user:writeup", Operation: "upload", Description: "上传 Writeup"},
	{Name: PermUserWriteupList, Resource: "user:writeup", Operation: "list", Description: "查看 Writeup 列表"},

	{Name: PermAdminPermissionRead, Resource: "admin:permission", Operation: "read", Description: "查看权限详情"},
	{Name: PermAdminPermissionUpdate, Resource: "admin:permission", Operation: "update", Description: "更新权限"},
	{Name: PermAdminPermissionList, Resource: "admin:permission", Operation: "list", Description: "查看权限列表"},

	{Name: PermAdminRoleCreate, Resource: "admin:role", Operation: "create", Description: "创建角色"},
	{Name: PermAdminRoleRead, Resource: "admin:role", Operation: "read", Description: "查看角色详情"},
	{Name: PermAdminRoleUpdate, Resource: "admin:role", Operation: "update", Description: "更新角色"},
	{Name: PermAdminRoleDelete, Resource: "admin:role", Operation: "delete", Description: "删除角色"},
	{Name: PermAdminRoleList, Resource: "admin:role", Operation: "list", Description: "查看角色列表"},
	{Name: PermAdminRoleAssign, Resource: "admin:role", Operation: "assign", Description: "分配角色"},
	{Name: PermAdminRoleRevoke, Resource: "admin:role", Operation: "revoke", Description: "移除角色"},

	{Name: PermAdminGroupCreate, Resource: "admin:group", Operation: "create", Description: "创建用户组"},
	{Name: PermAdminGroupRead, Resource: "admin:group", Operation: "read", Description: "查看用户组详情"},
	{Name: PermAdminGroupUpdate, Resource: "admin:group", Operation: "update", Description: "更新用户组"},
	{Name: PermAdminGroupDelete, Resource: "admin:group", Operation: "delete", Description: "删除用户组"},
	{Name: PermAdminGroupList, Resource: "admin:group", Operation: "list", Description: "查看用户组列表"},

	{Name: PermAdminUserCreate, Resource: "admin:user", Operation: "create", Description: "创建用户"},
	{Name: PermAdminUserRead, Resource: "admin:user", Operation: "read", Description: "查看用户详情"},
	{Name: PermAdminUserUpdate, Resource: "admin:user", Operation: "update", Description: "更新用户"},
	{Name: PermAdminUserDelete, Resource: "admin:user", Operation: "delete", Description: "删除用户"},
	{Name: PermAdminUserList, Resource: "admin:user", Operation: "list", Description: "查看用户列表"},
	{Name: PermAdminUserAssign, Resource: "admin:user", Operation: "assign", Description: "分配用户"},
	{Name: PermAdminUserRevoke, Resource: "admin:user", Operation: "revoke", Description: "移除用户"},

	{Name: PermAdminSystemRead, Resource: "admin:system", Operation: "read", Description: "查看系统配置"},
	{Name: PermAdminSystemUpdate, Resource: "admin:system", Operation: "update", Description: "更新系统配置"},
	{Name: PermAdminSystemRestart, Resource: "admin:system", Operation: "restart", Description: "重启系统"},

	{Name: PermAdminIPSearch, Resource: "admin:ip", Operation: "search", Description: "搜索 IP"},
	{Name: PermAdminModelsSearch, Resource: "admin:models", Operation: "search", Description: "搜索模型"},

	{Name: PermAdminOauthCreate, Resource: "admin:oauth", Operation: "create", Description: "创建 OAuth 配置"},
	{Name: PermAdminOauthRead, Resource: "admin:oauth", Operation: "read", Description: "查看 OAuth 配置"},
	{Name: PermAdminOauthUpdate, Resource: "admin:oauth", Operation: "update", Description: "更新 OAuth 配置"},
	{Name: PermAdminOauthDelete, Resource: "admin:oauth", Operation: "delete", Description: "删除 OAuth 配置"},
	{Name: PermAdminOauthList, Resource: "admin:oauth", Operation: "list", Description: "查看 OAuth 配置列表"},

	{Name: PermAdminSMTPCreate, Resource: "admin:smtp", Operation: "create", Description: "创建 SMTP 配置"},
	{Name: PermAdminSMTPRead, Resource: "admin:smtp", Operation: "read", Description: "查看 SMTP 配置"},
	{Name: PermAdminSMTPUpdate, Resource: "admin:smtp", Operation: "update", Description: "更新 SMTP 配置"},
	{Name: PermAdminSMTPDelete, Resource: "admin:smtp", Operation: "delete", Description: "删除 SMTP 配置"},
	{Name: PermAdminSMTPList, Resource: "admin:smtp", Operation: "list", Description: "查看 SMTP 配置列表"},

	{Name: PermAdminWebhookCreate, Resource: "admin:webhook", Operation: "create", Description: "创建 Webhook"},
	{Name: PermAdminWebhookRead, Resource: "admin:webhook", Operation: "read", Description: "查看 Webhook"},
	{Name: PermAdminWebhookUpdate, Resource: "admin:webhook", Operation: "update", Description: "更新 Webhook"},
	{Name: PermAdminWebhookDelete, Resource: "admin:webhook", Operation: "delete", Description: "删除 Webhook"},
	{Name: PermAdminWebhookList, Resource: "admin:webhook", Operation: "list", Description: "查看 Webhook 列表"},

	{Name: PermAdminChallengeCreate, Resource: "admin:challenge", Operation: "create", Description: "创建题目"},
	{Name: PermAdminChallengeRead, Resource: "admin:challenge", Operation: "read", Description: "查看题目详情"},
	{Name: PermAdminChallengeUpdate, Resource: "admin:challenge", Operation: "update", Description: "更新题目"},
	{Name: PermAdminChallengeDelete, Resource: "admin:challenge", Operation: "delete", Description: "删除题目"},
	{Name: PermAdminChallengeList, Resource: "admin:challenge", Operation: "list", Description: "查看题目列表"},
	{Name: PermAdminChallengeTest, Resource: "admin:challenge", Operation: "test", Description: "测试题目"},

	{Name: PermAdminContestCreate, Resource: "admin:contest", Operation: "create", Description: "创建比赛"},
	{Name: PermAdminContestRead, Resource: "admin:contest", Operation: "read", Description: "查看比赛详情"},
	{Name: PermAdminContestUpdate, Resource: "admin:contest", Operation: "update", Description: "更新比赛"},
	{Name: PermAdminContestDelete, Resource: "admin:contest", Operation: "delete", Description: "删除比赛"},
	{Name: PermAdminContestList, Resource: "admin:contest", Operation: "list", Description: "查看比赛列表"},
	{Name: PermAdminContestRank, Resource: "admin:contest", Operation: "rank", Description: "查看比赛排名"},

	{Name: PermAdminTeamRead, Resource: "admin:team", Operation: "read", Description: "查看队伍详情"},
	{Name: PermAdminTeamUpdate, Resource: "admin:team", Operation: "update", Description: "更新队伍"},
	{Name: PermAdminTeamDelete, Resource: "admin:team", Operation: "delete", Description: "删除队伍"},
	{Name: PermAdminTeamList, Resource: "admin:team", Operation: "list", Description: "查看队伍列表"},

	{Name: PermAdminTeamWriteupList, Resource: "admin:team_writeup", Operation: "list", Description: "查看队伍 Writeup 列表"},
	{Name: PermAdminTeamWriteupRead, Resource: "admin:team_writeup", Operation: "read", Description: "查看队伍 Writeup 详情"},

	{Name: PermAdminNoticeCreate, Resource: "admin:notice", Operation: "create", Description: "创建公告"},
	{Name: PermAdminNoticeUpdate, Resource: "admin:notice", Operation: "update", Description: "更新公告"},
	{Name: PermAdminNoticeDelete, Resource: "admin:notice", Operation: "delete", Description: "删除公告"},
	{Name: PermAdminNoticeList, Resource: "admin:notice", Operation: "list", Description: "查看公告列表"},

	{Name: PermAdminCheatCreate, Resource: "admin:cheat", Operation: "create", Description: "创建作弊记录"},
	{Name: PermAdminCheatUpdate, Resource: "admin:cheat", Operation: "update", Description: "更新作弊记录"},
	{Name: PermAdminCheatDelete, Resource: "admin:cheat", Operation: "delete", Description: "删除作弊记录"},
	{Name: PermAdminCheatList, Resource: "admin:cheat", Operation: "list", Description: "查看作弊记录列表"},

	{Name: PermAdminContestChallengeCreate, Resource: "admin:contest_challenge", Operation: "create", Description: "创建比赛题目关联"},
	{Name: PermAdminContestChallengeRead, Resource: "admin:contest_challenge", Operation: "read", Description: "查看比赛题目关联详情"},
	{Name: PermAdminContestChallengeUpdate, Resource: "admin:contest_challenge", Operation: "update", Description: "更新比赛题目关联"},
	{Name: PermAdminContestChallengeDelete, Resource: "admin:contest_challenge", Operation: "delete", Description: "删除比赛题目关联"},
	{Name: PermAdminContestChallengeList, Resource: "admin:contest_challenge", Operation: "list", Description: "查看比赛题目关联列表"},

	{Name: PermAdminContestChallengeFlagList, Resource: "admin:contest_challenge_flag", Operation: "list", Description: "查看比赛题目 Flag 列表"},
	{Name: PermAdminContestChallengeFlagRead, Resource: "admin:contest_challenge_flag", Operation: "read", Description: "查看比赛题目 Flag 详情"},
	{Name: PermAdminContestChallengeFlagUpdate, Resource: "admin:contest_challenge_flag", Operation: "update", Description: "更新比赛题目 Flag"},

	{Name: PermAdminImagePull, Resource: "admin:image", Operation: "pull", Description: "拉取镜像"},

	{Name: PermAdminVictimControl, Resource: "admin:victim", Operation: "control", Description: "控制靶机"},

	{Name: PermAdminFileList, Resource: "admin:file", Operation: "list", Description: "查看文件列表"},
	{Name: PermAdminFileRead, Resource: "admin:file", Operation: "read", Description: "查看文件详情"},
	{Name: PermAdminFileDelete, Resource: "admin:file", Operation: "delete", Description: "删除文件"},

	{Name: PermAdminLogRead, Resource: "admin:log", Operation: "read", Description: "查看日志"},
}

type Permission struct {
	Roles       []Role `gorm:"many2many:role_permissions" json:"-"`
	Name        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Resource    string `gorm:"type:varchar(255);index;not null" json:"resource"`
	Operation   string `gorm:"type:varchar(255);not null" json:"operation"`
	Description string `json:"description"`
	BaseModel
}

func (p Permission) TableName() string {
	return "permissions"
}

func (p Permission) ModelName() string {
	return "Permission"
}

func (p Permission) GetBaseModel() BaseModel {
	return p.BaseModel
}

func (p Permission) UniqueFields() []string {
	return []string{"id", "name"}
}

func (p Permission) QueryFields() []string {
	return []string{"id", "name", "resource", "operation", "description"}
}
