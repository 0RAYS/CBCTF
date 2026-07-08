package dto

import "CBCTF/internal/config"

type UpdateSettingForm struct {
	Host *string `form:"host" json:"host" binding:"omitempty,url"`

	AsyncQLogLevel              *string `form:"asyncq_log_level" json:"asyncq_log_level" binding:"omitempty,oneof=DEBUG INFO WARNING ERROR debug info warning error"`
	AsyncQVictimConcurrency     *int    `form:"asyncq_victim_concurrency" json:"asyncq_victim_concurrency" binding:"omitempty,gte=1"`
	AsyncQTrafficConcurrency    *int    `form:"asyncq_traffic_concurrency" json:"asyncq_traffic_concurrency" binding:"omitempty,gte=1"`
	AsyncQGeneratorConcurrency  *int    `form:"asyncq_generator_concurrency" json:"asyncq_generator_concurrency" binding:"omitempty,gte=1"`
	AsyncQAttachmentConcurrency *int    `form:"asyncq_attachment_concurrency" json:"asyncq_attachment_concurrency" binding:"omitempty,gte=1"`
	AsyncQEmailConcurrency      *int    `form:"asyncq_email_concurrency" json:"asyncq_email_concurrency" binding:"omitempty,gte=1"`
	AsyncQWebhookConcurrency    *int    `form:"asyncq_webhook_concurrency" json:"asyncq_webhook_concurrency" binding:"omitempty,gte=1"`
	AsyncQImageConcurrency      *int    `form:"asyncq_image_concurrency" json:"asyncq_image_concurrency" binding:"omitempty,gte=1"`

	GinMode               *string   `form:"gin_mode" json:"gin_mode" binding:"omitempty,oneof=DEBUG TEST RELEASE debug test release"`
	GinUploadPicture      *int      `form:"gin_upload_picture" json:"gin_upload_picture" binding:"omitempty,gte=1"`
	GinUploadChallenge    *int      `form:"gin_upload_challenge" json:"gin_upload_challenge" binding:"omitempty,gte=1"`
	GinUploadWriteup      *int      `form:"gin_upload_writeup" json:"gin_upload_writeup" binding:"omitempty,gte=1"`
	GinProxies            *[]string `form:"gin_proxies" json:"gin_proxies" binding:"omitempty,dive,ip|cidr"`
	GinRateLimitGlobal    *int      `form:"gin_ratelimit_global" json:"gin_ratelimit_global" binding:"omitempty,gte=1"`
	GinRateLimitWhitelist *[]string `form:"gin_ratelimit_whitelist" json:"gin_ratelimit_whitelist" binding:"omitempty,dive,ip|cidr"`
	GinOrigins            *[]string `form:"gin_origins" json:"gin_origins" binding:"omitempty,dive,url"`
	GinLogWhitelist       *[]string `form:"gin_log_whitelist" json:"gin_log_whitelist" binding:"omitempty,dive,uri"`
	GinJWTSecret          *string   `form:"gin_jwt_secret" json:"gin_jwt_secret" binding:"omitempty,min=11"`
	GinMetricsWhitelist   *[]string `form:"gin_metrics_whitelist" json:"gin_metrics_whitelist" binding:"omitempty,dive,ip|cidr"`
	GinPProfWhitelist     *[]string `form:"gin_pprof_whitelist" json:"gin_pprof_whitelist" binding:"omitempty,dive,ip|cidr"`

	K8SNamespace     *string              `form:"k8s_namespace" json:"k8s_namespace" binding:"omitempty,min=1,alphanum"`
	K8SCaptureImage  *string              `form:"k8s_capture" json:"k8s_capture" binding:"omitempty,min=1"`
	K8SFrpOn         *bool                `form:"k8s_frp_on" json:"k8s_frp_on"`
	K8SFrpFrpcImage  *string              `form:"k8s_frp_frpc" json:"k8s_frp_frpc" binding:"omitempty,min=1"`
	K8SFrpNginxImage *string              `form:"k8s_frp_nginx" json:"k8s_frp_nginx" binding:"omitempty,min=1"`
	K8SFrpFrps       *[]config.FrpsConfig `form:"k8s_frp_frps" json:"k8s_frp_frps" binding:"omitempty,dive"`

	CheatIPWhitelist         *[]string `form:"cheat_ip_whitelist" json:"cheat_ip_whitelist" binding:"omitempty,dive,ip|cidr"`
	WebhookWhitelist         *[]string `form:"webhook_whitelist" json:"webhook_whitelist" binding:"omitempty,dive,ip|cidr|hostname|hostname_port"`
	RegistrationEnabled      *bool     `form:"registration_enabled" json:"registration_enabled"`
	RegistrationDefaultGroup *uint     `form:"registration_default_group" json:"registration_default_group" binding:"omitempty,gte=0"`
}
