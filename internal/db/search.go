package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"reflect"
)

func Search(m model.Model, limit, offset int, options GetOptions) (any, int64, model.RetVal) {
	var count int64
	countOptions := CountOptions{
		Conditions: options.Conditions,
		Search:     options.Search,
		Deleted:    options.Deleted,
	}
	res := applyCountOptions(DB.Table(m.TableName()), countOptions).Count(&count)
	if res.Error != nil {
		log.Logger.Warningf("Failed to search %s: %s", m.ModelName(), res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": m.ModelName(), "Error": res.Error.Error()}}
	}
	ms := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(m)), 0, 0).Interface()
	res = applyGetOptions(DB.Table(m.TableName()), options).Limit(limit).Offset(offset).Find(&ms)
	if res.Error != nil {
		log.Logger.Warningf("Failed to search %s: %s", m.ModelName(), res.Error)
		return nil, 0, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": m.ModelName(), "Error": res.Error.Error()}}
	}
	return ms, count, model.SuccessRetVal()
}
