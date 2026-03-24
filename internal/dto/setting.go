package dto

import (
	"CBCTF/internal/config"
)

type UpdateSettingForm struct {
	Host *string `form:"host" json:"host" binding:"omitempty,url"`
	Path *string `form:"path" json:"path" binding:"omitempty,dir"`

	LogLevel *string `form:"log_level" json:"log_level" binding:"omitempty,oneof=DEBUG INFO WARNING ERROR"`
	LogSave  *bool   `form:"log_save" json:"log_save"`

	AsyncQLogLevel              *string `form:"asyncq_log_level" json:"asyncq_log_level" binding:"omitempty,oneof=DEBUG INFO WARNING ERROR"`
	AsyncQConcurrency           *int    `form:"asyncq_concurrency" json:"asyncq_concurrency" binding:"omitempty,gte=1"`
	AsyncQVictimConcurrency     *int    `form:"asyncq_victim_concurrency" json:"asyncq_victim_concurrency" binding:"omitempty,gte=1"`
	AsyncQGeneratorConcurrency  *int    `form:"asyncq_generator_concurrency" json:"asyncq_generator_concurrency" binding:"omitempty,gte=1"`
	AsyncQAttachmentConcurrency *int    `form:"asyncq_attachment_concurrency" json:"asyncq_attachment_concurrency" binding:"omitempty,gte=1"`
	AsyncQEmailConcurrency      *int    `form:"asyncq_email_concurrency" json:"asyncq_email_concurrency" binding:"omitempty,gte=1"`
	AsyncQWebhookConcurrency    *int    `form:"asyncq_webhook_concurrency" json:"asyncq_webhook_concurrency" binding:"omitempty,gte=1"`
	AsyncQImageConcurrency      *int    `form:"asyncq_image_concurrency" json:"asyncq_image_concurrency" binding:"omitempty,gte=1"`

	GinMode               *string   `form:"gin_mode" json:"gin_mode" binding:"omitempty,oneof=debug test release"`
	GinHost               *string   `form:"gin_host" json:"gin_host" binding:"omitempty,ip"`
	GinPort               *uint     `form:"gin_port" json:"gin_port" binding:"omitempty,port"`
	GinUploadMax          *int      `form:"gin_upload_max" json:"gin_upload_max" binding:"omitempty,gte=1"`
	GinProxies            *[]string `form:"gin_proxies" json:"gin_proxies" binding:"omitempty,dive,ip|cidr"`
	GinRateLimitGlobal    *int      `form:"gin_ratelimit_global" json:"gin_ratelimit_global" binding:"omitempty,gte=1"`
	GinRateLimitWhitelist *[]string `form:"gin_ratelimit_whitelist" json:"gin_ratelimit_whitelist" binding:"omitempty,dive,ip|cidr"`
	GinCORS               *[]string `form:"gin_cors" json:"gin_cors" binding:"omitempty,dive,url"`
	GinLogWhitelist       *[]string `form:"gin_log_whitelist" json:"gin_log_whitelist" binding:"omitempty,dive,uri"`
	GinJWTSecret          *string   `form:"gin_jwt_secret" json:"gin_jwt_secret" binding:"omitempty,min=11"`
	GinMetricsWhitelist   *[]string `form:"gin_metrics_whitelist" json:"gin_metrics_whitelist" binding:"omitempty,dive,ip|cidr"`

	GormPostgresHost    *string `form:"gorm_postgres_host" json:"gorm_postgres_host" binding:"omitempty,ip|hostname"`
	GormPostgresPort    *uint   `form:"gorm_postgres_port" json:"gorm_postgres_port" binding:"omitempty,port"`
	GormPostgresUser    *string `form:"gorm_postgres_user" json:"gorm_postgres_user" binding:"omitempty,min=1"`
	GormPostgresPwd     *string `form:"gorm_postgres_pwd" json:"gorm_postgres_pwd" binding:"omitempty,min=1,ascii"`
	GormPostgresDB      *string `form:"gorm_postgres_db" json:"gorm_postgres_db" binding:"omitempty,min=1"`
	GormPostgresSSLMode *bool   `form:"gorm_postgres_sslmode" json:"gorm_postgres_sslmode"`
	GormPostgresMXOpen  *int    `form:"gorm_postgres_mxopen" json:"gorm_postgres_mxopen" binding:"omitempty,gte=1"`
	GormPostgresMXIdle  *int    `form:"gorm_postgres_mxidle" json:"gorm_postgres_mxidle" binding:"omitempty,gte=1"`
	GormLogLevel        *string `form:"gorm_log_level" json:"gorm_log_level" binding:"omitempty,oneof=SILENT INFO WARNING ERROR"`

	RedisHost *string `form:"redis_host" json:"redis_host" binding:"omitempty,ip|hostname"`
	RedisPort *uint   `form:"redis_port" json:"redis_port" binding:"omitempty,port"`
	RedisPwd  *string `form:"redis_pwd" json:"redis_pwd" binding:"omitempty,min=1,ascii"`

	K8SConfig       *string `form:"k8s_config" json:"k8s_config" binding:"omitempty,filepath"`
	K8SNamespace    *string `form:"k8s_namespace" json:"k8s_namespace" binding:"omitempty,min=1,alphanum"`
	K8STCPDumpImage *string `form:"k8s_tcpdump" json:"k8s_tcpdump" binding:"omitempty,min=1"`

	K8SFrpOn         *bool                `form:"k8s_frp_on" json:"k8s_frp_on"`
	K8SFrpFrpcImage  *string              `form:"k8s_frp_frpc" json:"k8s_frp_frpc" binding:"omitempty,min=1"`
	K8SFrpNginxImage *string              `form:"k8s_frp_nginx" json:"k8s_frp_nginx" binding:"omitempty,min=1"`
	K8SFrpFrps       *[]config.FrpsConfig `form:"k8s_frp_frps" json:"k8s_frp_frps"`

	CheatIPWhitelist *[]string `form:"cheat_ip_whitelist" json:"cheat_ip_whitelist" binding:"omitempty,dive,ip|cidr"`

	WebhookWhitelist *[]string `form:"webhook_whitelist" json:"webhook_whitelist" binding:"omitempty,dive,ip|cidr|hostname|hostname_port"`

	RegistrationEnabled      *bool `form:"registration_enabled" json:"registration_enabled"`
	RegistrationDefaultGroup *uint `form:"registration_default_group" json:"registration_default_group" binding:"omitempty,gte=0"`

	GeoCityDB *string `form:"geocity_db" json:"geocity_db" binding:"omitempty,file"`
}
