package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
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

// UUID 生成随机uuid
func UUID() string {
	return uuid.New().String()
}

func EncryptMagic(magic string) string {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(magic)))[1:32]
	hash = fmt.Sprintf("%x", md5.Sum([]byte(hash)))
	return fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))
}

// UpdateOptions2Map 将Update请求的数据提取出被赋值的结果, 为区分默认的0值和赋值的0值, 表单中的字段都为指针类型
func UpdateOptions2Map(s interface{}) map[string]interface{} {
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

// TidyPaginate 实现内存分页逻辑, 处理各类边界场景
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
	return offset, end
}

// S2S 结构体转结构体, 有些低效, 但符合我预期, 先保留
func S2S[T any](o interface{}) (T, error) {
	var n T
	data, err := msgpack.Marshal(o)
	if err != nil {
		return n, err
	}
	err = msgpack.Unmarshal(data, &n)
	return n, err
}
