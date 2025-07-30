package i18n

import (
	"fmt"
	"strings"
)

// MessageCategory 消息分类
type MessageCategory struct {
	Name     string
	Messages []string
}

// GetMessageCategories 获取所有消息分类
func GetMessageCategories() []MessageCategory {
	return []MessageCategory{
		{
			Name: "通用消息",
			Messages: []string{
				"Success", "UnsupportedKey", "DeadLock",
			},
		},
		{
			Name: "HTTP状态码",
			Messages: []string{
				"BadRequest", "Unauthorized", "Forbidden", "TooManyRequests", "UnknownError",
			},
		},
		{
			Name: "管理员相关",
			Messages: []string{
				"CreateAdminError", "DeleteAdminError", "GetAdminError", "AdminNotFound", "UpdateAdminError",
			},
		},
		{
			Name: "题目相关",
			Messages: []string{
				"CreateChallengeError", "DeleteChallengeError", "GetChallengeError", "ChallengeNotFound", "UpdateChallengeError", "InvalidChallengeType",
			},
		},
		{
			Name: "题目Flag相关",
			Messages: []string{
				"CreateChallengeFlagError", "DeleteChallengeFlagError", "GetChallengeFlagError", "ChallengeFlagNotFound", "UpdateChallengeFlagError",
			},
		},
		{
			Name: "作弊记录相关",
			Messages: []string{
				"CreateCheatError", "DeleteCheatError", "GetCheatError", "CheatNotFound", "UpdateCheatError",
			},
		},
		{
			Name: "容器相关",
			Messages: []string{
				"CreateContainerError", "DeleteContainerError", "GetContainerError", "ContainerNotFound", "UpdateContainerError",
			},
		},
		{
			Name: "比赛相关",
			Messages: []string{
				"CreateContestError", "DeleteContestError", "GetContestError", "ContestNotFound", "UpdateContestError",
				"DuplicateContestName", "ContestCaptchaError", "ContestIsComing", "ContestIsRunning", "ContestIsOver",
			},
		},
		{
			Name: "比赛题目相关",
			Messages: []string{
				"CreateContestChallengeError", "DeleteContestChallengeError", "GetContestChallengeError", "ContestChallengeNotFound", "UpdateContestChallengeError",
				"AlreadySolved", "FlagNotMatch", "NotAllowSubmit",
			},
		},
		{
			Name: "比赛Flag相关",
			Messages: []string{
				"CreateContestFlagError", "DeleteContestFlagError", "GetContestFlagError", "ContestFlagNotFound", "UpdateContestFlagError", "InvalidScoreType",
			},
		},
		{
			Name: "设备相关",
			Messages: []string{
				"CreateDeviceError", "DeleteDeviceError", "GetDeviceError", "DeviceNotFound", "UpdateDeviceError",
			},
		},
		{
			Name: "Docker相关",
			Messages: []string{
				"CreateDockerError", "DeleteDockerError", "GetDockerError", "DockerNotFound", "UpdateDockerError", "InvalidDockerImage",
			},
		},
		{
			Name: "事件相关",
			Messages: []string{
				"CreateEventError", "DeleteEventError", "GetEventError", "EventNotFound", "UpdateEventError",
			},
		},
		{
			Name: "文件相关",
			Messages: []string{
				"CreateFileError", "DeleteFileError", "GetFileError", "FileNotFound", "UpdateFileError", "FileNotAllowed",
			},
		},
		{
			Name: "通知相关",
			Messages: []string{
				"CreateNoticeError", "DeleteNoticeError", "GetNoticeError", "NoticeNotFound", "UpdateNoticeError", "InvalidNoticeType",
			},
		},
		{
			Name: "Pod相关",
			Messages: []string{
				"CreatePodError", "DeletePodError", "GetPodError", "PodNotFound", "UpdatePodError",
			},
		},
		{
			Name: "请求相关",
			Messages: []string{
				"CreateRequestError", "DeleteRequestError", "GetRequestError", "RequestNotFound", "UpdateRequestError",
			},
		},
		{
			Name: "提交相关",
			Messages: []string{
				"CreateSubmissionError", "DeleteSubmissionError", "GetSubmissionError", "SubmissionNotFound", "UpdateSubmissionError",
			},
		},
		{
			Name: "战队相关",
			Messages: []string{
				"CreateTeamError", "DeleteTeamError", "GetTeamError", "TeamNotFound", "UpdateTeamError",
				"TeamIsBanned", "TeamIsFull", "DuplicateTeamName", "DuplicateMember", "UserNotInTeam", "CaptainCannotLeave", "TeamCaptchaError",
			},
		},
		{
			Name: "战队Flag相关",
			Messages: []string{
				"CreateTeamFlagError", "DeleteTeamFlagError", "GetTeamFlagError", "TeamFlagNotFound", "UpdateTeamFlagError",
			},
		},
		{
			Name: "流量相关",
			Messages: []string{
				"CreateTrafficError", "DeleteTrafficError", "GetTrafficError", "TrafficNotFound", "UpdateTrafficError",
				"ReadPcapError", "PcapNotFound", "HasNoTraffic",
			},
		},
		{
			Name: "用户相关",
			Messages: []string{
				"CreateUserError", "DeleteUserError", "GetUserError", "UserNotFound", "UpdateUserError",
				"InvalidEmail", "UnverifiedEmail", "DuplicateEmail", "DuplicateUserName", "WeakPassword",
				"NameOrPasswordError", "PasswordSame", "PasswordError",
			},
		},
		{
			Name: "靶机相关",
			Messages: []string{
				"CreateVictimError", "DeleteVictimError", "GetVictimError", "VictimNotFound", "UpdateVictimError", "HasMuchTime",
			},
		},
		{
			Name: "用户关联操作",
			Messages: []string{
				"AppendUserToTeamError", "AppendUserToContestError", "DeleteUserFromTeamError", "DeleteUserFromContestError",
			},
		},
		{
			Name: "排名相关",
			Messages: []string{
				"UpdateRankingError",
			},
		},
		{
			Name: "邮箱验证相关",
			Messages: []string{
				"SetEmailVerifyTokenError", "GetEmailVerifyTokenError", "DelEmailVerifyTokenError", "InvalidEmailVerifyToken", "SendEmailError", "RedisError",
			},
		},
		{
			Name: "文件操作相关",
			Messages: []string{
				"CreateDirError", "ReadDirError", "InvalidFileName", "InvalidDockerComposeYaml", "InvalidChallengeFlagInjectType", "CopyFileError", "ExecCommandError", "ZipError",
			},
		},
		{
			Name: "Kubernetes相关",
			Messages: []string{
				"CreateNamespaceError", "DeleteNamespaceError", "GetNamespaceError", "NamespaceNotFound",
				"CreateConfigMapError", "DeleteConfigMapError", "GetConfigMapError", "ConfigMapNotFound",
				"CreateNetworkPolicyError", "DeleteNetworkPolicyError", "GetNetworkPolicyError", "NetworkPolicyNotFound",
				"CreateServiceError", "DeleteServiceError", "GetServiceError", "ServiceNotFound",
				"CreateJobError", "DeleteJobError", "GetJobError", "JobNotFound",
				"CreateVPCError", "DeleteVPCError", "GetVPCError", "VPCNotFound",
				"CreateSubnetError", "DeleteSubnetError", "GetSubnetError", "SubnetNotFound",
				"CreateVPCNatGatewayError", "DeleteVPCNatGatewayError", "GetVPCNatGatewayError", "VPCNatGatewayNotFound",
				"CreateEIPError", "DeleteEIPError", "GetEIPError", "EIPNotFound",
				"CreateDNatError", "DeleteDNatError", "GetDNatError", "DNatNotFound",
				"CreateSNatError", "DeleteSNatError", "GetSNatError", "SNatNotFound",
				"CreateNetAttError", "DeleteNetAttError", "GetNetAttError", "NetAttNotFound",
				"CreateIPError", "DeleteIPError", "GetIPError", "IPNotFound",
				"CreatePVError", "DeletePVError", "GetPVError", "PVNotFound",
				"CreatePVCError", "DeletePVCError", "GetPVCError", "PVCNotFound",
				"GetNodeListError",
			},
		},
	}
}

// ValidateMessageID 验证消息ID是否有效
func ValidateMessageID(messageID string) bool {
	categories := GetMessageCategories()
	for _, category := range categories {
		for _, msg := range category.Messages {
			if msg == messageID {
				return true
			}
		}
	}
	return false
}

// GetMessageCategory 获取消息所属的分类
func GetMessageCategory(messageID string) string {
	categories := GetMessageCategories()
	for _, category := range categories {
		for _, msg := range category.Messages {
			if msg == messageID {
				return category.Name
			}
		}
	}
	return "未知分类"
}

// GetAllMessageIDs 获取所有消息ID
func GetAllMessageIDs() []string {
	var allMessages []string
	categories := GetMessageCategories()
	for _, category := range categories {
		allMessages = append(allMessages, category.Messages...)
	}
	return allMessages
}

// FormatMessageID 格式化消息ID为可读的格式
func FormatMessageID(messageID string) string {
	// 将驼峰命名转换为空格分隔的格式
	var result strings.Builder
	for i, char := range messageID {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteByte(' ')
		}
		result.WriteRune(char)
	}
	return result.String()
}

// GetStatusCodeDescription 获取状态码描述
func GetStatusCodeDescription(code int) string {
	descriptions := map[int]string{
		200: "成功",
		400: "请求错误",
		401: "未授权",
		403: "禁止访问",
		404: "未找到",
		429: "请求过于频繁",
		500: "服务器内部错误",
	}
	
	if desc, exists := descriptions[code]; exists {
		return desc
	}
	return fmt.Sprintf("未知状态码: %d", code)
}

// GetLanguageDisplayName 获取语言显示名称
func GetLanguageDisplayName(lang string) string {
	names := map[string]string{
		"zh-CN": "中文（简体）",
		"en-US": "English",
		"origin": "原始消息",
	}
	
	if name, exists := names[lang]; exists {
		return name
	}
	return lang
} 