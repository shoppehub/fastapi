package base

import (
	"reflect"
	"strings"
)

// 通过反射获取字段名称
func GetFieldName(field reflect.StructField) string {
	name := field.Tag.Get("bson")
	name = strings.Split(name, ",")[0]

	if field.Tag.Get("update") == "skip" {
		return ""
	}

	if name == "" {
		name = strings.ToLower(field.Name)
	}

	return name
}

// 获取字段的实际值
func GetFieldValue(val reflect.Value) interface{} {

	if val.CanInterface() {
		return val.Interface()
	}
	return nil
}
