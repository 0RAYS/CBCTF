package i18n

import "fmt"

var resp = map[string]map[string]interface{}{
	"Success":         {"zh-CN": "操作成功", "en-US": "Success", "code": 200},
	"ConfigNotChange": {"zh-CN": "配置未改变", "en-US": "Configuration unchanged", "code": 200},
	"BadRequest":      {"zh-CN": "参数错误", "en-US": "Bad request", "code": 400},
	"Unauthorized":    {"zh-CN": "未登录", "en-US": "Unauthorized", "code": 401},
	"Forbidden":       {"zh-CN": "禁止访问", "en-US": "Forbidden", "code": 403},
	"TooManyRequests": {"zh-CN": "请求过于频繁", "en-US": "Too many requests", "code": 429},
	"TooQuick":        {"zh-CN": "操作过于频繁", "en-US": "Operation too frequent", "code": 429},
	"UnknownError":    {"zh-CN": "未知错误, 请联系管理员", "en-US": "UnknownError, please contact the administrator", "code": 500},
	"DeadLock":        {"zh-CN": "失败次数过多", "en-US": "Failed too many times", "code": 500},
	"UnsupportedKey":  {"zh-CN": "不支持的字段", "en-US": "Unsupported column", "code": 400},

	"CreateModelError":   {"zh-CN": "创建失败", "en-US": "Failed to create", "code": 500},
	"UpdateModelError":   {"zh-CN": "更新失败", "en-US": "Failed to update", "code": 500},
	"DeleteModelError":   {"zh-CN": "删除失败", "en-US": "Failed to delete", "code": 500},
	"GetModelError":      {"zh-CN": "获取失败", "en-US": "Failed to get", "code": 500},
	"ModelNotFound":      {"zh-CN": "记录不存在", "en-US": "Record not found", "code": 404},
	"Options2ModelError": {"zh-CN": "参数转换失败", "en-US": "Failed to transform options", "code": 500},
	"CountModelError":    {"zh-CN": "统计失败", "en-US": "Failed to count", "code": 500},

	"CreateAdminError": {"zh-CN": "创建管理员失败", "en-US": "Failed to create admin", "code": 500},
	"DeleteAdminError": {"zh-CN": "删除管理员失败", "en-US": "Failed to delete admin", "code": 500},
	"UpdateAdminError": {"zh-CN": "更新管理员失败", "en-US": "Failed to update admin", "code": 500},
	"GetAdminError":    {"zh-CN": "获取管理员失败", "en-US": "Failed to get admin", "code": 500},
	"AdminNotFound":    {"zh-CN": "管理员不存在", "en-US": "Admin not found", "code": 404},

	"CreateUserError": {"zh-CN": "创建用户失败", "en-US": "Failed to create user", "code": 500},
	"DeleteUserError": {"zh-CN": "删除用户失败", "en-US": "Failed to delete user", "code": 500},
	"UpdateUserError": {"zh-CN": "更新用户失败", "en-US": "Failed to update user", "code": 500},
	"GetUserError":    {"zh-CN": "获取用户失败", "en-US": "Failed to get user", "code": 500},
	"UserNotFound":    {"zh-CN": "用户不存在", "en-US": "User not found", "code": 404},

	"CreateContestError": {"zh-CN": "创建赛事失败", "en-US": "Failed to create contest", "code": 500},
	"DeleteContestError": {"zh-CN": "删除赛事失败", "en-US": "Failed to delete contest", "code": 500},
	"UpdateContestError": {"zh-CN": "更新赛事失败", "en-US": "Failed to update contest", "code": 500},
	"GetContestError":    {"zh-CN": "获取赛事失败", "en-US": "Failed to get contest", "code": 500},
	"ContestNotFound":    {"zh-CN": "赛事不存在", "en-US": "Contest not found", "code": 404},

	"CreateTeamError": {"zh-CN": "创建队伍失败", "en-US": "Failed to create team", "code": 500},
	"DeleteTeamError": {"zh-CN": "删除队伍失败", "en-US": "Failed to delete team", "code": 500},
	"UpdateTeamError": {"zh-CN": "更新队伍失败", "en-US": "Failed to update team", "code": 500},
	"GetTeamError":    {"zh-CN": "获取队伍失败", "en-US": "Failed to get team", "code": 500},
	"TeamNotFound":    {"zh-CN": "队伍不存在", "en-US": "Team not found", "code": 404},

	"CreateChallengeError": {"zh-CN": "创建题目失败", "en-US": "Failed to create challenge", "code": 500},
	"DeleteChallengeError": {"zh-CN": "删除题目失败", "en-US": "Failed to delete challenge", "code": 500},
	"UpdateChallengeError": {"zh-CN": "更新题目失败", "en-US": "Failed to update challenge", "code": 500},
	"GetChallengeError":    {"zh-CN": "获取题目失败", "en-US": "Failed to get challenge", "code": 500},
	"ChallengeNotFound":    {"zh-CN": "题目不存在", "en-US": "Challenge not found", "code": 404},

	"CreateNoticeError": {"zh-CN": "创建公告失败", "en-US": "Failed to create notice", "code": 500},
	"DeleteNoticeError": {"zh-CN": "删除公告失败", "en-US": "Failed to delete notice", "code": 500},
	"UpdateNoticeError": {"zh-CN": "更新公告失败", "en-US": "Failed to update notice", "code": 500},
	"GetNoticeError":    {"zh-CN": "获取公告失败", "en-US": "Failed to get notice", "code": 500},
	"NoticeNotFound":    {"zh-CN": "公告不存在", "en-US": "Notice not found", "code": 404},

	"CreateUsageError": {"zh-CN": "添加题目到比赛失败", "en-US": "Failed to add challenge to contest", "code": 500},
	"UpdateUsageError": {"zh-CN": "更新题目失败", "en-US": "Failed to update challenge", "code": 500},
	"DeleteUsageError": {"zh-CN": "移除题目失败", "en-US": "Failed to delete challenge", "code": 500},
	"GetUsageError":    {"zh-CN": "获取题目失败", "en-US": "Failed to get challenge", "code": 500},
	"UsageNotFound":    {"zh-CN": "题目不存在", "en-US": "Challenge not found", "code": 404},

	"CreateFlagError": {"zh-CN": "创建flag失败", "en-US": "Failed to create flag", "code": 500},
	"DeleteFlagError": {"zh-CN": "重置flag失败", "en-US": "Failed to delete flag", "code": 500},
	"UpdateFlagError": {"zh-CN": "更新flag失败", "en-US": "Failed to update flag", "code": 500},
	"GetFlagError":    {"zh-CN": "获取flag失败", "en-US": "Failed to get flag", "code": 500},
	"FlagNotFound":    {"zh-CN": "Flag不存在", "en-US": "Flag not found", "code": 400},

	"CreateAnswerError": {"zh-CN": "创建答案失败", "en-US": "Failed to create answer", "code": 500},
	"DeleteAnswerError": {"zh-CN": "删除答案失败", "en-US": "Failed to delete answer", "code": 500},
	"UpdateAnswerError": {"zh-CN": "更新答案失败", "en-US": "Failed to update answer", "code": 500},
	"GetAnswerError":    {"zh-CN": "获取答案失败", "en-US": "Failed to get answer", "code": 500},
	"AnswerNotFound":    {"zh-CN": "题目未初始化", "en-US": "Challenge not initialize", "code": 404},

	"CreateContainerError": {"zh-CN": "创建容器失败", "en-US": "Failed to create container", "code": 500},
	"DeleteContainerError": {"zh-CN": "删除容器失败", "en-US": "Failed to delete container", "code": 500},
	"GetContainerError":    {"zh-CN": "获取容器失败", "en-US": "Failed to get container", "code": 500},
	"ContainerNotFound":    {"zh-CN": "容器不存在", "en-US": "Container not found", "code": 404},

	"CreateSubmissionError": {"zh-CN": "创建提交记录失败", "en-US": "Failed to create submission", "code": 500},
	"DeleteSubmissionError": {"zh-CN": "删除提交记录失败", "en-US": "Failed to delete submission", "code": 500},
	"GetSubmissionError":    {"zh-CN": "获取提交记录失败", "en-US": "Failed to get submission", "code": 500},
	"SubmissionNotFound":    {"zh-CN": "提交记录不存在", "en-US": "Submission not found", "code": 404},

	"CreateTrafficError": {"zh-CN": "读取流量失败", "en-US": "Failed to create traffic", "code": 500},
	"DeleteTrafficError": {"zh-CN": "删除流量失败", "en-US": "Failed to delete traffic", "code": 500},
	"GetTrafficError":    {"zh-CN": "获取流量失败", "en-US": "Failed to get traffic", "code": 500},
	"TrafficNotFound":    {"zh-CN": "流量不存在", "en-US": "Traffic not found", "code": 404},

	"CreateRequestError": {"zh-CN": "记录请求失败", "en-US": "Failed to record request", "code": 500},
	"DeleteRequestError": {"zh-CN": "删除请求失败", "en-US": "Failed to delete request", "code": 500},
	"GetRequestError":    {"zh-CN": "获取请求失败", "en-US": "Failed to get request", "code": 500},
	"RequestNotFound":    {"zh-CN": "请求不存在", "en-US": "Request not found", "code": 404},

	"CreateDeviceError": {"zh-CN": "记录设备失败", "en-US": "Failed to record device", "code": 500},
	"DeleteDeviceError": {"zh-CN": "删除设备失败", "en-US": "Failed to delete device", "code": 500},
	"UpdateDeviceError": {"zh-CN": "更新设备失败", "en-US": "Failed to update device", "code": 500},
	"GetDeviceError":    {"zh-CN": "获取设备失败", "en-US": "Failed to get device", "code": 500},
	"DeviceNotFound":    {"zh-CN": "设备不存在", "en-US": "Device not found", "code": 404},

	"DuplicateUsername":   {"zh-CN": "该用户名已注册", "en-US": "Username already registered", "code": 400},
	"DuplicateEmail":      {"zh-CN": "该邮箱已注册", "en-US": "Email already registered", "code": 400},
	"WeakPassword":        {"zh-CN": "密码过于简单", "en-US": "Password too simple", "code": 400},
	"NameOrPasswordError": {"zh-CN": "用户名或密码错误, 请重试", "en-US": "Username or password error, please try again", "code": 401},
	"TeamIsFull":          {"zh-CN": "队伍人数已满", "en-US": "Team is full", "code": 400},
	"TeamIsBanned":        {"zh-CN": "队伍已被禁止", "en-US": "Team is banned", "code": 400},
	"UnverifiedEmail":     {"zh-CN": "邮箱未验证", "en-US": "Email not verified", "code": 403},
	"CaptchaError":        {"zh-CN": "邀请码错误", "en-US": "Captcha error", "code": 400},
	"IsNotPlayer":         {"zh-CN": "该用户未加入任何队伍", "en-US": "User has not joined any team", "code": 400},
	"InvalidEmail":        {"zh-CN": "邮箱地址无效", "en-US": "Invalid Email address", "code": 400},
	"PasswordError":       {"zh-CN": "密码错误", "en-US": "Incorrect password", "code": 401},
	"PasswordSame":        {"zh-CN": "新密码与旧密码相同", "en-US": "New password is the same as the old one", "code": 400},
	"RepeatPlayer":        {"zh-CN": "该用户已在其他队伍中", "en-US": "User is already in another team", "code": 400},

	"UserNotInTeam": {"zh-CN": "用户不在队伍中", "en-US": "User is not in the team", "code": 400},

	"JoinTeamError":     {"zh-CN": "加入队伍失败", "en-US": "Failed to join team", "code": 500},
	"LeaveTeamError":    {"zh-CN": "退出队伍失败", "en-US": "Failed to leave team", "code": 500},
	"DuplicateTeamName": {"zh-CN": "队伍名已被占用", "en-US": "Team name already taken", "code": 400},
	"DuplicateMember":   {"zh-CN": "用户已加入此队伍或其他队伍", "en-US": "User has joined this team or other team", "code": 400},

	"InvalidChallengeType": {"zh-CN": "无效的题目类型", "en-US": "Invalid challenge type", "code": 400},
	"CreateDirError":       {"zh-CN": "创建目录失败", "en-US": "Failed to create directory", "code": 500},
	"ReadDirError":         {"zh-CN": "读取目录失败", "en-US": "Failed to read directory", "code": 500},
	"InvalidFileName":      {"zh-CN": "无效的文件名, 必须符合当前题目类型", "en-US": "Invalid file name, must be matched with challenge type", "code": 400},
	"EmptyGeneratorImage":  {"zh-CN": "未找到生成器镜像", "en-US": "Generator image not found", "code": 400},
	"EmptyContainerImage":  {"zh-CN": "未找到容器镜像", "en-US": "Container image not found", "code": 400},
	"UpdateContainerError": {"zh-CN": "更新容器失败", "en-US": "Failed to update Container", "code": 500},
	"StartContainerError":  {"zh-CN": "启动容器失败", "en-US": "Failed to start Container", "code": 500},
	"StopContainerError":   {"zh-CN": "停止容器失败", "en-US": "Failed to stop Container", "code": 500},
	"HasMuchTime":          {"zh-CN": "距容器关闭20分钟内才可延长时间", "en-US": "Can only extend time within 20 minutes before the container closes", "code": 400},

	"DuplicateContestName": {"zh-CN": "赛事名已被占用", "en-US": "Contest name already taken", "code": 400},
	"ContestNotRunning":    {"zh-CN": "赛事未开始", "en-US": "Contest not running", "code": 400},
	"ContestIsOver":        {"zh-CN": "赛事已结束", "en-US": "Contest is over", "code": 400},
	"ContestIsRunning":     {"zh-CN": "赛事进行中", "en-US": "Contest is running", "code": 400},

	"CreateFileRecordError": {"zh-CN": "保存文件失败", "en-US": "Failed to save file", "code": 500},
	"DeleteFileError":       {"zh-CN": "删除文件失败", "en-US": "Failed to delete file", "code": 500},
	"FileNotAllowed":        {"zh-CN": "不支持的文件类型", "en-US": "Unsupported file type", "code": 400},
	"FileNotFound":          {"zh-CN": "文件不存在", "en-US": "File not found", "code": 404},
	"UploadFileError":       {"zh-CN": "文件上传失败", "en-US": "File upload failed", "code": 500},

	"FlagNotMatch":   {"zh-CN": "flag错误", "en-US": "Flag not match", "code": 400},
	"AlreadySolved":  {"zh-CN": "已解决该题目", "en-US": "Challenge has already been solved", "code": 200},
	"NotAllowSubmit": {"zh-CN": "不允许提交flag", "en-US": "Not allowed to submit flag", "code": 400},

	"CopyFileError":            {"zh-CN": "复制文件失败", "en-US": "Failed to copy file", "code": 500},
	"CreatePodError":           {"zh-CN": "创建Pod失败", "en-US": "Failed to create Pod", "code": 500},
	"CreateServiceError":       {"zh-CN": "创建Service失败", "en-US": "Failed to create Service", "code": 500},
	"GetPodError":              {"zh-CN": "获取Pod失败", "en-US": "Failed to get Pod", "code": 500},
	"ServiceNotFound":          {"zh-CN": "Service不存在", "en-US": "Service not found", "code": 404},
	"GetServiceError":          {"zh-CN": "获取Service失败", "en-US": "Failed to get Service", "code": 500},
	"DeletePodError":           {"zh-CN": "删除Pod失败", "en-US": "Failed to delete Pod", "code": 500},
	"DeleteServiceError":       {"zh-CN": "删除Service失败", "en-US": "Failed to delete Service", "code": 500},
	"GetNodeError":             {"zh-CN": "获取Node失败", "en-US": "Failed to get Node", "code": 500},
	"PodNotFound":              {"zh-CN": "Pod不存在", "en-US": "Pod not found", "code": 404},
	"UpdatePodError":           {"zh-CN": "更新Pod失败", "en-US": "Failed to update Pod", "code": 500},
	"NetworkPolicyNotFound":    {"zh-CN": "NetworkPolicy不存在", "en-US": "NetworkPolicy not found", "code": 404},
	"GetNetworkPolicyError":    {"zh-CN": "获取NetworkPolicy失败", "en-US": "Failed to get NetworkPolicy", "code": 500},
	"CreateNetworkPolicyError": {"zh-CN": "创建NetworkPolicy失败", "en-US": "Failed to create NetworkPolicy", "code": 500},
	"DeleteNetworkPolicyError": {"zh-CN": "删除NetworkPolicy失败", "en-US": "Failed to delete NetworkPolicy", "code": 500},
	"ExecCommandError":         {"zh-CN": "执行POD命令失败", "en-US": "Failed to execute command", "code": 500},

	"AppendUserToTeamError":      {"zh-CN": "添加用户到队伍失败", "en-US": "Failed to add user to team", "code": 500},
	"AppendTeamToContestError":   {"zh-CN": "添加队伍到赛事失败", "en-US": "Failed to add team to contest", "code": 500},
	"AppendUserToContestError":   {"zh-CN": "添加用户到赛事失败", "en-US": "Failed to add user to contest", "code": 500},
	"DeleteUserFromTeamError":    {"zh-CN": "删除用户从队伍失败", "en-US": "Failed to delete user from team", "code": 500},
	"DeleteUserFromContestError": {"zh-CN": "删除用户从赛事失败", "en-US": "Failed to delete user from contest", "code": 500},
	"DeleteAssociatedDataError":  {"zh-CN": "删除关联数据失败", "en-US": "Failed to delete associated data", "code": 500},

	"SendEmailError":           {"zh-CN": "发送邮件失败", "en-US": "Failed to send email", "code": 500},
	"SetEmailVerifyTokenError": {"zh-CN": "缓存token失败", "en-US": "Failed to set token", "code": 500},
	"GetEmailVerifyTokenError": {"zh-CN": "获取token失败", "en-US": "Failed to get token", "code": 500},
	"DelEmailVerifyTokenError": {"zh-CN": "删除token失败", "en-US": "Failed to delete token", "code": 500},
	"InvalidEmailVerifyToken":  {"zh-CN": "无效的token", "en-US": "Invalid token", "code": 400},

	"PcapNotFound":     {"zh-CN": "pcap文件不存在, 请先停止容器", "en-US": "Pcap file not found, stop container first", "code": 404},
	"HasNoTraffic":     {"zh-CN": "容器未保存流量", "en-US": "No traffic saved by the container", "code": 404},
	"ReadPcapError":    {"zh-CN": "加载pcap文件失败", "en-US": "Failed to load pcap file", "code": 500},
	"SaveTrafficError": {"zh-CN": "加载流量失败", "en-US": "Failed to load traffic", "code": 500},
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
