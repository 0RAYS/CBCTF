package utils

import "reflect"

func GetFieldByJSONTag(obj any, tag string) any {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("json") == tag {
			return v.Field(i).Interface()
		}
	}
	return nil
}
