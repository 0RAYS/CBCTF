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

	PermAdminSystemStatus  = "admin:system:status"
	PermAdminSystemRead    = "admin:system:read"
	PermAdminSystemUpdate  = "admin:system:update"
	PermAdminSystemRestart = "admin:system:restart"

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

	PermAdminCronJobUpdate = "admin:cronjob:update"
	PermAdminCronJobList   = "admin:cronjob:list"

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

	PermAdminImagePull        = "admin:image:pull"
	PermAdminContestImagePull = "admin:contest_image:pull"

	PermAdminVictimControl        = "admin:victim:control"
	PermAdminContestVictimControl = "admin:contest_victim:control"

	PermAdminGeneratorControl        = "admin:generator:control"
	PermAdminContestGeneratorControl = "admin:contest_generator:control"

	PermAdminTrafficRead        = "admin:traffic:read"
	PermAdminContestTrafficRead = "admin:contest_traffic:read"

	PermAdminFileList   = "admin:file:list"
	PermAdminFileRead   = "admin:file:read"
	PermAdminFileDelete = "admin:file:delete"

	PermAdminTaskRead = "admin:task:read"

	PermAdminLogRead = "admin:log:read"
)

// RoutePermissions maps "METHOD /full-path" to the required permission name.
// Used by CheckPermission middleware and the /me/permissions handler.
var RoutePermissions = map[string]string{
	// /me
	"GET /me":             PermSelfRead,
	"GET /me/permissions": PermSelfRead,
	"PUT /me":             PermSelfUpdate,
	"PUT /me/password":    PermSelfUpdate,
	"DELETE /me":          PermSelfDelete,
	"POST /me/picture":    PermSelfUpdate,
	"POST /me/activate":   PermSelfActivate,

	// /contests/:contestID
	"GET /contests/:contestID":            PermUserContestRead,
	"GET /contests/:contestID/rank":       PermUserContestRank,
	"GET /contests/:contestID/scoreboard": PermUserContestRank,
	"GET /contests/:contestID/timeline":   PermUserContestRank,

	"POST /contests/:contestID/teams/join":   PermUserTeamJoin,
	"POST /contests/:contestID/teams/create": PermUserTeamCreate,

	"GET /contests/:contestID/teams/me":          PermUserTeamRead,
	"GET /contests/:contestID/teams/me/captcha":  PermUserTeamRead,
	"GET /contests/:contestID/teams/me/users":    PermUserTeamRead,
	"PUT /contests/:contestID/teams/me/captcha":  PermUserTeamUpdate,
	"PUT /contests/:contestID/teams/me":          PermUserTeamUpdate,
	"POST /contests/:contestID/teams/me/picture": PermUserTeamUpdate,
	"DELETE /contests/:contestID/teams/me":       PermUserTeamDelete,
	"POST /contests/:contestID/teams/me/kick":    PermUserTeamUpdate,
	"POST /contests/:contestID/teams/me/leave":   PermUserTeamRead,

	"GET /contests/:contestID/notices": PermUserNoticeList,

	"GET /contests/:contestID/challenges":                         PermUserChallengeList,
	"GET /contests/:contestID/challenges/categories":              PermUserChallengeList,
	"GET /contests/:contestID/challenges/:challengeID":            PermUserChallengeRead,
	"GET /contests/:contestID/challenges/:challengeID/attachment": PermUserChallengeRead,
	"POST /contests/:contestID/challenges/:challengeID/init":      PermUserChallengeInit,
	"POST /contests/:contestID/challenges/:challengeID/reset":     PermUserChallengeReset,
	"POST /contests/:contestID/challenges/:challengeID/start":     PermUserVictimControl,
	"POST /contests/:contestID/challenges/:challengeID/increase":  PermUserVictimControl,
	"POST /contests/:contestID/challenges/:challengeID/stop":      PermUserVictimControl,
	"POST /contests/:contestID/challenges/:challengeID/submit":    PermUserChallengeSubmit,

	"POST /contests/:contestID/writeups": PermUserWriteupUpload,
	"GET /contests/:contestID/writeups":  PermUserWriteupList,

	// /admin 基础
	"GET /admin/ip":     PermAdminIPSearch,
	"GET /admin/models": PermAdminModelsSearch,
	"GET /admin/search": PermAdminModelsSearch,

	// /admin/system
	"GET /admin/system/status":   PermAdminSystemStatus,
	"GET /admin/system/config":   PermAdminSystemRead,
	"PUT /admin/system/config":   PermAdminSystemUpdate,
	"POST /admin/system/restart": PermAdminSystemRestart,

	// /admin/permissions
	"GET /admin/permissions":               PermAdminPermissionList,
	"PUT /admin/permissions/:permissionID": PermAdminPermissionUpdate,

	// /admin/roles
	"GET /admin/roles":                        PermAdminRoleList,
	"POST /admin/roles":                       PermAdminRoleCreate,
	"GET /admin/roles/:roleID":                PermAdminRoleRead,
	"GET /admin/roles/:roleID/permissions":    PermAdminPermissionList,
	"PUT /admin/roles/:roleID":                PermAdminRoleUpdate,
	"DELETE /admin/roles/:roleID":             PermAdminRoleDelete,
	"POST /admin/roles/:roleID/permissions":   PermAdminRoleAssign,
	"DELETE /admin/roles/:roleID/permissions": PermAdminRoleRevoke,

	// /admin/groups
	"GET /admin/groups":                   PermAdminGroupList,
	"POST /admin/groups":                  PermAdminGroupCreate,
	"GET /admin/groups/:groupID":          PermAdminGroupRead,
	"GET /admin/groups/:groupID/users":    PermAdminUserList,
	"PUT /admin/groups/:groupID":          PermAdminGroupUpdate,
	"DELETE /admin/groups/:groupID":       PermAdminGroupDelete,
	"POST /admin/groups/:groupID/users":   PermAdminUserAssign,
	"DELETE /admin/groups/:groupID/users": PermAdminUserRevoke,

	// /admin/users
	"GET /admin/users":                  PermAdminUserList,
	"POST /admin/users":                 PermAdminUserCreate,
	"GET /admin/users/:userID":          PermAdminUserRead,
	"PUT /admin/users/:userID":          PermAdminUserUpdate,
	"DELETE /admin/users/:userID":       PermAdminUserDelete,
	"POST /admin/users/:userID/picture": PermAdminUserUpdate,

	// /admin/oauth
	"GET /admin/oauth":                   PermAdminOauthList,
	"POST /admin/oauth":                  PermAdminOauthCreate,
	"PUT /admin/oauth/:oauthID":          PermAdminOauthUpdate,
	"POST /admin/oauth/:oauthID/picture": PermAdminOauthUpdate,
	"DELETE /admin/oauth/:oauthID":       PermAdminOauthDelete,

	// /admin/email + /admin/smtp
	"GET /admin/email":              PermAdminSMTPList,
	"GET /admin/smtp":               PermAdminSMTPList,
	"POST /admin/smtp":              PermAdminSMTPCreate,
	"PUT /admin/smtp/:smtpID":       PermAdminSMTPUpdate,
	"DELETE /admin/smtp/:smtpID":    PermAdminSMTPDelete,
	"GET /admin/smtp/:smtpID/email": PermAdminSMTPList,

	// /admin/cronjobs
	"GET /admin/cronjobs":            PermAdminCronJobList,
	"GET /admin/cronjobs/:cronJobID": PermAdminCronJobList,
	"PUT /admin/cronjobs/:cronJobID": PermAdminCronJobUpdate,

	// /admin/webhook
	"GET /admin/webhook":                    PermAdminWebhookList,
	"GET /admin/webhook/events":             PermAdminWebhookList,
	"GET /admin/webhook/history":            PermAdminWebhookList,
	"POST /admin/webhook":                   PermAdminWebhookCreate,
	"PUT /admin/webhook/:webhookID":         PermAdminWebhookUpdate,
	"DELETE /admin/webhook/:webhookID":      PermAdminWebhookDelete,
	"GET /admin/webhook/:webhookID/history": PermAdminWebhookList,

	// /admin/challenges
	"GET /admin/challenges":                              PermAdminChallengeList,
	"GET /admin/challenges/categories":                   PermAdminChallengeList,
	"POST /admin/challenges":                             PermAdminChallengeCreate,
	"GET /admin/challenges/:challengeID/download":        PermAdminChallengeRead,
	"PUT /admin/challenges/:challengeID":                 PermAdminChallengeUpdate,
	"DELETE /admin/challenges/:challengeID":              PermAdminChallengeDelete,
	"POST /admin/challenges/:challengeID/upload":         PermAdminChallengeUpdate,
	"GET /admin/challenges/:challengeID/test":            PermAdminChallengeTest,
	"GET /admin/challenges/:challengeID/test/attachment": PermAdminChallengeTest,
	"POST /admin/challenges/:challengeID/test/start":     PermAdminChallengeTest,
	"POST /admin/challenges/:challengeID/test/stop":      PermAdminChallengeTest,

	"GET /admin/victims":                            PermAdminVictimControl,
	"DELETE /admin/victims":                         PermAdminVictimControl,
	"GET /admin/victims/:victimID/traffic":          PermAdminTrafficRead,
	"GET /admin/victims/:victimID/traffic/download": PermAdminTrafficRead,

	"GET /admin/generators":    PermAdminGeneratorControl,
	"POST /admin/generators":   PermAdminGeneratorControl,
	"DELETE /admin/generators": PermAdminGeneratorControl,

	// /admin/images
	"GET /admin/images":  PermAdminImagePull,
	"POST /admin/images": PermAdminImagePull,

	// /admin/contests
	"GET /admin/contests":                       PermAdminContestList,
	"POST /admin/contests":                      PermAdminContestCreate,
	"GET /admin/contests/:contestID":            PermAdminContestRead,
	"PUT /admin/contests/:contestID":            PermAdminContestUpdate,
	"DELETE /admin/contests/:contestID":         PermAdminContestDelete,
	"POST /admin/contests/:contestID/picture":   PermAdminContestUpdate,
	"GET /admin/contests/:contestID/rank":       PermAdminContestRank,
	"GET /admin/contests/:contestID/scoreboard": PermAdminContestRank,
	"GET /admin/contests/:contestID/timeline":   PermAdminContestRank,

	// /admin/contests/:contestID/teams
	"GET /admin/contests/:contestID/teams":                                            PermAdminTeamList,
	"GET /admin/contests/:contestID/teams/:teamID":                                    PermAdminTeamRead,
	"GET /admin/contests/:contestID/teams/:teamID/users":                              PermAdminTeamRead,
	"PUT /admin/contests/:contestID/teams/:teamID":                                    PermAdminTeamUpdate,
	"DELETE /admin/contests/:contestID/teams/:teamID":                                 PermAdminTeamDelete,
	"POST /admin/contests/:contestID/teams/:teamID/kick":                              PermAdminTeamUpdate,
	"POST /admin/contests/:contestID/teams/:teamID/picture":                           PermAdminTeamUpdate,
	"GET /admin/contests/:contestID/teams/:teamID/flags":                              PermAdminTeamRead,
	"GET /admin/contests/:contestID/teams/:teamID/submissions":                        PermAdminTeamRead,
	"GET /admin/contests/:contestID/teams/:teamID/victims":                            PermAdminContestTrafficRead,
	"GET /admin/contests/:contestID/teams/:teamID/victims/:victimID/traffic":          PermAdminContestTrafficRead,
	"GET /admin/contests/:contestID/teams/:teamID/victims/:victimID/traffic/download": PermAdminContestTrafficRead,
	"GET /admin/contests/:contestID/teams/:teamID/writeups":                           PermAdminTeamWriteupList,
	"GET /admin/contests/:contestID/teams/:teamID/writeups/:fileID":                   PermAdminTeamWriteupRead,

	// /admin/contests/:contestID/notices
	"GET /admin/contests/:contestID/notices":              PermAdminNoticeList,
	"POST /admin/contests/:contestID/notices":             PermAdminNoticeCreate,
	"PUT /admin/contests/:contestID/notices/:noticeID":    PermAdminNoticeUpdate,
	"DELETE /admin/contests/:contestID/notices/:noticeID": PermAdminNoticeDelete,

	// /admin/contests/:contestID/cheats
	"GET /admin/contests/:contestID/cheats":             PermAdminCheatList,
	"DELETE /admin/contests/:contestID/cheats":          PermAdminCheatDelete,
	"POST /admin/contests/:contestID/cheats":            PermAdminCheatCreate,
	"PUT /admin/contests/:contestID/cheats/:cheatID":    PermAdminCheatUpdate,
	"DELETE /admin/contests/:contestID/cheats/:cheatID": PermAdminCheatDelete,

	// /admin/contests/:contestID/challenges
	"GET /admin/contests/:contestID/challenges":                                    PermAdminContestChallengeList,
	"GET /admin/contests/:contestID/challenges/others":                             PermAdminChallengeList,
	"GET /admin/contests/:contestID/challenges/categories":                         PermAdminContestChallengeList,
	"POST /admin/contests/:contestID/challenges":                                   PermAdminContestChallengeCreate,
	"PUT /admin/contests/:contestID/challenges/:challengeID":                       PermAdminContestChallengeUpdate,
	"DELETE /admin/contests/:contestID/challenges/:challengeID":                    PermAdminContestChallengeDelete,
	"GET /admin/contests/:contestID/challenges/:challengeID/flags":                 PermAdminContestChallengeFlagList,
	"PUT /admin/contests/:contestID/challenges/:challengeID/flags/:flagID":         PermAdminContestChallengeFlagUpdate,
	"GET /admin/contests/:contestID/challenges/:challengeID/flags/:flagID/solvers": PermAdminContestChallengeFlagList,

	// /admin/contests/:contestID/images
	"GET /admin/contests/:contestID/images":  PermAdminContestImagePull,
	"POST /admin/contests/:contestID/images": PermAdminContestImagePull,

	// /admin/contests/:contestID/victims
	"GET /admin/contests/:contestID/victims":    PermAdminContestVictimControl,
	"POST /admin/contests/:contestID/victims":   PermAdminContestVictimControl,
	"DELETE /admin/contests/:contestID/victims": PermAdminContestVictimControl,

	// /admin/contests/:contestID/generators
	"GET /admin/contests/:contestID/generators":    PermAdminContestGeneratorControl,
	"POST /admin/contests/:contestID/generators":   PermAdminContestGeneratorControl,
	"DELETE /admin/contests/:contestID/generators": PermAdminContestGeneratorControl,

	// /admin/files
	"GET /admin/files":         PermAdminFileList,
	"DELETE /admin/files":      PermAdminFileDelete,
	"GET /admin/files/:fileID": PermAdminFileRead,

	// /admin/tasks
	"GET /admin/tasks":      PermAdminTaskRead,
	"GET /admin/tasks/live": PermAdminTaskRead,

	// /admin/logs
	"GET /admin/logs": PermAdminLogRead,
}

var Permissions = []Permission{
	{Name: PermAdminContestImagePull, Resource: "admin:contest_image", Operation: "pull", Description: "拉取比赛镜像"},
	{Name: PermAdminContestVictimControl, Resource: "admin:contest_victim", Operation: "control", Description: "控制比赛靶机"},
	{Name: PermAdminContestGeneratorControl, Resource: "admin:contest_generator", Operation: "control", Description: "控制比赛生成器"},
	{Name: PermAdminTrafficRead, Resource: "admin:traffic", Operation: "read", Description: "查看全局靶机流量"},
	{Name: PermAdminContestTrafficRead, Resource: "admin:contest_traffic", Operation: "read", Description: "查看比赛靶机流量"},
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

	{Name: PermAdminSystemStatus, Resource: "admin:system", Operation: "status", Description: "查看系统状态"},
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

	{Name: PermAdminCronJobUpdate, Resource: "admin:cronjob", Operation: "update", Description: "更新 Cron 任务"},
	{Name: PermAdminCronJobList, Resource: "admin:cronjob", Operation: "list", Description: "查看 Cron 任务列表"},

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

	{Name: PermAdminGeneratorControl, Resource: "admin:generator", Operation: "control", Description: "控制附件生成器"},

	{Name: PermAdminFileList, Resource: "admin:file", Operation: "list", Description: "查看文件列表"},
	{Name: PermAdminFileRead, Resource: "admin:file", Operation: "read", Description: "查看文件详情"},
	{Name: PermAdminFileDelete, Resource: "admin:file", Operation: "delete", Description: "删除文件"},

	{Name: PermAdminTaskRead, Resource: "admin:task", Operation: "read", Description: "查看任务队列"},

	{Name: PermAdminLogRead, Resource: "admin:log", Operation: "read", Description: "查看日志"},
}

// Permission 权限
// ManyToMany Role
type Permission struct {
	Roles       []Role `gorm:"many2many:role_permissions" json:"-"`
	Name        string `gorm:"type:varchar(255);uniqueIndex:idx_permissions_name_active,where:deleted_at IS NULL;not null" json:"name"`
	Resource    string `gorm:"type:varchar(255);index;not null" json:"resource"`
	Operation   string `gorm:"type:varchar(255);not null" json:"operation"`
	Description string `json:"description"`
	BaseModel
}
