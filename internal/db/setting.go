package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"encoding/json"
	"reflect"

	"gorm.io/gorm"
)

type SettingRepo struct {
	BaseRepo[model.Setting]
}

type CreateSettingOptions struct {
	Key   string
	Value model.SettingValue
}

func (c CreateSettingOptions) Convert2Model() model.Model {
	return model.Setting{
		Key:   c.Key,
		Value: c.Value,
	}
}

type UpdateSettingOptions struct {
	Value *model.SettingValue
}

func (u UpdateSettingOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Value != nil {
		options["value"] = *u.Value
	}
	return options
}

func InitSettingRepo(tx *gorm.DB) *SettingRepo {
	return &SettingRepo{
		BaseRepo: BaseRepo[model.Setting]{
			DB: tx,
		},
	}
}

func (s *SettingRepo) Get(key string, optionsL ...GetOptions) (model.Setting, model.RetVal) {
	return s.GetByUniqueField("key", key, optionsL...)
}

func (s *SettingRepo) Update(key string, options UpdateSettingOptions) model.RetVal {
	var count uint
	data := options.Convert2Map()
	if value, ok := data["value"]; !ok || value == nil || reflect.ValueOf(value.(model.SettingValue).V).IsNil() {
		return model.SuccessRetVal()
	}
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Setting: too many times failed due to optimistic lock")
			return model.RetVal{Msg: i18n.Model.Setting.DeadLock}
		}
		m, ret := s.Get(key)
		if !ret.OK {
			return ret
		}
		res := s.DB.Model(&m).Where("id = ?", m.ID).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Setting: %s", res.Error)
			return model.RetVal{Msg: i18n.Model.Setting.UpdateError, Attr: map[string]any{"Error": res.Error.Error()}}
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return model.SuccessRetVal()
}

func (s *SettingRepo) InitSettings() model.RetVal {
	for _, setting := range []model.Setting{
		{Key: model.HostSettingKey, Value: model.SettingValue{V: config.Env.Host}},
		{Key: model.PathSettingKey, Value: model.SettingValue{V: config.Env.Path}},

		{Key: model.LogLevelSettingKey, Value: model.SettingValue{V: config.Env.Log.Level}},
		{Key: model.LogSaveSettingKey, Value: model.SettingValue{V: config.Env.Log.Save}},

		{Key: model.AsyncQLogLevelSettingKey, Value: model.SettingValue{V: config.Env.AsyncQ.Log.Level}},
		{Key: model.AsyncQConcurrencySettingKey, Value: model.SettingValue{V: config.Env.AsyncQ.Concurrency}},
		{Key: model.AsyncQVictimConcurrencyKey, Value: model.SettingValue{V: config.Env.AsyncQ.Queues.Victim}},
		{Key: model.AsyncQGeneratorConcurrencyKey, Value: model.SettingValue{V: config.Env.AsyncQ.Queues.Generator}},
		{Key: model.AsyncQAttachmentConcurrencyKey, Value: model.SettingValue{V: config.Env.AsyncQ.Queues.Attachment}},
		{Key: model.AsyncQEmailConcurrencyKey, Value: model.SettingValue{V: config.Env.AsyncQ.Queues.Email}},
		{Key: model.AsyncQWebhookConcurrencyKey, Value: model.SettingValue{V: config.Env.AsyncQ.Queues.Webhook}},
		{Key: model.AsyncQImageConcurrencyKey, Value: model.SettingValue{V: config.Env.AsyncQ.Queues.Image}},

		{Key: model.GinModeSettingKey, Value: model.SettingValue{V: config.Env.Gin.Mode}},
		{Key: model.GinHostSettingKey, Value: model.SettingValue{V: config.Env.Gin.Host}},
		{Key: model.GinPortSettingKey, Value: model.SettingValue{V: config.Env.Gin.Port}},
		{Key: model.GinUploadMaxSettingKey, Value: model.SettingValue{V: config.Env.Gin.Upload.Max}},
		{Key: model.GinProxiesSettingKey, Value: model.SettingValue{V: config.Env.Gin.Proxies}},
		{Key: model.GinRateLimitGlobalSettingKey, Value: model.SettingValue{V: config.Env.Gin.RateLimit.Global}},
		{Key: model.GinRateLimitWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Gin.RateLimit.Whitelist}},
		{Key: model.GinCORSSettingKey, Value: model.SettingValue{V: config.Env.Gin.CORS}},
		{Key: model.GinLogWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Gin.Log.Whitelist}},
		{Key: model.GinJWTSecretSettingKey, Value: model.SettingValue{V: config.Env.Gin.JWT.Secret}},
		{Key: model.GinMetricsWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Gin.Metrics.Whitelist}},

		{Key: model.GormPostgresHostSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Postgres.Host}},
		{Key: model.GormPostgresPortSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Postgres.Port}},
		{Key: model.GormPostgresUserSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Postgres.User}},
		{Key: model.GormPostgresPwdSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Postgres.Pwd}},
		{Key: model.GormPostgresDBSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Postgres.DB}},
		{Key: model.GormPostgresSSLModeSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Postgres.SSLMode}},
		{Key: model.GormPostgresMXOpenSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Postgres.MaxOpenConns}},
		{Key: model.GormPostgresMXIdleSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Postgres.MaxIdleConns}},
		{Key: model.GormLogLevelSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Log.Level}},

		{Key: model.RedisHostSettingKey, Value: model.SettingValue{V: config.Env.Redis.Host}},
		{Key: model.RedisPortSettingKey, Value: model.SettingValue{V: config.Env.Redis.Port}},
		{Key: model.RedisPwdSettingKey, Value: model.SettingValue{V: config.Env.Redis.Pwd}},

		{Key: model.K8SConfigSettingKey, Value: model.SettingValue{V: config.Env.K8S.Config}},
		{Key: model.K8SNamespaceSettingKey, Value: model.SettingValue{V: config.Env.K8S.Namespace}},
		{Key: model.K8STCPDumpImageSettingKey, Value: model.SettingValue{V: config.Env.K8S.TCPDumpImage}},
		{Key: model.K8SFrpOnSettingKey, Value: model.SettingValue{V: config.Env.K8S.Frp.On}},
		{Key: model.K8SFrpFrpcImageSettingKey, Value: model.SettingValue{V: config.Env.K8S.Frp.FrpcImage}},
		{Key: model.K8SFrpNginxImageSettingKey, Value: model.SettingValue{V: config.Env.K8S.Frp.NginxImage}},
		{Key: model.K8SFrpFrpsSettingKey, Value: model.SettingValue{V: config.Env.K8S.Frp.Frps}},

		{Key: model.CheatIPWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Cheat.IP.Whitelist}},

		{Key: model.WebhookWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Webhook.Whitelist}},

		{Key: model.RegistrationEnabledSettingKey, Value: model.SettingValue{V: config.Env.Registration.Enabled}},
		{Key: model.RegistrationDefaultGroupSettingKey, Value: model.SettingValue{V: config.Env.Registration.DefaultGroup}},

		{Key: model.GeoCityDBSettingKey, Value: model.SettingValue{V: config.Env.GeoCityDB}},
	} {
		if _, ret := s.Create(CreateSettingOptions{
			Key:   setting.Key,
			Value: setting.Value,
		}); !ret.OK && ret.Msg != i18n.Model.DuplicateKeyValue {
			return ret
		}
	}
	return s.ReadSettings()
}

func (s *SettingRepo) ReadSettings() model.RetVal {
	var ret model.RetVal

	if config.Env.Host, ret = GetValue[string](s, model.HostSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Path, ret = GetValue[string](s, model.PathSettingKey); !ret.OK {
		return ret
	}

	if config.Env.Log.Level, ret = GetValue[string](s, model.LogLevelSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Log.Save, ret = GetValue[bool](s, model.LogSaveSettingKey); !ret.OK {
		return ret
	}

	if config.Env.AsyncQ.Log.Level, ret = GetValue[string](s, model.AsyncQLogLevelSettingKey); !ret.OK {
		return ret
	}
	if config.Env.AsyncQ.Concurrency, ret = GetValue[int](s, model.AsyncQConcurrencySettingKey); !ret.OK {
		return ret
	}
	if config.Env.AsyncQ.Queues.Victim, ret = GetValue[int](s, model.AsyncQVictimConcurrencyKey); !ret.OK {
		return ret
	}
	if config.Env.AsyncQ.Queues.Generator, ret = GetValue[int](s, model.AsyncQGeneratorConcurrencyKey); !ret.OK {
		return ret
	}
	if config.Env.AsyncQ.Queues.Attachment, ret = GetValue[int](s, model.AsyncQAttachmentConcurrencyKey); !ret.OK {
		return ret
	}
	if config.Env.AsyncQ.Queues.Email, ret = GetValue[int](s, model.AsyncQEmailConcurrencyKey); !ret.OK {
		return ret
	}
	if config.Env.AsyncQ.Queues.Webhook, ret = GetValue[int](s, model.AsyncQWebhookConcurrencyKey); !ret.OK {
		return ret
	}
	if config.Env.AsyncQ.Queues.Image, ret = GetValue[int](s, model.AsyncQImageConcurrencyKey); !ret.OK {
		return ret
	}

	if config.Env.Gin.Mode, ret = GetValue[string](s, model.GinModeSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.Host, ret = GetValue[string](s, model.GinHostSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.Port, ret = GetValue[uint](s, model.GinPortSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.Upload.Max, ret = GetValue[int](s, model.GinUploadMaxSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.Proxies, ret = GetValue[[]string](s, model.GinProxiesSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.RateLimit.Global, ret = GetValue[int](s, model.GinRateLimitGlobalSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.RateLimit.Whitelist, ret = GetValue[[]string](s, model.GinRateLimitWhitelistSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.CORS, ret = GetValue[[]string](s, model.GinCORSSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.Log.Whitelist, ret = GetValue[[]string](s, model.GinLogWhitelistSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.JWT.Secret, ret = GetValue[string](s, model.GinJWTSecretSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gin.Metrics.Whitelist, ret = GetValue[[]string](s, model.GinMetricsWhitelistSettingKey); !ret.OK {
		return ret
	}

	if config.Env.Gorm.Postgres.Host, ret = GetValue[string](s, model.GormPostgresHostSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.Postgres.Port, ret = GetValue[uint](s, model.GormPostgresPortSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.Postgres.User, ret = GetValue[string](s, model.GormPostgresUserSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.Postgres.Pwd, ret = GetValue[string](s, model.GormPostgresPwdSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.Postgres.DB, ret = GetValue[string](s, model.GormPostgresDBSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.Postgres.SSLMode, ret = GetValue[bool](s, model.GormPostgresSSLModeSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.Postgres.MaxOpenConns, ret = GetValue[int](s, model.GormPostgresMXOpenSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.Postgres.MaxIdleConns, ret = GetValue[int](s, model.GormPostgresMXIdleSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.Log.Level, ret = GetValue[string](s, model.GormLogLevelSettingKey); !ret.OK {
		return ret
	}

	if config.Env.Redis.Host, ret = GetValue[string](s, model.RedisHostSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Redis.Port, ret = GetValue[uint](s, model.RedisPortSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Redis.Pwd, ret = GetValue[string](s, model.RedisPwdSettingKey); !ret.OK {
		return ret
	}

	if config.Env.K8S.Config, ret = GetValue[string](s, model.K8SConfigSettingKey); !ret.OK {
		return ret
	}
	if config.Env.K8S.Namespace, ret = GetValue[string](s, model.K8SNamespaceSettingKey); !ret.OK {
		return ret
	}
	if config.Env.K8S.TCPDumpImage, ret = GetValue[string](s, model.K8STCPDumpImageSettingKey); !ret.OK {
		return ret
	}
	if config.Env.K8S.Frp.On, ret = GetValue[bool](s, model.K8SFrpOnSettingKey); !ret.OK {
		return ret
	}
	if config.Env.K8S.Frp.FrpcImage, ret = GetValue[string](s, model.K8SFrpFrpcImageSettingKey); !ret.OK {
		return ret
	}
	if config.Env.K8S.Frp.NginxImage, ret = GetValue[string](s, model.K8SFrpNginxImageSettingKey); !ret.OK {
		return ret
	}
	if config.Env.K8S.Frp.Frps, ret = GetValue[[]config.FrpsConfig](s, model.K8SFrpFrpsSettingKey); !ret.OK {
		return ret
	}

	if config.Env.Cheat.IP.Whitelist, ret = GetValue[[]string](s, model.CheatIPWhitelistSettingKey); !ret.OK {
		return ret
	}

	if config.Env.Webhook.Whitelist, ret = GetValue[[]string](s, model.WebhookWhitelistSettingKey); !ret.OK {
		return ret
	}

	if config.Env.Registration.Enabled, ret = GetValue[bool](s, model.RegistrationEnabledSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Registration.DefaultGroup, ret = GetValue[uint](s, model.RegistrationDefaultGroupSettingKey); !ret.OK {
		return ret
	}

	if config.Env.GeoCityDB, ret = GetValue[string](s, model.GeoCityDBSettingKey); !ret.OK {
		return ret
	}
	config.Tidy()
	if err := config.Save(); err != nil {
		log.Logger.Warningf("Failed to save config: %s, but it's not important, all config will be read from database", err.Error())
	}

	return model.SuccessRetVal()
}

func GetValue[T any](s *SettingRepo, key string) (T, model.RetVal) {
	data, ret := s.Get(key)
	if !ret.OK {
		var zero T
		return zero, ret
	}

	var zero T
	targetType := reflect.TypeOf(zero)

	var out T
	bytes, err := json.Marshal(data.Value.V)
	if err != nil {
		return zero, model.RetVal{
			Msg: i18n.Model.Setting.InvalidType,
			Attr: map[string]any{
				"Key":         data.Key,
				"Type":        targetType.String(),
				"InvalidType": reflect.TypeOf(data.Value.V).String(),
			},
		}
	}
	if err = json.Unmarshal(bytes, &out); err != nil {
		return zero, model.RetVal{
			Msg: i18n.Model.Setting.InvalidType,
			Attr: map[string]any{
				"Key":         data.Key,
				"Type":        targetType.String(),
				"InvalidType": reflect.TypeOf(data.Value.V).String(),
			},
		}
	}

	return out, model.SuccessRetVal()
}
