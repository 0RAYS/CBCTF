package i18n

import (
	"embed"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var LocaleFS embed.FS

// Manager i18n管理器
type Manager struct {
	bundle *i18n.Bundle
}

// NewManager 创建新的i18n管理器
func NewManager() *Manager {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 加载所有语言文件
	loadLocales(bundle)

	return &Manager{
		bundle: bundle,
	}
}

// loadLocales 加载所有语言文件
func loadLocales(bundle *i18n.Bundle) {
	// 加载中文语言文件
	if _, err := bundle.LoadMessageFileFS(LocaleFS, "locales/zh-CN.toml"); err != nil {
		// 如果加载失败，记录错误但不中断程序
		// 这里可以添加日志记录
	}

	// 加载英文语言文件
	if _, err := bundle.LoadMessageFileFS(LocaleFS, "locales/en-US.toml"); err != nil {
		// 如果加载失败，记录错误但不中断程序
		// 这里可以添加日志记录
	}
}

// GetLocalizer 根据语言偏好获取本地化器
func (m *Manager) GetLocalizer(languages ...string) *i18n.Localizer {
	return i18n.NewLocalizer(m.bundle, languages...)
}

// T 翻译消息
func (m *Manager) T(messageID string, languages ...string) string {
	localizer := m.GetLocalizer(languages...)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID
	}
	return message
}

// TWithData 翻译消息并传入数据
func (m *Manager) TWithData(messageID string, data map[string]interface{}, languages ...string) string {
	localizer := m.GetLocalizer(languages...)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID
	}
	return message
}

// TPlural 翻译复数消息
func (m *Manager) TPlural(messageID string, count int, languages ...string) string {
	localizer := m.GetLocalizer(languages...)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:   messageID,
		PluralCount: count,
	})
	if err != nil {
		return messageID
	}
	return message
}

// ParseAcceptLanguage 解析Accept-Language头
func (m *Manager) ParseAcceptLanguage(acceptLanguage string) []string {
	if acceptLanguage == "" {
		return []string{"zh-CN"}
	}

	// 解析语言偏好
	parts := strings.Split(acceptLanguage, ",")
	var languages []string

	for _, part := range parts {
		lang := strings.TrimSpace(strings.Split(part, ";")[0])
		if lang == "en-US" || strings.HasPrefix(lang, "en") {
			languages = append(languages, "en-US")
		} else if lang == "origin" {
			languages = append(languages, "origin")
		} else {
			languages = append(languages, "zh-CN")
		}
	}

	if len(languages) == 0 {
		languages = append(languages, "zh-CN")
	}

	return languages
}

// GetResponse 获取带状态码的响应
func (m *Manager) GetResponse(messageID string, languages ...string) (string, int) {
	// 从消息ID中提取状态码
	code := extractStatusCode(messageID)

	// 如果是origin语言，返回原始消息ID
	if len(languages) > 0 && languages[0] == "origin" {
		return messageID, code
	}

	message := m.T(messageID, languages...)
	return message, code
}

// extractStatusCode 从消息ID中提取状态码
func extractStatusCode(messageID string) int {
	// 状态码映射表
	statusCodes := map[string]int{
		"Success":         200,
		"BadRequest":      400,
		"Unauthorized":    401,
		"Forbidden":       403,
		"TooManyRequests": 429,
		"UnknownError":    500,

		// 404错误
		"AdminNotFound": 404, "ChallengeNotFound": 404, "ChallengeFlagNotFound": 404, "CheatNotFound": 404,
		"ContainerNotFound": 404, "ContestNotFound": 404, "ContestChallengeNotFound": 404, "ContestFlagNotFound": 404,
		"DeviceNotFound": 404, "DockerNotFound": 404, "EventNotFound": 404, "FileNotFound": 404, "NoticeNotFound": 404,
		"PodNotFound": 404, "RequestNotFound": 404, "SubmissionNotFound": 404, "TeamNotFound": 404, "TeamFlagNotFound": 404,
		"TrafficNotFound": 404, "UserNotFound": 404, "VictimNotFound": 404, "PcapNotFound": 404, "NamespaceNotFound": 404,
		"ConfigMapNotFound": 404, "NetworkPolicyNotFound": 404, "ServiceNotFound": 404, "JobNotFound": 404, "VPCNotFound": 404,
		"SubnetNotFound": 404, "VPCNatGatewayNotFound": 404, "EIPNotFound": 404, "DNatNotFound": 404, "SNatNotFound": 404,
		"NetAttNotFound": 404, "IPNotFound": 404, "PVNotFound": 404, "PVCNotFound": 404,

		// 400错误
		"UnsupportedKey": 400, "InvalidChallengeType": 400, "DuplicateContestName": 400, "ContestCaptchaError": 400,
		"ContestIsComing": 400, "ContestIsRunning": 400, "ContestIsOver": 400, "AlreadySolved": 400, "FlagNotMatch": 400,
		"NotAllowSubmit": 400, "InvalidScoreType": 400, "InvalidDockerImage": 400, "FileNotAllowed": 400,
		"InvalidNoticeType": 400, "TeamIsFull": 400, "DuplicateTeamName": 400, "DuplicateMember": 400,
		"UserNotInTeam": 400, "CaptainCannotLeave": 400, "TeamCaptchaError": 400, "InvalidEmail": 400, "UnverifiedEmail": 400,
		"DuplicateEmail": 400, "DuplicateUserName": 400, "WeakPassword": 400, "PasswordSame": 400, "HasMuchTime": 400,
		"InvalidEmailVerifyToken": 400, "InvalidFileName": 400, "InvalidDockerComposeYaml": 400,
		"InvalidChallengeFlagInjectType": 400,

		// 401错误
		"NameOrPasswordError": 401, "PasswordError": 401,

		// 403错误
		"TeamIsBanned": 403,

		// 500错误（默认）
		"DeadLock": 500, "CreateAdminError": 500, "DeleteAdminError": 500, "GetAdminError": 500, "UpdateAdminError": 500,
		"CreateChallengeError": 500, "DeleteChallengeError": 500, "GetChallengeError": 500, "UpdateChallengeError": 500,
		"CreateChallengeFlagError": 500, "DeleteChallengeFlagError": 500, "GetChallengeFlagError": 500, "UpdateChallengeFlagError": 500,
		"CreateCheatError": 500, "DeleteCheatError": 500, "GetCheatError": 500, "UpdateCheatError": 500,
		"CreateContainerError": 500, "DeleteContainerError": 500, "GetContainerError": 500, "UpdateContainerError": 500,
		"CreateContestError": 500, "DeleteContestError": 500, "GetContestError": 500, "UpdateContestError": 500,
		"CreateContestChallengeError": 500, "DeleteContestChallengeError": 500, "GetContestChallengeError": 500, "UpdateContestChallengeError": 500,
		"CreateContestFlagError": 500, "DeleteContestFlagError": 500, "GetContestFlagError": 500, "UpdateContestFlagError": 500,
		"CreateDeviceError": 500, "DeleteDeviceError": 500, "GetDeviceError": 500, "UpdateDeviceError": 500,
		"CreateDockerError": 500, "DeleteDockerError": 500, "GetDockerError": 500, "UpdateDockerError": 500,
		"CreateEventError": 500, "DeleteEventError": 500, "GetEventError": 500, "UpdateEventError": 500,
		"CreateFileError": 500, "DeleteFileError": 500, "GetFileError": 500, "UpdateFileError": 500,
		"CreateNoticeError": 500, "DeleteNoticeError": 500, "GetNoticeError": 500, "UpdateNoticeError": 500,
		"CreatePodError": 500, "DeletePodError": 500, "GetPodError": 500, "UpdatePodError": 500,
		"CreateRequestError": 500, "DeleteRequestError": 500, "GetRequestError": 500, "UpdateRequestError": 500,
		"CreateSubmissionError": 500, "DeleteSubmissionError": 500, "GetSubmissionError": 500, "UpdateSubmissionError": 500,
		"CreateTeamError": 500, "DeleteTeamError": 500, "GetTeamError": 500, "UpdateTeamError": 500,
		"CreateTeamFlagError": 500, "DeleteTeamFlagError": 500, "GetTeamFlagError": 500, "UpdateTeamFlagError": 500,
		"CreateTrafficError": 500, "DeleteTrafficError": 500, "GetTrafficError": 500, "UpdateTrafficError": 500,
		"CreateUserError": 500, "DeleteUserError": 500, "GetUserError": 500, "UpdateUserError": 500,
		"CreateVictimError": 500, "DeleteVictimError": 500, "GetVictimError": 500, "UpdateVictimError": 500,
		"AppendUserToTeamError": 500, "AppendUserToContestError": 500, "DeleteUserFromTeamError": 500, "DeleteUserFromContestError": 500,
		"UpdateRankingError": 500, "SetEmailVerifyTokenError": 500, "GetEmailVerifyTokenError": 500, "DelEmailVerifyTokenError": 500,
		"SendEmailError": 500, "RedisError": 500, "CreateDirError": 500, "ReadDirError": 500, "CopyFileError": 500,
		"ExecCommandError": 500, "ZipError": 500, "CreateNamespaceError": 500, "DeleteNamespaceError": 500, "GetNamespaceError": 500,
		"CreateConfigMapError": 500, "DeleteConfigMapError": 500, "GetConfigMapError": 500, "CreateNetworkPolicyError": 500,
		"DeleteNetworkPolicyError": 500, "GetNetworkPolicyError": 500, "CreateServiceError": 500, "DeleteServiceError": 500,
		"GetServiceError": 500, "CreateJobError": 500, "DeleteJobError": 500, "GetJobError": 500, "CreateVPCError": 500,
		"DeleteVPCError": 500, "GetVPCError": 500, "CreateSubnetError": 500, "DeleteSubnetError": 500, "GetSubnetError": 500,
		"CreateVPCNatGatewayError": 500, "DeleteVPCNatGatewayError": 500, "GetVPCNatGatewayError": 500, "CreateEIPError": 500,
		"DeleteEIPError": 500, "GetEIPError": 500, "CreateDNatError": 500, "DeleteDNatError": 500, "GetDNatError": 500,
		"CreateSNatError": 500, "DeleteSNatError": 500, "GetSNatError": 500, "CreateNetAttError": 500, "DeleteNetAttError": 500,
		"GetNetAttError": 500, "CreateIPError": 500, "DeleteIPError": 500, "GetIPError": 500, "CreatePVError": 500,
		"DeletePVError": 500, "GetPVError": 500, "CreatePVCError": 500, "DeletePVCError": 500, "GetPVCError": 500,
		"GetNodeListError": 500, "ReadPcapError": 500, "HasNoTraffic": 500,
	}

	if code, exists := statusCodes[messageID]; exists {
		return code
	}

	// 默认返回500
	return 500
}

// 全局管理器实例
var GlobalManager *Manager

// Init 初始化全局i18n管理器
func Init() {
	GlobalManager = NewManager()
}

// T 全局翻译函数
func T(messageID string, languages ...string) string {
	if GlobalManager == nil {
		return messageID
	}
	return GlobalManager.T(messageID, languages...)
}

// TWithData 全局翻译函数（带数据）
func TWithData(messageID string, data map[string]interface{}, languages ...string) string {
	if GlobalManager == nil {
		return messageID
	}
	return GlobalManager.TWithData(messageID, data, languages...)
}

// GetResponse 全局获取响应函数
func GetResponse(messageID string, languages ...string) (string, int) {
	if GlobalManager == nil {
		return messageID, 500
	}
	return GlobalManager.GetResponse(messageID, languages...)
}

// ParseAcceptLanguage 全局解析语言函数
func ParseAcceptLanguage(acceptLanguage string) []string {
	if GlobalManager == nil {
		return []string{"zh-CN"}
	}
	return GlobalManager.ParseAcceptLanguage(acceptLanguage)
}
