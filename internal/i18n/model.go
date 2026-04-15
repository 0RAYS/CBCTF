package i18n

type Challenge struct {
	InvalidType string
	EmptyImage  string
	GetError    string
	DeleteError string
}

type ChallengeFlag struct {
	InvalidType string
	DeleteError string
}

type Cheat struct {
	InvalidType string
	CreateError string
	DeleteError string
}

type Contest struct {
	CaptchaWrong    string
	DuplicateMember string
	IsComing        string
	IsRunning       string
	IsOver          string
	GetError        string
	DeleteError     string
	NotFound        string
}

type ContestChallenge struct {
	GetError    string
	DeleteError string
}

type ContestFlag struct {
	InvalidType string
	GetError    string
	DeleteError string
	NotFound    string
}

type Device struct {
	CreateError string
}

type Docker struct {
	// Error
	InvalidComposeYaml string
	DeleteError        string
}

type Email struct {
	InvalidVerifyToken string
	CreateError        string
}

type File struct {
	NotAllowed string
	// Error
	ReadPcapError string
	// Error
	CreateDirError string
	CreateError    string
	DeleteError    string
	NotFound       string
}

type Generator struct {
	NotAvailable string
	NotStoppable string
	DeleteError  string
}

type Group struct {
	CannotUpdateDefault string
	CannotDeleteDefault string
	GetError            string
	DeleteError         string
}

type Notice struct {
	InvalidType string
}

type Permission struct {
	GetError string
	NotFound string
}

type Pod struct {
	DeleteError string
}

type Request struct {
	GetError string
}

type Role struct {
	CannotUpdateDefault string
	CannotDeleteDefault string
	GetError            string
	DeleteError         string
}

type RolePermission struct {
	CreateError string
	DeleteError string
}

type Setting struct {
	InvalidType string
	DeadLock    string
	UpdateError string
}

type Submission struct {
	NotAllowed  string
	GetError    string
	DeleteError string
}

type Team struct {
	Banned             string
	Full               string
	CaptchaWrong       string
	NotHasMember       string
	CaptainCannotLeave string
	GetError           string
	NotFound           string
	DeleteError        string
}

type TeamFlag struct {
	NotMatch      string
	AlreadySolved string
	GetError      string
	NotFound      string
}

type Traffic struct {
	GetError string
}

type User struct {
	WeakPassword       string
	SamePassword       string
	PasswordWrong      string
	NamePasswordWrong  string
	UnverifiedEmail    string
	AlreadyVerified    string
	InContest          string
	NotAllowedRegister string
	GetError           string
	DeleteError        string
}

type UserContest struct {
	CreateError string
	DeleteError string
}

type UserGroup struct {
	GetError    string
	CreateError string
	DeleteError string
}

type UserTeam struct {
	CreateError string
	DeleteError string
}

type Victim struct {
	Limited      string
	HasMuchTime  string
	NotStoppable string
	DeleteError  string
}

type Webhook struct {
	InvalidMethod    string
	NotAllowedTarget string
}

type WebhookHistory struct {
	CreateError string
}

var Model = struct {
	// Model Key
	DuplicateKeyValue string
	// Model Key
	NotUniqueKey string
	// Model
	DeadLock string
	// Model Error
	CreateError string
	// Model Error
	GetError string
	// Model
	NotFound string
	// Model Error
	UpdateError string
	// Model Error
	DeleteError      string
	Challenge        Challenge
	ChallengeFlag    ChallengeFlag
	Cheat            Cheat
	Contest          Contest
	ContestChallenge ContestChallenge
	ContestFlag      ContestFlag
	Device           Device
	Docker           Docker
	Email            Email
	File             File
	Generator        Generator
	Group            Group
	Notice           Notice
	Permission       Permission
	Pod              Pod
	Request          Request
	Role             Role
	RolePermission   RolePermission
	Setting          Setting
	Submission       Submission
	Team             Team
	TeamFlag         TeamFlag
	Traffic          Traffic
	User             User
	UserContest      UserContest
	UserGroup        UserGroup
	UserTeam         UserTeam
	Victim           Victim
	Webhook          Webhook
	WebhookHistory   WebhookHistory
}{
	DuplicateKeyValue: "model.duplicateKeyValue",
	NotUniqueKey:      "model.notUniqueKey",
	DeadLock:          "model.deadLock",
	CreateError:       "model.createError",
	GetError:          "model.getError",
	NotFound:          "model.notFound",
	UpdateError:       "model.updateError",
	DeleteError:       "model.deleteError",
	Challenge: Challenge{
		InvalidType: "model.challenge.invalidType",
		EmptyImage:  "model.challenge.emptyImage",
		GetError:    "model.challenge.getError",
		DeleteError: "model.challenge.deleteError",
	},
	ChallengeFlag: ChallengeFlag{
		InvalidType: "model.challenge.invalidType",
		DeleteError: "model.challengeFlag.deleteError",
	},
	Cheat: Cheat{
		InvalidType: "model.cheat.invalidType",
		CreateError: "model.cheat.createError",
		DeleteError: "model.cheat.deleteError",
	},
	Contest: Contest{
		CaptchaWrong:    "model.contest.captchaWrong",
		DuplicateMember: "model.contest.duplicateMember",
		IsComing:        "model.contest.isComing",
		IsRunning:       "model.contest.isRunning",
		IsOver:          "model.contest.isOver",
		GetError:        "model.contest.getError",
		DeleteError:     "model.contest.deleteError",
		NotFound:        "model.contest.notFound",
	},
	ContestChallenge: ContestChallenge{
		GetError:    "model.contestChallenge.getError",
		DeleteError: "model.contestChallenge.deleteError",
	},
	ContestFlag: ContestFlag{
		InvalidType: "model.contestFlag.invalidType",
		GetError:    "model.contestFlag.getError",
		DeleteError: "model.contestFlag.deleteError",
		NotFound:    "model.contestFlag.notFound",
	},
	Device: Device{
		CreateError: "model.device.createError",
	},
	Docker: Docker{
		InvalidComposeYaml: "model.docker.invalidComposeYaml",
		DeleteError:        "model.docker.deleteError",
	},
	Email: Email{
		InvalidVerifyToken: "model.email.invalidVerifyToken",
		CreateError:        "model.email.createError",
	},
	File: File{
		NotAllowed:     "model.file.notAllowed",
		ReadPcapError:  "model.file.readPcapError",
		CreateDirError: "model.file.createDirError",
		CreateError:    "model.file.createError",
		DeleteError:    "model.file.deleteError",
		NotFound:       "model.file.notFound",
	},
	Generator: Generator{
		NotAvailable: "model.generator.notAvailable",
		NotStoppable: "model.generator.notStoppable",
		DeleteError:  "model.generator.deleteError",
	},
	Group: Group{
		CannotUpdateDefault: "model.group.cannotUpdateDefault",
		CannotDeleteDefault: "model.group.cannotDeleteDefault",
		GetError:            "model.group.getError",
		DeleteError:         "model.group.deleteError",
	},
	Notice: Notice{
		InvalidType: "model.notice.invalidType",
	},
	Permission: Permission{
		GetError: "model.permission.getError",
		NotFound: "model.permission.notFound",
	},
	Pod: Pod{
		DeleteError: "model.pod.deleteError",
	},
	Request: Request{
		GetError: "model.request.getError",
	},
	Role: Role{
		CannotUpdateDefault: "model.role.cannotUpdateDefault",
		CannotDeleteDefault: "model.role.cannotDeleteDefault",
		GetError:            "model.role.getError",
		DeleteError:         "model.role.deleteError",
	},
	RolePermission: RolePermission{
		CreateError: "model.rolePermission.createError",
		DeleteError: "model.rolePermission.deleteError",
	},
	Setting: Setting{
		InvalidType: "model.setting.invalidType",
		DeadLock:    "model.setting.deadLock",
		UpdateError: "model.setting.updateError",
	},
	Submission: Submission{
		NotAllowed:  "model.submission.notAllowed",
		GetError:    "model.submission.getError",
		DeleteError: "model.submission.deleteError",
	},
	Team: Team{
		Banned:             "model.team.banned",
		Full:               "model.team.full",
		CaptchaWrong:       "model.team.captchaWrong",
		NotHasMember:       "model.team.notHasMember",
		CaptainCannotLeave: "model.team.captainCannotLeave",
		GetError:           "model.team.getError",
		NotFound:           "model.team.notFound",
		DeleteError:        "model.team.deleteError",
	},
	TeamFlag: TeamFlag{
		NotMatch:      "model.teamFlag.notMatch",
		AlreadySolved: "model.teamFlag.alreadySolved",
		GetError:      "model.teamFlag.getError",
		NotFound:      "model.teamFlag.notFound",
	},
	Traffic: Traffic{
		GetError: "model.traffic.getError",
	},
	User: User{
		WeakPassword:       "model.user.weakPassword",
		SamePassword:       "model.user.samePassword",
		PasswordWrong:      "model.user.passwordWrong",
		NamePasswordWrong:  "model.user.namePasswordWrong",
		UnverifiedEmail:    "model.user.unverifiedEmail",
		AlreadyVerified:    "model.user.alreadyVerified",
		InContest:          "model.user.inContest",
		NotAllowedRegister: "model.user.notAllowedRegister",
		GetError:           "model.user.getError",
		DeleteError:        "model.user.deleteError",
	},
	UserContest: UserContest{
		CreateError: "model.userContest.createError",
		DeleteError: "model.userContest.deleteError",
	},
	UserGroup: UserGroup{
		GetError:    "model.userGroup.getError",
		CreateError: "model.userGroup.createError",
		DeleteError: "model.userGroup.deleteError",
	},
	UserTeam: UserTeam{
		CreateError: "model.userTeam.createError",
		DeleteError: "model.userTeam.deleteError",
	},
	Victim: Victim{
		Limited:      "model.victim.limited",
		HasMuchTime:  "model.victim.hasMuchTime",
		NotStoppable: "model.victim.notStoppable",
		DeleteError:  "model.victim.deleteError",
	},
	Webhook: Webhook{
		InvalidMethod:    "model.webhook.invalidMethod",
		NotAllowedTarget: "model.webhook.notAllowedTarget",
	},
	WebhookHistory: WebhookHistory{
		CreateError: "model.webhookHistory.createError",
	},
}
