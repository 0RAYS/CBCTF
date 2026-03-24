package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	HostSettingKey = "host"
	PathSettingKey = "path"

	LogLevelSettingKey = "log.level"
	LogSaveSettingKey  = "log.save"

	AsyncQLogLevelSettingKey       = "asynq.log.level"
	AsyncQConcurrencySettingKey    = "asynq.concurrency"
	AsyncQVictimConcurrencyKey     = "asynq.queues.victim"
	AsyncQGeneratorConcurrencyKey  = "asynq.queues.generator"
	AsyncQAttachmentConcurrencyKey = "asynq.queues.attachment"
	AsyncQEmailConcurrencyKey      = "asynq.queues.email"
	AsyncQWebhookConcurrencyKey    = "asynq.queues.webhook"
	AsyncQImageConcurrencyKey      = "asynq.queues.image"

	GinModeSettingKey               = "gin.mode"
	GinHostSettingKey               = "gin.host"
	GinPortSettingKey               = "gin.port"
	GinUploadMaxSettingKey          = "gin.upload.max"
	GinProxiesSettingKey            = "gin.proxies"
	GinRateLimitGlobalSettingKey    = "gin.ratelimit.global"
	GinRateLimitWhitelistSettingKey = "gin.ratelimit.whitelist"
	GinCORSSettingKey               = "gin.cors"
	GinLogWhitelistSettingKey       = "gin.log.whitelist"
	GinJWTSecretSettingKey          = "gin.jwt.secret"
	GinMetricsWhitelistSettingKey   = "gin.metrics.whitelist"

	GormPostgresHostSettingKey    = "gorm.postgres.host"
	GormPostgresPortSettingKey    = "gorm.postgres.port"
	GormPostgresUserSettingKey    = "gorm.postgres.user"
	GormPostgresPwdSettingKey     = "gorm.postgres.pwd"
	GormPostgresDBSettingKey      = "gorm.postgres.db"
	GormPostgresSSLModeSettingKey = "gorm.postgres.sslmode"
	GormPostgresMXOpenSettingKey  = "gorm.postgres.mxopen"
	GormPostgresMXIdleSettingKey  = "gorm.postgres.mxidle"
	GormLogLevelSettingKey        = "gorm.log.level"

	RedisHostSettingKey = "redis.host"
	RedisPortSettingKey = "redis.port"
	RedisPwdSettingKey  = "redis.pwd"

	K8SConfigSettingKey        = "k8s.config"
	K8SNamespaceSettingKey     = "k8s.namespace"
	K8STCPDumpImageSettingKey  = "k8s.tcpdump"
	K8SFrpOnSettingKey         = "k8s.frp.on"
	K8SFrpFrpcImageSettingKey  = "k8s.frp.frpc"
	K8SFrpNginxImageSettingKey = "k8s.frp.nginx"
	K8SFrpFrpsSettingKey       = "k8s.frp.frps"

	CheatIPWhitelistSettingKey = "cheat.ip.whitelist"

	WebhookWhitelistSettingKey = "webhook.whitelist"

	RegistrationEnabledSettingKey      = "registration.enabled"
	RegistrationDefaultGroupSettingKey = "registration.default_group"

	GeoCityDBSettingKey = "geocity_db"
)

type Setting struct {
	Key   string       `gorm:"size:255;uniqueIndex" json:"key"`
	Value SettingValue `gorm:"type:jsonb" json:"value"`
	BaseModel
}

type SettingValue struct {
	V any
}

func (s SettingValue) Value() (driver.Value, error) {
	if s.V == nil {
		return nil, nil
	}
	return json.Marshal(s.V)
}

func (s *SettingValue) Scan(value any) error {
	if value == nil {
		s.V = nil
		return nil
	}
	bs, err := scanBytes(value)
	if err != nil {
		return fmt.Errorf("failed to scan SettingValue value")
	}
	decoder := json.NewDecoder(bytes.NewReader(bs))
	decoder.UseNumber()
	var data any
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	s.V = data
	return nil
}
