package i18n

import "fmt"

var resp = map[string]map[string]any{
	Success:        {"zh-CN": "操作成功", "en-US": "Success", "code": 200},
	UnsupportedKey: {"zh-CN": "不支持的键", "en-US": "Unsupported key", "code": 400},
	DeadLock:       {"zh-CN": "数据写入失败过多", "en-US": "Database deadlock", "code": 500},

	BadRequest:   {"zh-CN": "请求错误", "en-US": "Bad request", "code": 400},
	Unauthorized: {"zh-CN": "未登录", "en-US": "Unauthorized", "code": 401},
	Forbidden:    {"zh-CN": "禁止访问", "en-US": "Forbidden", "code": 403},
	UnknownError: {"zh-CN": "未知错误", "en-US": "Unknown error", "code": 500},

	CreateAdminError: {"zh-CN": "创建管理员失败", "en-US": "Create admin failed", "code": 500},
	DeleteAdminError: {"zh-CN": "删除管理员失败", "en-US": "Delete admin failed", "code": 500},
	GetAdminError:    {"zh-CN": "获取管理员失败", "en-US": "Get admin failed", "code": 500},
	AdminNotFound:    {"zh-CN": "管理员不存在", "en-US": "Admin not found", "code": 404},
	UpdateAdminError: {"zh-CN": "更新管理员失败", "en-US": "Update admin failed", "code": 500},

	CreateChallengeError: {"zh-CN": "创建题目失败", "en-US": "Create challenge failed", "code": 500},
	DeleteChallengeError: {"zh-CN": "删除题目失败", "en-US": "Delete challenge failed", "code": 500},
	GetChallengeError:    {"zh-CN": "获取题目失败", "en-US": "Get challenge failed", "code": 500},
	ChallengeNotFound:    {"zh-CN": "题目不存在", "en-US": "Challenge not found", "code": 404},
	UpdateChallengeError: {"zh-CN": "更新题目失败", "en-US": "Update challenge failed", "code": 500},

	CreateChallengeFlagError: {"zh-CN": "创建题目flag失败", "en-US": "Create challenge flag failed", "code": 500},
	DeleteChallengeFlagError: {"zh-CN": "删除题目flag失败", "en-US": "Delete challenge flag failed", "code": 500},
	GetChallengeFlagError:    {"zh-CN": "获取题目flag失败", "en-US": "Get challenge flag failed", "code": 500},
	ChallengeFlagNotFound:    {"zh-CN": "题目flag不存在", "en-US": "Challenge flag not found", "code": 404},
	UpdateChallengeFlagError: {"zh-CN": "更新题目flag失败", "en-US": "Update challenge flag failed", "code": 500},

	CreateContestError:   {"zh-CN": "创建比赛失败", "en-US": "Create contest failed", "code": 500},
	DeleteContestError:   {"zh-CN": "删除比赛失败", "en-US": "Delete contest failed", "code": 500},
	GetContestError:      {"zh-CN": "获取比赛失败", "en-US": "Get contest failed", "code": 500},
	ContestNotFound:      {"zh-CN": "比赛不存在", "en-US": "Contest not found", "code": 404},
	UpdateContestError:   {"zh-CN": "更新比赛失败", "en-US": "Update contest failed", "code": 500},
	DuplicateContestName: {"zh-CN": "比赛名称已被使用", "en-US": "Contest name already in use", "code": 400},
	ContestCaptchaError:  {"zh-CN": "比赛邀请码错误", "en-US": "Contest captcha error", "code": 400},
	ContestIsComing:      {"zh-CN": "比赛未开始", "en-US": "Contest is coming", "code": 400},
	ContestIsRunning:     {"zh-CN": "比赛正在进行中", "en-US": "Contest is running", "code": 400},
	ContestIsOver:        {"zh-CN": "比赛已结束", "en-US": "Contest is over", "code": 400},

	CreateContestChallengeError: {"zh-CN": "添加比赛题目失败", "en-US": "Add contest challenge failed", "code": 500},
	DeleteContestChallengeError: {"zh-CN": "删除比赛题目失败", "en-US": "Delete contest challenge failed", "code": 500},
	GetContestChallengeError:    {"zh-CN": "获取比赛题目失败", "en-US": "Get contest challenge failed", "code": 500},
	ContestChallengeNotFound:    {"zh-CN": "比赛题目不存在", "en-US": "Contest challenge not found", "code": 404},
	UpdateContestChallengeError: {"zh-CN": "更新比赛题目失败", "en-US": "Update contest challenge failed", "code": 500},

	CreateContestFlagError: {"zh-CN": "创建比赛flag失败", "en-US": "Create contest flag failed", "code": 500},
	DeleteContestFlagError: {"zh-CN": "删除比赛flag失败", "en-US": "Delete contest flag failed", "code": 500},
	GetContestFlagError:    {"zh-CN": "获取比赛flag失败", "en-US": "Get contest flag failed", "code": 500},
	ContestFlagNotFound:    {"zh-CN": "比赛flag不存在", "en-US": "Contest flag not found", "code": 404},
	UpdateContestFlagError: {"zh-CN": "更新比赛flag失败", "en-US": "Update contest flag failed", "code": 500},

	CreateDeviceError: {"zh-CN": "创建设备失败", "en-US": "Create device failed", "code": 500},
	DeleteDeviceError: {"zh-CN": "删除设备失败", "en-US": "Delete device failed", "code": 500},
	GetDeviceError:    {"zh-CN": "获取设备失败", "en-US": "Get device failed", "code": 500},
	DeviceNotFound:    {"zh-CN": "设备不存在", "en-US": "Device not found", "code": 404},
	UpdateDeviceError: {"zh-CN": "更新设备失败", "en-US": "Update device failed", "code": 500},

	CreateDockerError: {"zh-CN": "创建Docker失败", "en-US": "Create Docker failed", "code": 500},
	DeleteDockerError: {"zh-CN": "删除Docker失败", "en-US": "Delete Docker failed", "code": 500},
	GetDockerError:    {"zh-CN": "获取Docker失败", "en-US": "Get Docker failed", "code": 500},
	DockerNotFound:    {"zh-CN": "Docker不存在", "en-US": "Docker not found", "code": 404},
	UpdateDockerError: {"zh-CN": "更新Docker失败", "en-US": "Update Docker failed", "code": 500},

	CreateDockerGroupError: {"zh-CN": "创建Docker组失败", "en-US": "Create Docker group failed", "code": 500},
	DeleteDockerGroupError: {"zh-CN": "删除Docker组失败", "en-US": "Delete Docker group failed", "code": 500},
	GetDockerGroupError:    {"zh-CN": "获取Docker组失败", "en-US": "Get Docker group failed", "code": 500},
	DockerGroupNotFound:    {"zh-CN": "Docker组不存在", "en-US": "Docker group not found", "code": 404},
	UpdateDockerGroupError: {"zh-CN": "更新Docker组失败", "en-US": "Update Docker group failed", "code": 500},

	CreateEventError: {"zh-CN": "创建事件失败", "en-US": "Create event failed", "code": 500},
	DeleteEventError: {"zh-CN": "删除事件失败", "en-US": "Delete event failed", "code": 500},
	GetEventError:    {"zh-CN": "获取事件失败", "en-US": "Get event failed", "code": 500},
	EventNotFound:    {"zh-CN": "事件不存在", "en-US": "Event not found", "code": 404},
	UpdateEventError: {"zh-CN": "更新事件失败", "en-US": "Update event failed", "code": 500},

	CreateFileError: {"zh-CN": "创建文件失败", "en-US": "Create file failed", "code": 500},
	DeleteFileError: {"zh-CN": "删除文件失败", "en-US": "Delete file failed", "code": 500},
	GetFileError:    {"zh-CN": "获取文件失败", "en-US": "Get file failed", "code": 500},
	FileNotFound:    {"zh-CN": "文件不存在", "en-US": "File not found", "code": 404},
	UpdateFileError: {"zh-CN": "更新文件失败", "en-US": "Update file failed", "code": 500},
	FileNotAllowed:  {"zh-CN": "不允许的文件类型", "en-US": "File type not allowed", "code": 400},

	CreateNoticeError: {"zh-CN": "创建通知失败", "en-US": "Create notice failed", "code": 500},
	DeleteNoticeError: {"zh-CN": "删除通知失败", "en-US": "Delete notice failed", "code": 500},
	GetNoticeError:    {"zh-CN": "获取通知失败", "en-US": "Get notice failed", "code": 500},
	NoticeNotFound:    {"zh-CN": "通知不存在", "en-US": "Notice not found", "code": 404},
	UpdateNoticeError: {"zh-CN": "更新通知失败", "en-US": "Update notice failed", "code": 500},

	CreateRequestError: {"zh-CN": "创建请求失败", "en-US": "Create request failed", "code": 500},
	DeleteRequestError: {"zh-CN": "删除请求失败", "en-US": "Delete request failed", "code": 500},
	GetRequestError:    {"zh-CN": "获取请求失败", "en-US": "Get request failed", "code": 500},
	RequestNotFound:    {"zh-CN": "请求不存在", "en-US": "Request not found", "code": 404},
	UpdateRequestError: {"zh-CN": "更新请求失败", "en-US": "Update request failed", "code": 500},

	CreateTeamError:    {"zh-CN": "创建战队失败", "en-US": "Create team failed", "code": 500},
	DeleteTeamError:    {"zh-CN": "删除战队失败", "en-US": "Delete team failed", "code": 500},
	GetTeamError:       {"zh-CN": "获取战队失败", "en-US": "Get team failed", "code": 500},
	TeamNotFound:       {"zh-CN": "战队不存在", "en-US": "Team not found", "code": 404},
	UpdateTeamError:    {"zh-CN": "更新战队失败", "en-US": "Update team failed", "code": 500},
	TeamIsBanned:       {"zh-CN": "战队已被封禁", "en-US": "Team is banned", "code": 403},
	TeamIsFull:         {"zh-CN": "战队已满员", "en-US": "Team is full", "code": 400},
	DuplicateTeamName:  {"zh-CN": "战队名称已被使用", "en-US": "Team name already in use", "code": 400},
	DuplicateMember:    {"zh-CN": "用户已在战队中", "en-US": "User already in team", "code": 400},
	UserNotInTeam:      {"zh-CN": "用户不在战队中", "en-US": "User not in team", "code": 400},
	CaptainCannotLeave: {"zh-CN": "队长不能离开战队", "en-US": "Captain cannot leave team", "code": 400},
	TeamCaptchaError:   {"zh-CN": "战队邀请码错误", "en-US": "Team captcha error", "code": 400},

	CreateUserError:     {"zh-CN": "创建用户失败", "en-US": "Create user failed", "code": 500},
	DeleteUserError:     {"zh-CN": "删除用户失败", "en-US": "Delete user failed", "code": 500},
	GetUserError:        {"zh-CN": "获取用户失败", "en-US": "Get user failed", "code": 500},
	UserNotFound:        {"zh-CN": "用户不存在", "en-US": "User not found", "code": 404},
	UpdateUserError:     {"zh-CN": "更新用户失败", "en-US": "Update user failed", "code": 500},
	InvalidEmail:        {"zh-CN": "无效的邮箱", "en-US": "Invalid email", "code": 400},
	UnverifiedEmail:     {"zh-CN": "邮箱未验证", "en-US": "Email not verified", "code": 400},
	DuplicateEmail:      {"zh-CN": "邮箱已被使用", "en-US": "Email already in use", "code": 400},
	DuplicateUserName:   {"zh-CN": "用户名已被使用", "en-US": "Username already in use", "code": 400},
	WeakPassword:        {"zh-CN": "密码过于简单", "en-US": "Weak password", "code": 400},
	NameOrPasswordError: {"zh-CN": "用户名或密码错误", "en-US": "Username or password error", "code": 401},
	PasswordError:       {"zh-CN": "密码错误", "en-US": "Password error", "code": 401},
	PasswordSame:        {"zh-CN": "新密码与旧密码相同", "en-US": "New password is the same as old password", "code": 400},

	AppendUserToTeamError:      {"zh-CN": "添加用户到战队失败", "en-US": "Append user to team failed", "code": 500},
	AppendUserToContestError:   {"zh-CN": "添加用户到比赛失败", "en-US": "Append user to contest failed", "code": 500},
	DeleteUserFromTeamError:    {"zh-CN": "从战队中删除用户失败", "en-US": "Delete user from team failed", "code": 500},
	DeleteUserFromContestError: {"zh-CN": "从比赛中删除用户失败", "en-US": "Delete user from contest failed", "code": 500},

	SetEmailVerifyTokenError: {"zh-CN": "设置邮箱验证令牌失败", "en-US": "Set email verify token failed", "code": 500},
	GetEmailVerifyTokenError: {"zh-CN": "获取邮箱验证令牌失败", "en-US": "Get email verify token failed", "code": 500},
	DelEmailVerifyTokenError: {"zh-CN": "删除邮箱验证令牌失败", "en-US": "Delete email verify token failed", "code": 500},
	SendEmailError:           {"zh-CN": "发送邮件失败", "en-US": "Send email failed", "code": 500},
}

// I18N 获取翻译与状态码, 非http响应状态码
func I18N(key string, language string) (string, int) {
	if v, ok := resp[key]; !ok {
		switch language {
		case "en-US":
			return fmt.Sprintf("I18N configuration is incomplete: %s", key), 400
		default:
			return fmt.Sprintf("I18N 配置不完全: %s", key), 400
		}
	} else {
		if language == "origin" {
			return key, v["code"].(int)
		}
		return v[language].(string), v["code"].(int)
	}
}
