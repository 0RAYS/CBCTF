package utils

import (
	"github.com/google/uuid"
	"math/rand"
	"reflect"
	"strings"
	"time"
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

// UUID 生成随机uuid
func UUID() string {
	return uuid.New().String()
}

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandomString(n int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var builder strings.Builder
	builder.Grow(n)
	charsetLength := len(charset)
	for i := 0; i < n; i++ {
		builder.WriteByte(charset[rand.Intn(charsetLength)])
	}
	return builder.String()
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

func ToTitle(s string) string {
	if len(s) == 0 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// TidyPaginate 实现内存分页逻辑，处理各类边界场景
func TidyPaginate(length, limit, offset int) (int, int) {
	if length == 0 {
		return 0, 0
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= length {
		return 0, 0
	}
	var end int
	if limit <= 0 {
		end = length // limit<=0 时视为取全部剩余数据
	} else {
		end = offset + limit
	}
	if end > length {
		end = length
	}
	return end, limit
}
