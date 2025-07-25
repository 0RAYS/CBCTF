package i18n

var resp = map[string]map[string]any{
	Success:        {"zh-CN": "操作成功", "en-US": "Success", "code": 200},
	UnsupportedKey: {"zh-CN": "不支持的键", "en-US": "Unsupported key", "code": 400},
	DeadLock:       {"zh-CN": "数据写入失败过多", "en-US": "Database deadlock", "code": 500},

	BadRequest:      {"zh-CN": "参数缺失或错误", "en-US": "Bad request", "code": 400},
	Unauthorized:    {"zh-CN": "未登录", "en-US": "Unauthorized", "code": 401},
	Forbidden:       {"zh-CN": "禁止访问", "en-US": "Forbidden", "code": 403},
	TooManyRequests: {"zh-CN": "请求过快", "en-US": "Too many requests", "code": 429},
	UnknownError:    {"zh-CN": "未知错误", "en-US": "Unknown error", "code": 500},

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
	InvalidChallengeType: {"zh-CN": "无效的题目类型", "en-US": "Invalid challenge type", "code": 400},

	CreateChallengeFlagError: {"zh-CN": "创建题目flag失败", "en-US": "Create challenge flag failed", "code": 500},
	DeleteChallengeFlagError: {"zh-CN": "删除题目flag失败", "en-US": "Delete challenge flag failed", "code": 500},
	GetChallengeFlagError:    {"zh-CN": "获取题目flag失败", "en-US": "Get challenge flag failed", "code": 500},
	ChallengeFlagNotFound:    {"zh-CN": "题目flag不存在", "en-US": "Challenge flag not found", "code": 404},
	UpdateChallengeFlagError: {"zh-CN": "更新题目flag失败", "en-US": "Update challenge flag failed", "code": 500},

	CreateCheatError: {"zh-CN": "创建作弊记录失败", "en-US": "Create cheat failed", "code": 500},
	DeleteCheatError: {"zh-CN": "删除作弊记录失败", "en-US": "Delete cheat failed", "code": 500},
	GetCheatError:    {"zh-CN": "获取作弊记录失败", "en-US": "Get cheat failed", "code": 500},
	CheatNotFound:    {"zh-CN": "作弊记录不存在", "en-US": "Cheat not found", "code": 404},
	UpdateCheatError: {"zh-CN": "更新作弊记录失败", "en-US": "Update cheat failed", "code": 500},

	CreateContainerError: {"zh-CN": "创建容器失败", "en-US": "Create container failed", "code": 500},
	DeleteContainerError: {"zh-CN": "删除容器失败", "en-US": "Delete container failed", "code": 500},
	GetContainerError:    {"zh-CN": "获取容器失败", "en-US": "Get container failed", "code": 500},
	ContainerNotFound:    {"zh-CN": "容器不存在", "en-US": "Container not found", "code": 404},
	UpdateContainerError: {"zh-CN": "更新容器失败", "en-US": "Update container failed", "code": 500},

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
	AlreadySolved:               {"zh-CN": "题目已被解决", "en-US": "Challenge already solved", "code": 400},
	FlagNotMatch:                {"zh-CN": "flag错误", "en-US": "Flag does not match", "code": 400},
	NotAllowSubmit:              {"zh-CN": "不允许提交", "en-US": "Not allowed to submit", "code": 403},

	CreateContestFlagError: {"zh-CN": "创建比赛flag失败", "en-US": "Create contest flag failed", "code": 500},
	DeleteContestFlagError: {"zh-CN": "删除比赛flag失败", "en-US": "Delete contest flag failed", "code": 500},
	GetContestFlagError:    {"zh-CN": "获取比赛flag失败", "en-US": "Get contest flag failed", "code": 500},
	ContestFlagNotFound:    {"zh-CN": "比赛flag不存在", "en-US": "Contest flag not found", "code": 404},
	UpdateContestFlagError: {"zh-CN": "更新比赛flag失败", "en-US": "Update contest flag failed", "code": 500},
	InvalidScoreType:       {"zh-CN": "无效的分数类型", "en-US": "Invalid score type", "code": 500},

	CreateDeviceError: {"zh-CN": "创建设备失败", "en-US": "Create device failed", "code": 500},
	DeleteDeviceError: {"zh-CN": "删除设备失败", "en-US": "Delete device failed", "code": 500},
	GetDeviceError:    {"zh-CN": "获取设备失败", "en-US": "Get device failed", "code": 500},
	DeviceNotFound:    {"zh-CN": "设备不存在", "en-US": "Device not found", "code": 404},
	UpdateDeviceError: {"zh-CN": "更新设备失败", "en-US": "Update device failed", "code": 500},

	CreateDockerError:  {"zh-CN": "创建Docker失败", "en-US": "Create Docker failed", "code": 500},
	DeleteDockerError:  {"zh-CN": "删除Docker失败", "en-US": "Delete Docker failed", "code": 500},
	GetDockerError:     {"zh-CN": "获取Docker失败", "en-US": "Get Docker failed", "code": 500},
	DockerNotFound:     {"zh-CN": "Docker不存在", "en-US": "Docker not found", "code": 404},
	UpdateDockerError:  {"zh-CN": "更新Docker失败", "en-US": "Update Docker failed", "code": 500},
	InvalidDockerImage: {"zh-CN": "无效的Docker镜像", "en-US": "Invalid Docker image", "code": 400},

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
	InvalidNoticeType: {"zh-CN": "无效的通知类型", "en-US": "Invalid notice type", "code": 400},

	CreatePodError: {"zh-CN": "创建Pod失败", "en-US": "Create pod failed", "code": 500},
	DeletePodError: {"zh-CN": "删除Pod失败", "en-US": "Delete pod failed", "code": 500},
	GetPodError:    {"zh-CN": "获取Pod失败", "en-US": "Get pod failed", "code": 500},
	PodNotFound:    {"zh-CN": "Pod不存在", "en-US": "Pod not found", "code": 404},
	UpdatePodError: {"zh-CN": "更新Pod失败", "en-US": "Update pod failed", "code": 500},

	CreateRequestError: {"zh-CN": "创建请求失败", "en-US": "Create request failed", "code": 500},
	DeleteRequestError: {"zh-CN": "删除请求失败", "en-US": "Delete request failed", "code": 500},
	GetRequestError:    {"zh-CN": "获取请求失败", "en-US": "Get request failed", "code": 500},
	RequestNotFound:    {"zh-CN": "请求不存在", "en-US": "Request not found", "code": 404},
	UpdateRequestError: {"zh-CN": "更新请求失败", "en-US": "Update request failed", "code": 500},

	CreateSubmissionError: {"zh-CN": "创建提交记录失败", "en-US": "Record submission failed", "code": 500},
	DeleteSubmissionError: {"zh-CN": "删除提交记录失败", "en-US": "Delete submission failed", "code": 500},
	GetSubmissionError:    {"zh-CN": "获取提交记录失败", "en-US": "Get submission failed", "code": 500},
	SubmissionNotFound:    {"zh-CN": "提交记录不存在", "en-US": "Submission not found", "code": 404},
	UpdateSubmissionError: {"zh-CN": "更新提交记录失败", "en-US": "Update submission failed", "code": 500},

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

	CreateTeamFlagError: {"zh-CN": "创建战队flag失败", "en-US": "Create team flag failed", "code": 500},
	DeleteTeamFlagError: {"zh-CN": "删除战队flag失败", "en-US": "Delete team flag failed", "code": 500},
	GetTeamFlagError:    {"zh-CN": "获取战队flag失败", "en-US": "Get team flag failed", "code": 500},
	TeamFlagNotFound:    {"zh-CN": "战队flag不存在", "en-US": "Team flag not found", "code": 404},
	UpdateTeamFlagError: {"zh-CN": "更新战队flag失败", "en-US": "Update team flag failed", "code": 500},

	CreateTrafficError: {"zh-CN": "创建流量记录失败", "en-US": "Create traffic failed", "code": 500},
	DeleteTrafficError: {"zh-CN": "删除流量记录失败", "en-US": "Delete traffic failed", "code": 500},
	GetTrafficError:    {"zh-CN": "获取流量记录失败", "en-US": "Get traffic failed", "code": 500},
	TrafficNotFound:    {"zh-CN": "流量记录不存在", "en-US": "Traffic not found", "code": 404},
	UpdateTrafficError: {"zh-CN": "更新流量记录失败", "en-US": "Update traffic failed", "code": 500},
	ReadPcapError:      {"zh-CN": "读取PCAP文件失败", "en-US": "Read PCAP file failed", "code": 500},
	PcapNotFound:       {"zh-CN": "PCAP文件不存在", "en-US": "PCAP file not found", "code": 404},
	HasNoTraffic:       {"zh-CN": "没有流量记录", "en-US": "No traffic records found", "code": 404},

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

	CreateVictimError: {"zh-CN": "创建靶机失败", "en-US": "Create victim failed", "code": 500},
	DeleteVictimError: {"zh-CN": "关闭靶机失败", "en-US": "Delete victim failed", "code": 500},
	GetVictimError:    {"zh-CN": "获取靶机失败", "en-US": "Get victim failed", "code": 500},
	VictimNotFound:    {"zh-CN": "靶机不存在", "en-US": "Victim not found", "code": 404},
	UpdateVictimError: {"zh-CN": "更新靶机失败", "en-US": "Update victim failed", "code": 500},
	HasMuchTime:       {"zh-CN": "距容器关闭20分钟内才可延长时间", "en-US": "Can only extend time within 20 minutes before the container closes", "code": 400},

	AppendUserToTeamError:      {"zh-CN": "添加用户到战队失败", "en-US": "Append user to team failed", "code": 500},
	AppendUserToContestError:   {"zh-CN": "添加用户到比赛失败", "en-US": "Append user to contest failed", "code": 500},
	DeleteUserFromTeamError:    {"zh-CN": "从战队中删除用户失败", "en-US": "Delete user from team failed", "code": 500},
	DeleteUserFromContestError: {"zh-CN": "从比赛中删除用户失败", "en-US": "Delete user from contest failed", "code": 500},

	UpdateRankingError: {"zh-CN": "更新排名失败", "en-US": "Update ranking failed", "code": 500},

	SetEmailVerifyTokenError: {"zh-CN": "设置邮箱验证令牌失败", "en-US": "Set email verify token failed", "code": 500},
	GetEmailVerifyTokenError: {"zh-CN": "获取邮箱验证令牌失败", "en-US": "Get email verify token failed", "code": 500},
	DelEmailVerifyTokenError: {"zh-CN": "删除邮箱验证令牌失败", "en-US": "Delete email verify token failed", "code": 500},
	InvalidEmailVerifyToken:  {"zh-CN": "无效的邮箱验证令牌", "en-US": "Invalid email verify token", "code": 400},
	SendEmailError:           {"zh-CN": "发送邮件失败", "en-US": "Send email failed", "code": 500},
	RedisError:               {"zh-CN": "缓存失败", "en-US": "Redis operation failed", "code": 500},

	CreateDirError:                 {"zh-CN": "创建题目根目录失败", "en-US": "Create directory failed", "code": 500},
	ReadDirError:                   {"zh-CN": "读取题目根目录失败", "en-US": "Read directory failed", "code": 500},
	InvalidFileName:                {"zh-CN": "无效的文件名", "en-US": "Invalid file name", "code": 400},
	InvalidDockerComposeYaml:       {"zh-CN": "无效的Docker Compose YAML文件", "en-US": "Invalid Docker Compose YAML file", "code": 400},
	InvalidChallengeFlagInjectType: {"zh-CN": "不支持的flag注入方式", "en-US": "Invalid challenge flag inject type", "code": 400},
	CopyFileError:                  {"zh-CN": "复制文件失败", "en-US": "Copy file failed", "code": 500},
	ExecCommandError:               {"zh-CN": "执行命令失败", "en-US": "Execute command failed", "code": 500},
	ZipError:                       {"zh-CN": "压缩文件失败", "en-US": "Zip file failed", "code": 500},

	CreateNamespaceError: {"zh-CN": "创建Namespace失败", "en-US": "Create namespace failed", "code": 500},
	DeleteNamespaceError: {"zh-CN": "删除Namespace失败", "en-US": "Delete namespace failed", "code": 500},
	GetNamespaceError:    {"zh-CN": "获取Namespace失败", "en-US": "Get namespace failed", "code": 500},
	NamespaceNotFound:    {"zh-CN": "Namespace不存在", "en-US": "Namespace not found", "code": 404},

	CreateConfigMapError: {"zh-CN": "创建ConfigMap失败", "en-US": "Create config map failed", "code": 500},
	DeleteConfigMapError: {"zh-CN": "删除ConfigMap失败", "en-US": "Delete config map failed", "code": 500},
	GetConfigMapError:    {"zh-CN": "获取ConfigMap失败", "en-US": "Get config map failed", "code": 500},
	ConfigMapNotFound:    {"zh-CN": "ConfigMap不存在", "en-US": "Config map not found", "code": 404},

	CreateNetworkPolicyError: {"zh-CN": "创建NetworkPolicy失败", "en-US": "Create network policy failed", "code": 500},
	DeleteNetworkPolicyError: {"zh-CN": "删除NetworkPolicy失败", "en-US": "Delete network policy failed", "code": 500},
	GetNetworkPolicyError:    {"zh-CN": "获取NetworkPolicy失败", "en-US": "Get network policy failed", "code": 500},
	NetworkPolicyNotFound:    {"zh-CN": "NetworkPolicy不存在", "en-US": "Network policy not found", "code": 404},

	CreateServiceError: {"zh-CN": "创建Service失败", "en-US": "Create service failed", "code": 500},
	DeleteServiceError: {"zh-CN": "删除Service失败", "en-US": "Delete service failed", "code": 500},
	GetServiceError:    {"zh-CN": "获取Service失败", "en-US": "Get service failed", "code": 500},
	ServiceNotFound:    {"zh-CN": "Service不存在", "en-US": "Service not found", "code": 404},

	CreateJobError: {"zh-CN": "创建Job失败", "en-US": "Create Job failed", "code": 500},
	DeleteJobError: {"zh-CN": "删除Job失败", "en-US": "Delete Job failed", "code": 500},
	GetJobError:    {"zh-CN": "获取Job失败", "en-US": "Get Job failed", "code": 500},
	JobNotFound:    {"zh-CN": "Job不存在", "en-US": "Job not found", "code": 404},

	CreateVPCError: {"zh-CN": "创建VPC失败", "en-US": "Create VPC failed", "code": 500},
	DeleteVPCError: {"zh-CN": "删除VPC失败", "en-US": "Delete VPC failed", "code": 500},
	VPCNotFound:    {"zh-CN": "VPC不存在", "en-US": "VPC not found", "code": 404},
	GetVPCError:    {"zh-CN": "获取VPC失败", "en-US": "Get VPC failed", "code": 500},

	CreateSubnetError: {"zh-CN": "创建子网失败", "en-US": "Create subnet failed", "code": 500},
	DeleteSubnetError: {"zh-CN": "删除子网失败", "en-US": "Delete subnet failed", "code": 500},
	SubnetNotFound:    {"zh-CN": "子网不存在", "en-US": "Subnet not found", "code": 404},
	GetSubnetError:    {"zh-CN": "获取子网失败", "en-US": "Get subnet failed", "code": 500},

	CreateVPCNatGatewayError: {"zh-CN": "创建VPC NAT网关失败", "en-US": "Create VPC NAT gateway failed", "code": 500},
	DeleteVPCNatGatewayError: {"zh-CN": "删除VPC NAT网关失败", "en-US": "Delete VPC NAT gateway failed", "code": 500},
	GetVPCNatGatewayError:    {"zh-CN": "获取VPC NAT网关失败", "en-US": "Get VPC NAT gateway failed", "code": 500},
	VPCNatGatewayNotFound:    {"zh-CN": "VPC NAT网关不存在", "en-US": "VPC NAT gateway not found", "code": 404},

	CreateEIPError: {"zh-CN": "创建弹性IP失败", "en-US": "Create EIP failed", "code": 500},
	DeleteEIPError: {"zh-CN": "删除弹性IP失败", "en-US": "Delete EIP failed", "code": 500},
	GetEIPError:    {"zh-CN": "获取弹性IP失败", "en-US": "Get EIP failed", "code": 500},
	EIPNotFound:    {"zh-CN": "弹性IP不存在", "en-US": "EIP not found", "code": 404},

	CreateDNatError: {"zh-CN": "创建DNat规则失败", "en-US": "Create DNat failed", "code": 500},
	DeleteDNatError: {"zh-CN": "删除DNat规则失败", "en-US": "Delete DNat failed", "code": 500},
	GetDNatError:    {"zh-CN": "获取DNat规则失败", "en-US": "Get DNat failed", "code": 500},
	DNatNotFound:    {"zh-CN": "DNat规则不存在", "en-US": "DNat failed", "code": 404},

	CreateSNatError: {"zh-CN": "创建SNat规则失败", "en-US": "Create SNat failed", "code": 500},
	DeleteSNatError: {"zh-CN": "删除SNat规则失败", "en-US": "Delete SNat failed", "code": 500},
	GetSNatError:    {"zh-CN": "获取SNat规则失败", "en-US": "Get SNat failed", "code": 500},
	SNatNotFound:    {"zh-CN": "SNat规则不存在", "en-US": "SNat failed", "code": 404},

	CreateNetAttError: {"zh-CN": "创建附属网卡失败", "en-US": "Create network attachment failed", "code": 500},
	DeleteNetAttError: {"zh-CN": "删除附属网卡失败", "en-US": "Delete network attachment failed", "code": 500},
	GetNetAttError:    {"zh-CN": "获取附属网卡失败", "en-US": "Get network attachment failed", "code": 500},
	NetAttNotFound:    {"zh-CN": "附属网卡不存在", "en-US": "Network attachment not found", "code": 404},

	CreateIPError: {"zh-CN": "创建IP失败", "en-US": "Create IP failed", "code": 500},
	DeleteIPError: {"zh-CN": "删除IP失败", "en-US": "Delete IP failed", "code": 500},
	GetIPError:    {"zh-CN": "获取IP失败", "en-US": "Get IP failed", "code": 500},
	IPNotFound:    {"zh-CN": "IP不存在", "en-US": "IP not found", "code": 404},

	GetNodeListError: {"zh-CN": "获取K8S节点失败", "en-US": "Get node list failed", "code": 500},
}

// I18N 获取翻译与状态码, 非http响应状态码
func I18N(key string, language string) (string, int) {
	if v, ok := resp[key]; !ok {
		return key, 500
	} else {
		if language == "origin" {
			return key, v["code"].(int)
		}
		return v[language].(string), v["code"].(int)
	}
}
