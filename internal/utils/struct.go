package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

type Data interface {
	Tidy() struct{}
}

// TidyResponse 递归处理结构体，确保 admin:"true" 的字段被忽略
func TidyResponse(ctx context.Context, v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil
		}
		return TidyResponse(ctx, rv.Elem().Interface())
	case reflect.Slice, reflect.Array:
		length := rv.Len()
		result := make([]interface{}, length)
		for i := 0; i < length; i++ {
			val := TidyResponse(ctx, rv.Index(i).Interface())
			result[i] = val
		}
		return result
	case reflect.Map:
		result := make(map[string]interface{})
		for _, key := range rv.MapKeys() {
			val := TidyResponse(ctx, rv.MapIndex(key).Interface())
			result[fmt.Sprintf("%v", key.Interface())] = val
		}
		return result
	case reflect.Struct:
		result := make(map[string]interface{})

		for i := 0; i < rv.NumField(); i++ {
			field := rt.Field(i)
			jsonTag := field.Tag.Get("json")
			showTag := field.Tag.Get("admin")
			fieldValue := rv.Field(i)
			if showTag == "true" || jsonTag == "-" {
				continue
			}
			val := TidyResponse(ctx, fieldValue.Interface())
			if jsonTag != "" {
				result[jsonTag] = val
			} else {
				result[field.Name] = val
			}
		}
		return result
	default:
		return v
	}
}

func MarshalJSON(v interface{}) ([]byte, error) {
	data := TidyResponse(context.Background(), v)
	return json.Marshal(data)
}
