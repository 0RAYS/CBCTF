package utils

import (
	i18n2 "CBCTF/internal/i18n"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

// In 实现 in 判断
func In(value interface{}, slice interface{}) bool {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return false
	}
	valueReflect := reflect.ValueOf(value)
	for i := 0; i < v.Len(); i++ {
		if reflect.DeepEqual(v.Index(i).Interface(), valueReflect.Interface()) {
			return true
		}
	}
	return false
}

// RandomString 生成随机uuid
func RandomString() string {
	return uuid.New().String()
}

// M i18n 将db返回的msg转换为gin返回的msg
func M(ctx *gin.Context, msg string) string {
	acceptLang := ctx.GetHeader("Accept-Language")
	return i18n.NewLocalizer(i18n2.Bundle, acceptLang).MustLocalize(&i18n.LocalizeConfig{
		MessageID: msg,
	})
}

// TidyRetData struct 转 json，同时去除一些敏感字段和结构体字段
func TidyRetData(data interface{}, bannedL ...string) []interface{} {
	var ret []interface{}
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Struct:
		banned := []string{"teams", "contests", "users", "model"}
		for _, v := range bannedL {
			banned = append(banned, v)
		}
		tmp := map[string]interface{}{}
		sType := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := sType.Field(i)
			fieldName := field.Name
			fieldValue := value.Field(i).Interface()
			if strings.ToLower(fieldName) == "model" {
				tmp["id"] = fieldValue.(gorm.Model).ID
			}
			if In(strings.ToLower(fieldName), banned) {
				continue
			}
			tmp[strings.ToLower(fieldName)] = fieldValue
		}
		ret = append(ret, tmp)
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			item := value.Index(i).Interface()
			for _, v := range TidyRetData(item, bannedL...) {
				ret = append(ret, v)
			}
		}
	default:
		return nil
	}
	return ret
}

// Form2Map 将Update请求的数据提取出被赋值的结果，为区分默认的0值和赋值的0值，表单中的字段都为指针类型
func Form2Map(s interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	types := reflect.TypeOf(s)
	values := reflect.ValueOf(s)
	n := values.NumField()
	for i := 0; i < n; i++ {
		key := types.Field(i).Tag.Get("json")
		if values.Field(i).Elem().IsValid() {
			value := values.Field(i).Elem().Interface()
			data[key] = value
		}
	}
	return data
}
