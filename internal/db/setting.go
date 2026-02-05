package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
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
	Value *string
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

func (s *SettingRepo) InitSettings() model.RetVal {
	for _, setting := range []model.Setting{
		{Key: model.AsyncQLogLevelSettingKey, Value: model.SettingValue{V: config.Env.AsyncQ.Log.Level}},
		{Key: model.AsyncQConcurrencySettingKey, Value: model.SettingValue{V: config.Env.AsyncQ.Concurrency}},

		{Key: model.GinRateLimitGlobalSettingKey, Value: model.SettingValue{V: config.Env.Gin.RateLimit.Global}},
		{Key: model.GinRateLimitWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Gin.RateLimit.Whitelist}},
		{Key: model.GinCORSSettingKey, Value: model.SettingValue{V: config.Env.Gin.CORS}},
		{Key: model.GinLogWhitelistSettingKey, Value: model.SettingValue{V: config.Env.Gin.Log.Whitelist}},
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

	if config.Env.AsyncQ.Log.Level, ret = GetValue[string](s, model.AsyncQLogLevelSettingKey); !ret.OK {
		return ret
	}
	if config.Env.AsyncQ.Concurrency, ret = GetValue[float64](s, model.AsyncQConcurrencySettingKey); !ret.OK {
		return ret
	}

	if config.Env.Gin.RateLimit.Global, ret = GetValue[float64](s, model.GinRateLimitGlobalSettingKey); !ret.OK {
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

	if targetType.Kind() == reflect.Slice {
		sliceVal, ok := data.Value.V.([]any)
		if !ok {
			return zero, model.RetVal{
				Msg: i18n.Model.Setting.InvalidType,
				Attr: map[string]any{
					"Key":         data.Key,
					"Type":        targetType.String(),
					"InvalidType": reflect.TypeOf(data.Value.V).String(),
				},
			}
		}

		elemType := targetType.Elem()

		resultSlice := reflect.MakeSlice(targetType, 0, len(sliceVal))

		for _, item := range sliceVal {
			itemVal := reflect.ValueOf(item)

			if !itemVal.Type().AssignableTo(elemType) {
				return zero, model.RetVal{
					Msg: i18n.Model.Setting.InvalidType,
					Attr: map[string]any{
						"Key":         data.Key,
						"Type":        elemType.String(),
						"InvalidType": itemVal.Type().String(),
					},
				}
			}

			resultSlice = reflect.Append(resultSlice, itemVal)
		}

		return resultSlice.Interface().(T), model.SuccessRetVal()
	}

	val, ok := data.Value.V.(T)
	if !ok {
		return zero, model.RetVal{
			Msg: i18n.Model.Setting.InvalidType,
			Attr: map[string]any{
				"Key":         data.Key,
				"Type":        targetType.String(),
				"InvalidType": reflect.TypeOf(data.Value.V).String(),
			},
		}
	}

	return val, model.SuccessRetVal()
}
