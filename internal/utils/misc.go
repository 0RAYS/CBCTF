package utils

import (
	"github.com/google/uuid"
	"reflect"
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

// Form2Map 将Update请求的数据提取出被赋值的结果, 为区分默认的0值和赋值的0值, 表单中的字段都为指针类型
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
