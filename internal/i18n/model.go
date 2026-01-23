package i18n

type Challenge struct {
	InvalidType string
	EmptyImage  string
}

type ChallengeFlag struct {
	InvalidType string
}

type Cheat struct {
	InvalidType string
}

type Contest struct {
	CaptchaWrong    string
	DuplicateMember string
	IsComing        string
	IsRunning       string
	IsOver          string
}

type ContestFlag struct {
	InvalidType string
}

type Docker struct {
	// Error
	InvalidComposeYaml string
}

type Email struct {
	InvalidVerifyToken string
}

type File struct {
	NotAllowed string
	// Error
	ReadPcapError string
	// Error
	CreateDirError string
}

type Notice struct {
	InvalidType string
}

type Submission struct {
	NotAllowed string
}

type Team struct {
	Banned             string
	Full               string
	CaptchaWrong       string
	NotHasMember       string
	CaptainCannotLeave string
}

type TeamFlag struct {
	NotMatch      string
	AlreadySolved string
}

type User struct {
	WeakPassword      string
	SamePassword      string
	PasswordWrong     string
	NamePasswordWrong string
	UnverifiedEmail   string
	AlreadyVerified   string
}

type Victim struct {
	Limited     string
	HasMuchTime string
}

type Webhook struct {
	InvalidMethod string
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
	DeleteError   string
	Challenge     Challenge
	ChallengeFlag ChallengeFlag
	Cheat         Cheat
	Contest       Contest
	ContestFlag   ContestFlag
	Docker        Docker
	Email         Email
	File          File
	Notice        Notice
	Submission    Submission
	Team          Team
	TeamFlag      TeamFlag
	User          User
	Victim        Victim
	Webhook       Webhook
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
	},
	ChallengeFlag: ChallengeFlag{
		InvalidType: "model.challenge.invalidType",
	},
	Cheat: Cheat{
		InvalidType: "model.cheat.invalidType",
	},
	Contest: Contest{
		CaptchaWrong:    "model.contest.captchaWrong",
		DuplicateMember: "model.contest.duplicateMember",
		IsComing:        "model.contest.isComing",
		IsRunning:       "model.contest.isRunning",
		IsOver:          "model.contest.isOver",
	},
	ContestFlag: ContestFlag{
		InvalidType: "model.contestFlag.invalidType",
	},
	Docker: Docker{
		InvalidComposeYaml: "model.docker.invalidComposeYaml",
	},
	Email: Email{
		InvalidVerifyToken: "model.email.invalidVerifyToken",
	},
	File: File{
		NotAllowed:     "model.file.notAllowed",
		ReadPcapError:  "model.file.readPcapError",
		CreateDirError: "model.file.createDirError",
	},
	Notice: Notice{
		InvalidType: "model.notice.invalidType",
	},
	Submission: Submission{
		NotAllowed: "model.submission.notAllowed",
	},
	Team: Team{
		Banned:             "model.team.banned",
		Full:               "model.team.full",
		CaptchaWrong:       "model.team.captchaWrong",
		NotHasMember:       "model.team.notHasMember",
		CaptainCannotLeave: "model.team.captainCannotLeave",
	},
	TeamFlag: TeamFlag{
		NotMatch:      "model.teamFlag.notMatch",
		AlreadySolved: "model.teamFlag.alreadySolved",
	},
	User: User{
		WeakPassword:      "model.user.weakPassword",
		SamePassword:      "model.user.samePassword",
		PasswordWrong:     "model.user.passwordWrong",
		NamePasswordWrong: "model.user.namePasswordWrong",
		UnverifiedEmail:   "model.user.unverifiedEmail",
		AlreadyVerified:   "model.user.alreadyVerified",
	},
	Victim: Victim{
		Limited:     "model.victim.limited",
		HasMuchTime: "model.victim.hasMuchTime",
	},
	Webhook: Webhook{
		InvalidMethod: "model.webhook.invalidMethod",
	},
}
