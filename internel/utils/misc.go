package utils

import "reflect"

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
	if limit < 0 {
		end = length // limit<=0 时视为取全部剩余数据
	} else {
		end = offset + limit
	}
	if end > length {
		end = length
	}
	return offset, end
}

func Ptr[T any](v T) *T {
	return &v
}
