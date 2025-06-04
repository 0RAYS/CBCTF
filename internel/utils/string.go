package utils

import (
	"github.com/google/uuid"
	"math/rand"
	"reflect"
	"time"
)

// In 实现 in 判断
func In(value any, slice any) bool {
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

// RandStr 生成随机字符串
func RandStr(n int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	result := make([]byte, n)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}
