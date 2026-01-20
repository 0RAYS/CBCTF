package i18n

type BaseModel struct {
	// CreateError has fields: Error
	CreateError string

	// DeleteError has fields: Error
	DeleteError string

	// GetError has fields: Error
	GetError string

	NotFound string

	// UpdateError has fields: Error
	UpdateError string
}

type AssociationModel struct {
	// GetError has fields: Error
	GetError string

	// AppendError has fields: Error
	AppendError string

	// DeleteError has fields: Error
	DeleteError string
}

type Email struct {
	InvalidVerifyToken   string
	AlreadyVerifiedEmail string
	BaseModel
}

type File struct {
	// UnacceptedSuffix has fields: Accepted Suffix
	UnacceptedSuffix string

	BaseModel
}

type Setting struct {
	// InvalidType has fields: Key Type InvalidType
	InvalidType string

	BaseModel
}

type User struct {
	WeakPassword        string
	PasswordError       string
	NameOrPasswordError string
	SamePassword        string
	InvalidEmail        string
	UnverifiedEmail     string
	DuplicateEmail      string
	DuplicateUserName   string
	PasswordSame        string
	BaseModel
}

type Challenge struct {
	InvalidType string
	BaseModel
}

type ChallengeFlag struct {
	InvalidInjectType string
	BaseModel
}

type Contest struct {
	DuplicateName string
	CaptchaError  string
	IsComing      string
	IsRunning     string
	IsOver        string
	BaseModel
}

type ContestChallenge struct {
	AlreadySolved  string
	FlagNotMatch   string
	NotAllowSubmit string
	BaseModel
}

type ContestFlag struct {
	InvalidScoreType string
	BaseModel
}

type Docker struct {
	InvalidImage string
	BaseModel
}

type Notice struct {
	InvalidType string
	BaseModel
}

type Team struct {
	IsBanned           string
	IsFull             string
	DuplicateName      string
	DuplicateMember    string
	UserNotInTeam      string
	CaptainCannotLeave string
	CaptchaError       string
	BaseModel
}

type Traffic struct {
	ReadPcapError string
	PcapNotFound  string
	HasNoTraffic  string
	BaseModel
}

type Victim struct {
	HasMuchTime string
	Limited     string
	BaseModel
}

var Public = struct {
	Success string

	// UnknownError has fields: Error
	UnknownError string
}{
	Success:      "public.success",
	UnknownError: "public.unknownError",
}

var Request = struct {
	// BadRequest has fields: Error
	BadRequest string

	TooManyRequests string
	Unauthorized    string
	Forbidden       string
}{
	BadRequest:      "request.badRequest",
	TooManyRequests: "request.tooManyRequests",
	Unauthorized:    "request.unauthorized",
	Forbidden:       "request.forbidden",
}

var Redis = struct {
	// SetError has fields: Error
	SetError string

	// GetError has fields: Error
	GetError string

	// DeleteError has fields: Error
	DeleteError string
}{
	SetError:    "redis.setError",
	GetError:    "redis.getError",
	DeleteError: "redis.deleteError",
}

var Task = struct {
	// EnqueueError has fields: Error
	EnqueueError string
}{
	EnqueueError: "task.enqueueError",
}

var Model = struct {
	// NotUniqueKey has fields: Model Column
	NotUniqueKey string

	// UnsupportedKey has fields: Model Column
	UnsupportedKey string

	// UpdateDeadLock has fields: Model
	UpdateDeadLock string

	// DuplicateField has fields: Model Key Value
	DuplicateField string

	UserGroup         AssociationModel
	UserRole          AssociationModel
	GroupRole         AssociationModel
	RolePermission    AssociationModel
	Email             Email
	File              File
	Group             BaseModel
	OIDC              BaseModel
	Permission        BaseModel
	Request           BaseModel
	Role              BaseModel
	Setting           Setting
	Smtp              BaseModel
	User              User
	Admin             BaseModel
	Challenge         Challenge
	ChallengeFlag     ChallengeFlag
	Cheat             BaseModel
	Container         BaseModel
	Contest           Contest
	ContestChallenge  ContestChallenge
	ContestFlag       ContestFlag
	Device            BaseModel
	Docker            Docker
	EmailRecord       BaseModel
	Event             BaseModel
	Notice            Notice
	Oauth             BaseModel
	Pod               BaseModel
	Submission        BaseModel
	Team              Team
	TeamFlag          BaseModel
	Traffic           Traffic
	Victim            Victim
	Webhook           BaseModel
	WebhookHistory    BaseModel
	Namespace         BaseModel
	ConfigMap         BaseModel
	NetworkPolicy     BaseModel
	Service           BaseModel
	Job               BaseModel
	VPC               BaseModel
	Subnet            BaseModel
	VPCNatGateway     BaseModel
	EIP               BaseModel
	DNat              BaseModel
	SNat              BaseModel
	NetworkAttachment BaseModel
	IP                BaseModel
	PV                BaseModel
	PVC               BaseModel
	Endpoint          BaseModel
}{
	NotUniqueKey:   "model.notUniqueKey",
	UnsupportedKey: "model.unsupportedKey",
	UpdateDeadLock: "model.updateDeadLock",
	DuplicateField: "model.duplicateField",
	UserGroup: AssociationModel{
		GetError:    "model.userGroup.getError",
		AppendError: "model.userGroup.appendError",
		DeleteError: "model.userGroup.deleteError",
	},
	UserRole: AssociationModel{
		GetError:    "model.userRole.getError",
		AppendError: "model.userRole.appendError",
		DeleteError: "model.userRole.deleteError",
	},
	GroupRole: AssociationModel{
		GetError:    "model.groupRole.getError",
		AppendError: "model.groupRole.appendError",
		DeleteError: "model.groupRole.deleteError",
	},
	RolePermission: AssociationModel{
		GetError:    "model.rolePermission.getError",
		AppendError: "model.rolePermission.appendError",
		DeleteError: "model.rolePermission.deleteError",
	},
	Email: Email{
		"model.email.invalidVerifyToken",
		"model.email.alreadyVerifiedEmail",
		BaseModel{
			CreateError: "model.email.createError",
			DeleteError: "model.email.deleteError",
			GetError:    "model.email.getError",
			NotFound:    "model.email.notFound",
			UpdateError: "model.email.updateError",
		},
	},
	File: File{
		"model.file.unacceptedSuffix",
		BaseModel{
			CreateError: "model.file.CreateError",
			DeleteError: "model.file.DeleteError",
			GetError:    "model.file.getError",
			NotFound:    "model.file.notFound",
			UpdateError: "model.file.updateError",
		},
	},
	Group: BaseModel{
		CreateError: "model.group.createError",
		DeleteError: "model.group.deleteError",
		GetError:    "model.group.getError",
		NotFound:    "model.group.notFound",
		UpdateError: "model.group.updateError",
	},
	OIDC: BaseModel{
		CreateError: "model.oidc.createError",
		DeleteError: "model.oidc.deleteError",
		GetError:    "model.oidc.getError",
		NotFound:    "model.oidc.notFound",
		UpdateError: "model.oidc.updateError",
	},
	Permission: BaseModel{
		CreateError: "model.permission.createError",
		DeleteError: "model.permission.deleteError",
		GetError:    "model.permission.getError",
		NotFound:    "model.permission.notFound",
		UpdateError: "model.permission.updateError",
	},
	Request: BaseModel{
		"model.request.createError",
		"model.request.deleteError",
		"model.request.getError",
		"model.request.notFound",
		"model.request.updateError",
	},
	Role: BaseModel{
		"model.role.createError",
		"model.role.deleteError",
		"model.role.getError",
		"model.role.notFound",
		"model.role.updateError",
	},
	Setting: Setting{
		"model.setting.invalidType",
		BaseModel{
			"model.setting.createError",
			"model.setting.deleteError",
			"model.setting.getError",
			"model.setting.notFound",
			"model.setting.updateError",
		},
	},
	Smtp: BaseModel{
		"model.smtp.createError",
		"model.smtp.deleteError",
		"model.smtp.getError",
		"model.smtp.notFound",
		"model.smtp.updateError",
	},
	User: User{
		WeakPassword:        "model.user.weakPassword",
		PasswordError:       "model.user.passwordError",
		NameOrPasswordError: "model.user.nameOrPasswordError",
		SamePassword:        "model.user.samePassword",
		InvalidEmail:        "model.user.invalidEmail",
		UnverifiedEmail:     "model.user.unverifiedEmail",
		DuplicateEmail:      "model.user.duplicateEmail",
		DuplicateUserName:   "model.user.duplicateUserName",
		PasswordSame:        "model.user.passwordSame",
		BaseModel: BaseModel{
			CreateError: "model.user.createError",
			DeleteError: "model.user.deleteError",
			GetError:    "model.user.getError",
			NotFound:    "model.user.notFound",
			UpdateError: "model.user.updateError",
		},
	},
	Admin: BaseModel{
		CreateError: "model.admin.createError",
		DeleteError: "model.admin.deleteError",
		GetError:    "model.admin.getError",
		NotFound:    "model.admin.notFound",
		UpdateError: "model.admin.updateError",
	},
	Challenge: Challenge{
		InvalidType: "model.challenge.invalidType",
		BaseModel: BaseModel{
			CreateError: "model.challenge.createError",
			DeleteError: "model.challenge.deleteError",
			GetError:    "model.challenge.getError",
			NotFound:    "model.challenge.notFound",
			UpdateError: "model.challenge.updateError",
		},
	},
	ChallengeFlag: ChallengeFlag{
		InvalidInjectType: "model.challengeFlag.invalidInjectType",
		BaseModel: BaseModel{
			CreateError: "model.challengeFlag.createError",
			DeleteError: "model.challengeFlag.deleteError",
			GetError:    "model.challengeFlag.getError",
			NotFound:    "model.challengeFlag.notFound",
			UpdateError: "model.challengeFlag.updateError",
		},
	},
	Cheat: BaseModel{
		CreateError: "model.cheat.createError",
		DeleteError: "model.cheat.deleteError",
		GetError:    "model.cheat.getError",
		NotFound:    "model.cheat.notFound",
		UpdateError: "model.cheat.updateError",
	},
	Container: BaseModel{
		CreateError: "model.container.createError",
		DeleteError: "model.container.deleteError",
		GetError:    "model.container.getError",
		NotFound:    "model.container.notFound",
		UpdateError: "model.container.updateError",
	},
	Contest: Contest{
		DuplicateName: "model.contest.duplicateName",
		CaptchaError:  "model.contest.captchaError",
		IsComing:      "model.contest.isComing",
		IsRunning:     "model.contest.isRunning",
		IsOver:        "model.contest.isOver",
		BaseModel: BaseModel{
			CreateError: "model.contest.createError",
			DeleteError: "model.contest.deleteError",
			GetError:    "model.contest.getError",
			NotFound:    "model.contest.notFound",
			UpdateError: "model.contest.updateError",
		},
	},
	ContestChallenge: ContestChallenge{
		AlreadySolved:  "model.contestChallenge.alreadySolved",
		FlagNotMatch:   "model.contestChallenge.flagNotMatch",
		NotAllowSubmit: "model.contestChallenge.notAllowSubmit",
		BaseModel: BaseModel{
			CreateError: "model.contestChallenge.createError",
			DeleteError: "model.contestChallenge.deleteError",
			GetError:    "model.contestChallenge.getError",
			NotFound:    "model.contestChallenge.notFound",
			UpdateError: "model.contestChallenge.updateError",
		},
	},
	ContestFlag: ContestFlag{
		InvalidScoreType: "model.contestFlag.invalidScoreType",
		BaseModel: BaseModel{
			CreateError: "model.contestFlag.createError",
			DeleteError: "model.contestFlag.deleteError",
			GetError:    "model.contestFlag.getError",
			NotFound:    "model.contestFlag.notFound",
			UpdateError: "model.contestFlag.updateError",
		},
	},
	Device: BaseModel{
		CreateError: "model.device.createError",
		DeleteError: "model.device.deleteError",
		GetError:    "model.device.getError",
		NotFound:    "model.device.notFound",
		UpdateError: "model.device.updateError",
	},
	Docker: Docker{
		InvalidImage: "model.docker.invalidImage",
		BaseModel: BaseModel{
			CreateError: "model.docker.createError",
			DeleteError: "model.docker.deleteError",
			GetError:    "model.docker.getError",
			NotFound:    "model.docker.notFound",
			UpdateError: "model.docker.updateError",
		},
	},
	EmailRecord: BaseModel{
		CreateError: "model.emailRecord.createError",
		DeleteError: "model.emailRecord.deleteError",
		GetError:    "model.emailRecord.getError",
		NotFound:    "model.emailRecord.notFound",
		UpdateError: "model.emailRecord.updateError",
	},
	Event: BaseModel{
		CreateError: "model.event.createError",
		DeleteError: "model.event.deleteError",
		GetError:    "model.event.getError",
		NotFound:    "model.event.notFound",
		UpdateError: "model.event.updateError",
	},
	Notice: Notice{
		InvalidType: "model.notice.invalidType",
		BaseModel: BaseModel{
			CreateError: "model.notice.createError",
			DeleteError: "model.notice.deleteError",
			GetError:    "model.notice.getError",
			NotFound:    "model.notice.notFound",
			UpdateError: "model.notice.updateError",
		},
	},
	Oauth: BaseModel{
		CreateError: "model.oauth.createError",
		DeleteError: "model.oauth.deleteError",
		GetError:    "model.oauth.getError",
		NotFound:    "model.oauth.notFound",
		UpdateError: "model.oauth.updateError",
	},
	Pod: BaseModel{
		CreateError: "model.pod.createError",
		DeleteError: "model.pod.deleteError",
		GetError:    "model.pod.getError",
		NotFound:    "model.pod.notFound",
		UpdateError: "model.pod.updateError",
	},
	Submission: BaseModel{
		CreateError: "model.submission.createError",
		DeleteError: "model.submission.deleteError",
		GetError:    "model.submission.getError",
		NotFound:    "model.submission.notFound",
		UpdateError: "model.submission.updateError",
	},
	Team: Team{
		IsBanned:           "model.team.isBanned",
		IsFull:             "model.team.isFull",
		DuplicateName:      "model.team.duplicateName",
		DuplicateMember:    "model.team.duplicateMember",
		UserNotInTeam:      "model.team.userNotInTeam",
		CaptainCannotLeave: "model.team.captainCannotLeave",
		CaptchaError:       "model.team.captchaError",
		BaseModel: BaseModel{
			CreateError: "model.team.createError",
			DeleteError: "model.team.deleteError",
			GetError:    "model.team.getError",
			NotFound:    "model.team.notFound",
			UpdateError: "model.team.updateError",
		},
	},
	TeamFlag: BaseModel{
		CreateError: "model.teamFlag.createError",
		DeleteError: "model.teamFlag.deleteError",
		GetError:    "model.teamFlag.getError",
		NotFound:    "model.teamFlag.notFound",
		UpdateError: "model.teamFlag.updateError",
	},
	Traffic: Traffic{
		ReadPcapError: "model.traffic.readPcapError",
		PcapNotFound:  "model.traffic.pcapNotFound",
		HasNoTraffic:  "model.traffic.hasNoTraffic",
		BaseModel: BaseModel{
			CreateError: "model.traffic.createError",
			DeleteError: "model.traffic.deleteError",
			GetError:    "model.traffic.getError",
			NotFound:    "model.traffic.notFound",
			UpdateError: "model.traffic.updateError",
		},
	},
	Victim: Victim{
		HasMuchTime: "model.victim.hasMuchTime",
		Limited:     "model.victim.limited",
		BaseModel: BaseModel{
			CreateError: "model.victim.createError",
			DeleteError: "model.victim.deleteError",
			GetError:    "model.victim.getError",
			NotFound:    "model.victim.notFound",
			UpdateError: "model.victim.updateError",
		},
	},
	Webhook: BaseModel{
		CreateError: "model.webhook.createError",
		DeleteError: "model.webhook.deleteError",
		GetError:    "model.webhook.getError",
		NotFound:    "model.webhook.notFound",
		UpdateError: "model.webhook.updateError",
	},
	WebhookHistory: BaseModel{
		CreateError: "model.webhookHistory.createError",
		DeleteError: "model.webhookHistory.deleteError",
		GetError:    "model.webhookHistory.getError",
		NotFound:    "model.webhookHistory.notFound",
		UpdateError: "model.webhookHistory.updateError",
	},
	Namespace: BaseModel{
		CreateError: "model.namespace.createError",
		DeleteError: "model.namespace.deleteError",
		GetError:    "model.namespace.getError",
		NotFound:    "model.namespace.notFound",
	},
	ConfigMap: BaseModel{
		CreateError: "model.configMap.createError",
		DeleteError: "model.configMap.deleteError",
		GetError:    "model.configMap.getError",
		NotFound:    "model.configMap.notFound",
	},
	NetworkPolicy: BaseModel{
		CreateError: "model.networkPolicy.createError",
		DeleteError: "model.networkPolicy.deleteError",
		GetError:    "model.networkPolicy.getError",
		NotFound:    "model.networkPolicy.notFound",
	},
	Service: BaseModel{
		CreateError: "model.service.createError",
		DeleteError: "model.service.deleteError",
		GetError:    "model.service.getError",
		NotFound:    "model.service.notFound",
	},
	Job: BaseModel{
		CreateError: "model.job.createError",
		DeleteError: "model.job.deleteError",
		GetError:    "model.job.getError",
		NotFound:    "model.job.notFound",
	},
	VPC: BaseModel{
		CreateError: "model.vpc.createError",
		DeleteError: "model.vpc.deleteError",
		GetError:    "model.vpc.getError",
		NotFound:    "model.vpc.notFound",
	},
	Subnet: BaseModel{
		CreateError: "model.subnet.createError",
		DeleteError: "model.subnet.deleteError",
		GetError:    "model.subnet.getError",
		NotFound:    "model.subnet.notFound",
	},
	VPCNatGateway: BaseModel{
		CreateError: "model.vpcNatGateway.createError",
		DeleteError: "model.vpcNatGateway.deleteError",
		GetError:    "model.vpcNatGateway.getError",
		NotFound:    "model.vpcNatGateway.notFound",
	},
	EIP: BaseModel{
		CreateError: "model.eip.createError",
		DeleteError: "model.eip.deleteError",
		GetError:    "model.eip.getError",
		NotFound:    "model.eip.notFound",
	},
	DNat: BaseModel{
		CreateError: "model.dnat.createError",
		DeleteError: "model.dnat.deleteError",
		GetError:    "model.dnat.getError",
		NotFound:    "model.dnat.notFound",
	},
	SNat: BaseModel{
		CreateError: "model.snat.createError",
		DeleteError: "model.snat.deleteError",
		GetError:    "model.snat.getError",
		NotFound:    "model.snat.notFound",
	},
	NetworkAttachment: BaseModel{
		CreateError: "model.networkAttachment.createError",
		DeleteError: "model.networkAttachment.deleteError",
		GetError:    "model.networkAttachment.getError",
		NotFound:    "model.networkAttachment.notFound",
	},
	IP: BaseModel{
		CreateError: "model.ip.createError",
		DeleteError: "model.ip.deleteError",
		GetError:    "model.ip.getError",
		NotFound:    "model.ip.notFound",
	},
	PV: BaseModel{
		CreateError: "model.pv.createError",
		DeleteError: "model.pv.deleteError",
		GetError:    "model.pv.getError",
		NotFound:    "model.pv.notFound",
	},
	PVC: BaseModel{
		CreateError: "model.pvc.createError",
		DeleteError: "model.pvc.deleteError",
		GetError:    "model.pvc.getError",
		NotFound:    "model.pvc.notFound",
	},
	Endpoint: BaseModel{
		CreateError: "model.endpoint.createError",
		DeleteError: "model.endpoint.deleteError",
		GetError:    "model.endpoint.getError",
		NotFound:    "model.endpoint.notFound",
	},
}

var Association = struct {
	AppendUserToTeam      string
	AppendUserToContest   string
	DeleteUserFromTeam    string
	DeleteUserFromContest string
}{
	AppendUserToTeam:      "association.appendUserToTeam",
	AppendUserToContest:   "association.appendUserToContest",
	DeleteUserFromTeam:    "association.deleteUserFromTeam",
	DeleteUserFromContest: "association.deleteUserFromContest",
}

var Ranking = struct {
	UpdateError string
}{
	UpdateError: "ranking.updateError",
}

var EmailToken = struct {
	InvalidToken string
	SendError    string
}{
	InvalidToken: "emailToken.invalidToken",
	SendError:    "emailToken.sendError",
}

var FileOps = struct {
	CreateDirError           string
	ReadDirError             string
	InvalidDockerComposeYaml string
	CopyFileError            string
	ExecCommandError         string
	ZipError                 string
}{
	CreateDirError:           "fileOps.createDirError",
	ReadDirError:             "fileOps.readDirError",
	InvalidDockerComposeYaml: "fileOps.invalidDockerComposeYaml",
	CopyFileError:            "fileOps.copyFileError",
	ExecCommandError:         "fileOps.execCommandError",
	ZipError:                 "fileOps.zipError",
}

var K8s = struct {
	GetNodeListError string
}{
	GetNodeListError: "k8s.getNodeListError",
}
