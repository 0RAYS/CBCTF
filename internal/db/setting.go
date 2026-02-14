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
	return s.GetByUniqueKey("key", key, optionsL...)
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
			return model.RetVal{Msg: i18n.Model.DeadLock, Attr: map[string]any{"Model": model.Setting{}.ModelName()}}
		}
		m, ret := s.Get(key, GetOptions{Selects: []string{"id", "version"}})
		if !ret.OK {
			return ret
		}
		res := s.DB.Model(&m).Where("id = ?", m.ID).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Setting: %s", res.Error)
			return model.RetVal{Msg: i18n.Model.UpdateError, Attr: map[string]any{"Model": m.ModelName(), "Error": res.Error.Error()}}
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

		{Key: model.GinModeSettingKey, Value: model.SettingValue{V: config.Env.Gin.Mode}},
		{Key: model.GinHostSettingKey, Value: model.SettingValue{V: config.Env.Gin.Host}},
		{Key: model.GinPortSettingKey, Value: model.SettingValue{V: config.Env.Gin.Port}},
		{Key: model.GinUploadMaxSettingKey, Value: model.SettingValue{V: config.Env.Gin.Upload.Max}},
		{Key: model.GinProxiesSettingKey, Value: model.SettingValue{V: config.Env.Gin.Proxies}},
		{Key: model.GinRateLimitGlobalSettingKey, Value: model.SettingValue{V: config.Env.Gin.RateLimit.Global}},
		{Key: model.GinRateLimitWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Gin.RateLimit.Whitelist}},
		{Key: model.GinCORSSettingKey, Value: model.SettingValue{V: config.Env.Gin.CORS}},
		{Key: model.GinLogWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Gin.Log.Whitelist}},

		{Key: model.GormMySQLHostSettingKey, Value: model.SettingValue{V: config.Env.Gorm.MySQL.Host}},
		{Key: model.GormMySQLPortSettingKey, Value: model.SettingValue{V: config.Env.Gorm.MySQL.Port}},
		{Key: model.GormMySQLUserSettingKey, Value: model.SettingValue{V: config.Env.Gorm.MySQL.User}},
		{Key: model.GormMySQLPwdSettingKey, Value: model.SettingValue{V: config.Env.Gorm.MySQL.Pwd}},
		{Key: model.GormMySQLDBSettingKey, Value: model.SettingValue{V: config.Env.Gorm.MySQL.DB}},
		{Key: model.GormMySQLMXOpenSettingKey, Value: model.SettingValue{V: config.Env.Gorm.MySQL.MaxOpenConns}},
		{Key: model.GormMySQLMXIdleSettingKey, Value: model.SettingValue{V: config.Env.Gorm.MySQL.MaxIdleConns}},
		{Key: model.GormLogLevelSettingKey, Value: model.SettingValue{V: config.Env.Gorm.Log.Level}},

		{Key: model.RedisHostSettingKey, Value: model.SettingValue{V: config.Env.Redis.Host}},
		{Key: model.RedisPortSettingKey, Value: model.SettingValue{V: config.Env.Redis.Port}},
		{Key: model.RedisPwdSettingKey, Value: model.SettingValue{V: config.Env.Redis.Pwd}},

		{Key: model.K8SConfigSettingKey, Value: model.SettingValue{V: config.Env.K8S.Config}},
		{Key: model.K8SNamespaceSettingKey, Value: model.SettingValue{V: config.Env.K8S.Namespace}},
		{Key: model.K8SExternalNetworkCIDRSettingKey, Value: model.SettingValue{V: config.Env.K8S.ExternalNetwork.CIDR}},
		{Key: model.K8SExternalNetworkGatewaySettingKey, Value: model.SettingValue{V: config.Env.K8S.ExternalNetwork.Gateway}},
		{Key: model.K8SExternalNetworkInterfaceSettingKey, Value: model.SettingValue{V: config.Env.K8S.ExternalNetwork.Interface}},
		{Key: model.K8SExternalNetworkExcludeIPsSettingKey, Value: model.SettingValue{V: config.Env.K8S.ExternalNetwork.ExcludeIPs}},
		{Key: model.K8STCPDumpImageSettingKey, Value: model.SettingValue{V: config.Env.K8S.TCPDumpImage}},
		{Key: model.K8SFrpOnSettingKey, Value: model.SettingValue{V: config.Env.K8S.Frp.On}},
		{Key: model.K8SFrpFrpcImageSettingKey, Value: model.SettingValue{V: config.Env.K8S.Frp.FrpcImage}},
		{Key: model.K8SFrpNginxImageSettingKey, Value: model.SettingValue{V: config.Env.K8S.Frp.NginxImage}},
		{Key: model.K8SFrpFrpsSettingKey, Value: model.SettingValue{V: config.Env.K8S.Frp.Frps}},
		{Key: model.K8SGeneratorWorkerSettingKey, Value: model.SettingValue{V: config.Env.K8S.GeneratorWorker}},

		{Key: model.NFSServerSettingKey, Value: model.SettingValue{V: config.Env.NFS.Server}},
		{Key: model.NFSPathSettingKey, Value: model.SettingValue{V: config.Env.NFS.Path}},
		{Key: model.NFSStorageSettingKey, Value: model.SettingValue{V: config.Env.NFS.Storage}},

		{Key: model.CheatIPWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Cheat.IP.Whitelist}},

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

	if config.Env.Gorm.MySQL.Host, ret = GetValue[string](s, model.GormMySQLHostSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.MySQL.Port, ret = GetValue[uint](s, model.GormMySQLPortSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.MySQL.User, ret = GetValue[string](s, model.GormMySQLUserSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.MySQL.Pwd, ret = GetValue[string](s, model.GormMySQLPwdSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.MySQL.DB, ret = GetValue[string](s, model.GormMySQLDBSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.MySQL.MaxOpenConns, ret = GetValue[int](s, model.GormMySQLMXOpenSettingKey); !ret.OK {
		return ret
	}
	if config.Env.Gorm.MySQL.MaxIdleConns, ret = GetValue[int](s, model.GormMySQLMXIdleSettingKey); !ret.OK {
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
	if config.Env.K8S.ExternalNetwork.CIDR, ret = GetValue[string](s, model.K8SExternalNetworkCIDRSettingKey); !ret.OK {
		return ret
	}
	if config.Env.K8S.ExternalNetwork.Gateway, ret = GetValue[string](s, model.K8SExternalNetworkGatewaySettingKey); !ret.OK {
		return ret
	}
	if config.Env.K8S.ExternalNetwork.Interface, ret = GetValue[string](s, model.K8SExternalNetworkInterfaceSettingKey); !ret.OK {
		return ret
	}
	if config.Env.K8S.ExternalNetwork.ExcludeIPs, ret = GetValue[[]string](s, model.K8SExternalNetworkExcludeIPsSettingKey); !ret.OK {
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
	if config.Env.K8S.GeneratorWorker, ret = GetValue[int](s, model.K8SGeneratorWorkerSettingKey); !ret.OK {
		return ret
	}

	if config.Env.NFS.Server, ret = GetValue[string](s, model.NFSServerSettingKey); !ret.OK {
		return ret
	}
	if config.Env.NFS.Path, ret = GetValue[string](s, model.NFSPathSettingKey); !ret.OK {
		return ret
	}
	if config.Env.NFS.Storage, ret = GetValue[string](s, model.NFSStorageSettingKey); !ret.OK {
		return ret
	}

	if config.Env.Cheat.IP.Whitelist, ret = GetValue[[]string](s, model.CheatIPWhitelistSettingKey); !ret.OK {
		return ret
	}

	if config.Env.GeoCityDB, ret = GetValue[string](s, model.GeoCityDBSettingKey); !ret.OK {
		return ret
	}
	config.Tidy()
	if err := config.Save(); err != nil {
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
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
