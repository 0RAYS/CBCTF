package i18n

var resp = map[string]map[string]interface{}{
	"Success":         {"zh-CN": "操作成功", "en-US": "Success", "code": 200},
	"ConfigNotChange": {"zh-CN": "配置未改变", "en-US": "Configuration unchanged", "code": 200},
	"BadRequest":      {"zh-CN": "参数错误", "en-US": "Bad request", "code": 400},
	"Unauthorized":    {"zh-CN": "未登录", "en-US": "Unauthorized", "code": 401},
	"Forbidden":       {"zh-CN": "禁止访问", "en-US": "Forbidden", "code": 403},
	"TooManyRequests": {"zh-CN": "请求过于频繁", "en-US": "Too many requests", "code": 429},
	"UnknownError":    {"zh-CN": "未知错误，请联系管理员", "en-US": "UnknownError, please contact the administrator", "code": 500},

	"EmailExists":         {"zh-CN": "该邮箱已注册", "en-US": "Email already registered", "code": 400},
	"NameOrPasswordError": {"zh-CN": "用户名或密码错误，请重试", "en-US": "Username or password error, please try again", "code": 401},
	"TeamIsFull":          {"zh-CN": "队伍人数已满", "en-US": "Team is full", "code": 400},
	"TeamIsBanned":        {"zh-CN": "队伍已被禁止", "en-US": "Team is banned", "code": 400},
	"UnverifiedEmail":     {"zh-CN": "邮箱未验证", "en-US": "Email not verified", "code": 403},
	"CaptchaError":        {"zh-CN": "邀请码错误", "en-US": "Captcha error", "code": 400},
	"IsNotPlayer":         {"zh-CN": "该用户未加入任何队伍", "en-US": "User has not joined any team", "code": 400},
	"InvalidEmail":        {"zh-CN": "邮箱地址无效", "en-US": "Invalid Email address", "code": 400},
	"PasswordError":       {"zh-CN": "密码错误", "en-US": "Incorrect password", "code": 401},
	"PasswordSame":        {"zh-CN": "新密码与旧密码相同", "en-US": "New password is the same as the old one", "code": 400},
	"RepeatPlayer":        {"zh-CN": "该用户已在其他队伍中", "en-US": "User is already in another team", "code": 400},

	"CreateUserError": {"zh-CN": "创建用户失败", "en-US": "Failed to create user", "code": 500},
	"DeleteUserError": {"zh-CN": "删除用户失败", "en-US": "Failed to delete user", "code": 500},
	"UpdateUserError": {"zh-CN": "更新用户失败", "en-US": "Failed to update user", "code": 500},
	"UserNameExists":  {"zh-CN": "用户名已被占用", "en-US": "Username already taken", "code": 400},
	"UserNotFound":    {"zh-CN": "用户不存在", "en-US": "User not found", "code": 404},
	"UserNotInTeam":   {"zh-CN": "用户不在队伍中", "en-US": "User is not in the team", "code": 400},

	"CreateTeamError":  {"zh-CN": "创建队伍失败", "en-US": "Failed to create team", "code": 500},
	"DeleteTeamError":  {"zh-CN": "删除队伍失败", "en-US": "Failed to delete team", "code": 500},
	"UpdateTeamError":  {"zh-CN": "更新队伍失败", "en-US": "Failed to update team", "code": 500},
	"JoinTeamError":    {"zh-CN": "加入队伍失败", "en-US": "Failed to join team", "code": 500},
	"LeaveTeamError":   {"zh-CN": "退出队伍失败", "en-US": "Failed to leave team", "code": 500},
	"TeamNotFound":     {"zh-CN": "队伍不存在", "en-US": "Team not found", "code": 404},
	"TeamNameExists":   {"zh-CN": "队伍名已被占用", "en-US": "Team name already taken", "code": 400},
	"TeamMemberExists": {"zh-CN": "用户已加入此队伍或其他队伍", "en-US": "User has joined this team or other team", "code": 400},

	"CreateAdminError": {"zh-CN": "创建管理员失败", "en-US": "Failed to create admin", "code": 500},
	"DeleteAdminError": {"zh-CN": "删除管理员失败", "en-US": "Failed to delete admin", "code": 500},
	"UpdateAdminError": {"zh-CN": "更新管理员失败", "en-US": "Failed to update admin", "code": 500},

	"CreateChallengeError": {"zh-CN": "创建题目失败", "en-US": "Failed to create challenge", "code": 500},
	"ChallengeNotFound":    {"zh-CN": "题目不存在", "en-US": "Challenge not found", "code": 404},
	"DeleteChallengeError": {"zh-CN": "删除题目失败", "en-US": "Failed to delete challenge", "code": 500},
	"UpdateChallengeError": {"zh-CN": "更新题目失败", "en-US": "Failed to update challenge", "code": 500},
	"InvalidChallengeType": {"zh-CN": "无效的题目类型", "en-US": "Invalid challenge type", "code": 400},
	"CreateDirError":       {"zh-CN": "创建目录失败", "en-US": "Failed to create directory", "code": 500},
	"ReadDirError":         {"zh-CN": "读取目录失败", "en-US": "Failed to read directory", "code": 500},
	"InvalidFileName":      {"zh-CN": "无效的文件名，必须符合当前题目类型", "en-US": "Invalid file name, must be matched with challenge type", "code": 400},
	"EmptyGeneratorImage":  {"zh-CN": "未找到生成器镜像", "en-US": "Generator image not found", "code": 400},
	"EmptyDockerImage":     {"zh-CN": "未找到Docker镜像", "en-US": "Docker image not found", "code": 400},
	"CreateDockerError":    {"zh-CN": "创建容器失败", "en-US": "Failed to create Docker", "code": 500},
	"GetDockerError":       {"zh-CN": "获取容器失败", "en-US": "Failed to get Docker", "code": 500},
	"DeleteDockerError":    {"zh-CN": "删除容器失败", "en-US": "Failed to delete Docker", "code": 500},
	"UpdateDockerError":    {"zh-CN": "更新容器失败", "en-US": "Failed to update Docker", "code": 500},
	"StartDockerError":     {"zh-CN": "启动容器失败", "en-US": "Failed to start Docker", "code": 500},
	"StopDockerError":      {"zh-CN": "停止容器失败", "en-US": "Failed to stop Docker", "code": 500},
	"DockerNotFound":       {"zh-CN": "容器不存在", "en-US": "Docker not found", "code": 404},
	"HasMuchTime":          {"zh-CN": "距容器关闭20分钟内才可延长时间", "en-US": "Can only extend time within 20 minutes before the container closes", "code": 400},

	"CreateContestError": {"zh-CN": "创建赛事失败", "en-US": "Failed to create contest", "code": 500},
	"DeleteContestError": {"zh-CN": "删除赛事失败", "en-US": "Failed to delete contest", "code": 500},
	"UpdateContestError": {"zh-CN": "更新赛事失败", "en-US": "Failed to update contest", "code": 500},
	"ContestNameExists":  {"zh-CN": "赛事名已被占用", "en-US": "Contest name already taken", "code": 400},
	"ContestNotFound":    {"zh-CN": "赛事不存在", "en-US": "Contest not found", "code": 404},
	"ContestNotRunning":  {"zh-CN": "赛事未开始", "en-US": "Contest not running", "code": 400},
	"ContestIsOver":      {"zh-CN": "赛事已结束", "en-US": "Contest is over", "code": 400},
	"ContestIsRunning":   {"zh-CN": "赛事进行中", "en-US": "Contest is running", "code": 400},

	"CreateUsageError": {"zh-CN": "添加题目到比赛失败", "en-US": "Failed to add challenge to contest", "code": 500},
	"GetUsageError":    {"zh-CN": "获取题目失败", "en-US": "Failed to get challenge", "code": 500},
	"UsageNotFound":    {"zh-CN": "题目不存在", "en-US": "Challenge not found", "code": 404},
	"UpdateUsageError": {"zh-CN": "更新题目失败", "en-US": "Failed to update challenge", "code": 500},
	"DeleteUsageError": {"zh-CN": "移除题目失败", "en-US": "Failed to delete challenge", "code": 500},

	"CreateFileRecordError": {"zh-CN": "保存文件失败", "en-US": "Failed to save file", "code": 500},
	"DeleteFileError":       {"zh-CN": "删除文件失败", "en-US": "Failed to delete file", "code": 500},
	"FileNotAllowed":        {"zh-CN": "不支持的文件类型", "en-US": "Unsupported file type", "code": 400},
	"FileNotFound":          {"zh-CN": "文件不存在", "en-US": "File not found", "code": 404},
	"UploadFileError":       {"zh-CN": "文件上传失败", "en-US": "File upload failed", "code": 500},

	"CreateFlagError":       {"zh-CN": "创建flag失败", "en-US": "Failed to create flag", "code": 500},
	"FlagNotFound":          {"zh-CN": "题目未初始化", "en-US": "Challenge not initialized", "code": 404},
	"UpdateFlagError":       {"zh-CN": "更新flag失败", "en-US": "Failed to update flag", "code": 500},
	"FlagNotMatch":          {"zh-CN": "flag错误", "en-US": "Flag not match", "code": 400},
	"SubmissionNotFound":    {"zh-CN": "提交记录不存在", "en-US": "Submission not found", "code": 404},
	"CreateSubmissionError": {"zh-CN": "提交失败", "en-US": "Failed to create submission", "code": 500},
	"AlreadySolved":         {"zh-CN": "已提交过flag", "en-US": "Flag already submitted", "code": 200},
	"NotAllowSubmit":        {"zh-CN": "不允许提交flag", "en-US": "Not allowed to submit flag", "code": 400},

	"CreatePodError":     {"zh-CN": "创建Pod失败", "en-US": "Failed to create Pod", "code": 500},
	"GetPodError":        {"zh-CN": "获取Pod失败", "en-US": "Failed to get Pod", "code": 500},
	"CopyFileError":      {"zh-CN": "复制文件失败", "en-US": "Failed to copy file", "code": 500},
	"CreateServiceError": {"zh-CN": "创建Service失败", "en-US": "Failed to create Service", "code": 500},
	"DeletePodError":     {"zh-CN": "删除Pod失败", "en-US": "Failed to delete Pod", "code": 500},
	"DeleteServiceError": {"zh-CN": "删除Service失败", "en-US": "Failed to delete Service", "code": 500},

	"AppendUserToTeamError":      {"zh-CN": "添加用户到队伍失败", "en-US": "Failed to add user to team", "code": 500},
	"AppendTeamToContestError":   {"zh-CN": "添加队伍到赛事失败", "en-US": "Failed to add team to contest", "code": 500},
	"AppendContestToUserError":   {"zh-CN": "添加用户到赛事失败", "en-US": "Failed to add user to contest", "code": 500},
	"DeleteUserFromTeamError":    {"zh-CN": "删除用户从队伍失败", "en-US": "Failed to delete user from team", "code": 500},
	"DeleteUserFromContestError": {"zh-CN": "删除用户从赛事失败", "en-US": "Failed to delete user from contest", "code": 500},
	"DeleteAssociatedDataError":  {"zh-CN": "删除关联数据失败", "en-US": "Failed to delete associated data", "code": 500},
}

// I18N rewrite the response message
func I18N(key string, language string) (string, int) {
	if v, ok := resp[key]; !ok {
		return key, 400
	} else {
		return v[language].(string), v["code"].(int)
	}
}
