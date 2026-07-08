package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	HostSettingKey = "host"

	AsyncQLogLevelSettingKey       = "asynq.log.level"
	AsyncQVictimConcurrencyKey     = "asynq.queues.victim"
	AsyncQTrafficConcurrencyKey    = "asynq.queues.traffic"
	AsyncQGeneratorConcurrencyKey  = "asynq.queues.generator"
	AsyncQAttachmentConcurrencyKey = "asynq.queues.attachment"
	AsyncQEmailConcurrencyKey      = "asynq.queues.email"
	AsyncQWebhookConcurrencyKey    = "asynq.queues.webhook"
	AsyncQImageConcurrencyKey      = "asynq.queues.image"

	GinModeSettingKey               = "gin.mode"
	GinUploadPictureSettingKey      = "gin.upload.picture"
	GinUploadChallengeSettingKey    = "gin.upload.challenge"
	GinUploadWriteupSettingKey      = "gin.upload.writeup"
	GinProxiesSettingKey            = "gin.proxies"
	GinRateLimitGlobalSettingKey    = "gin.ratelimit.global"
	GinRateLimitWhitelistSettingKey = "gin.ratelimit.whitelist"
	GinOriginsSettingKey            = "gin.origins"
	GinLogWhitelistSettingKey       = "gin.log.whitelist"
	GinJWTSecretSettingKey          = "gin.jwt.secret"
	GinMetricsWhitelistSettingKey   = "gin.metrics.whitelist"

	K8SNamespaceSettingKey     = "k8s.namespace"
	K8SCaptureImageSettingKey  = "k8s.capture"
	K8SFrpOnSettingKey         = "k8s.frp.on"
	K8SFrpFrpcImageSettingKey  = "k8s.frp.frpc"
	K8SFrpNginxImageSettingKey = "k8s.frp.nginx"
	K8SFrpFrpsSettingKey       = "k8s.frp.frps"

	CheatIPWhitelistSettingKey         = "cheat.ip.whitelist"
	WebhookWhitelistSettingKey         = "webhook.whitelist"
	RegistrationEnabledSettingKey      = "registration.enabled"
	RegistrationDefaultGroupSettingKey = "registration.default_group"
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
	if err = decoder.Decode(&data); err != nil {
		return err
	}
	s.V = data
	return nil
}
