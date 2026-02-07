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

	AsyncQLogLevelSettingKey    = "asynq.log.level"
	AsyncQConcurrencySettingKey = "asynq.concurrency"

	GinModeSettingKey               = "gin.mode"
	GinHostSettingKey               = "gin.host"
	GinPortSettingKey               = "gin.port"
	GinUploadMaxSettingKey          = "gin.upload.max"
	GinProxiesSettingKey            = "gin.proxies"
	GinRateLimitGlobalSettingKey    = "gin.ratelimit.global"
	GinRateLimitWhitelistSettingKey = "gin.ratelimit.whitelist"
	GinCORSSettingKey               = "gin.cors"
	GinLogWhitelistSettingKey       = "gin.log.whitelist"

	GormMySQLHostSettingKey   = "gorm.mysql.host"
	GormMySQLPortSettingKey   = "gorm.mysql.port"
	GormMySQLUserSettingKey   = "gorm.mysql.user"
	GormMySQLPwdSettingKey    = "gorm.mysql.pwd"
	GormMySQLDBSettingKey     = "gorm.mysql.db"
	GormMySQLMXOpenSettingKey = "gorm.mysql.mxopen"
	GormMySQLMXIdleSettingKey = "gorm.mysql.mxidle"
	GormLogLevelSettingKey    = "gorm.log.level"

	RedisHostSettingKey = "redis.host"
	RedisPortSettingKey = "redis.port"
	RedisPwdSettingKey  = "redis.pwd"

	K8SConfigSettingKey                    = "k8s.config"
	K8SNamespaceSettingKey                 = "k8s.namespace"
	K8SExternalNetworkCIDRSettingKey       = "k8s.external_network.cidr"
	K8SExternalNetworkGatewaySettingKey    = "k8s.external_network.gateway"
	K8SExternalNetworkInterfaceSettingKey  = "k8s.external_network.interface"
	K8SExternalNetworkExcludeIPsSettingKey = "k8s.external_network.exclude_ips"
	K8STCPDumpImageSettingKey              = "k8s.tcpdump"
	K8SFrpOnSettingKey                     = "k8s.frp.on"
	K8SFrpFrpcImageSettingKey              = "k8s.frp.frpc"
	K8SFrpNginxImageSettingKey             = "k8s.frp.nginx"
	K8SFrpFrpsSettingKey                   = "k8s.frp.frps"
	K8SGeneratorWorkerSettingKey           = "k8s.generator_worker"

	NFSServerSettingKey  = "nfs.server"
	NFSPathSettingKey    = "nfs.path"
	NFSStorageSettingKey = "nfs.storage"

	CheatIPWhitelistSettingKey = "cheat.ip.whitelist"
)

type Setting struct {
	Key   string       `gorm:"size:255;uniqueIndex" json:"key"`
	Value SettingValue `gorm:"type:json" json:"value"`
	BaseModel
}

func (s Setting) ModelName() string {
	return "Setting"
}

func (s Setting) GetBaseModel() BaseModel {
	return s.BaseModel
}

func (s Setting) UniqueFields() []string {
	return []string{"key"}
}

func (s Setting) QueryFields() []string {
	return []string{"id", "key", "value"}
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
	bs, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SettingValue value")
	}
	if len(bs) == 0 {
		s.V = nil
		return nil
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
