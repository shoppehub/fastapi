package fastapi

import (
	"reflect"
	"sort"
	"strings"
)

// 通过反射获取字段名称
func getFieldName(field reflect.StructField) string {
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
func getFieldValue(val reflect.Value) interface{} {

	if val.CanInterface() {
		return val.Interface()
	}
	return nil
}

// 判断字符串是否包含在数组内
func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	//index的取值：[0,len(str_array)]
	if index < len(str_array) && str_array[index] == target { //需要注意此处的判断，先判断 &&左侧的条件，如果不满足则结束此处判断，不会再进行右侧的判断
		return true
	}
	return false
}
