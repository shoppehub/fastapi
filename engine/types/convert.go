package types

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shoppehub/fastapi/collection"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 转换数据类型
func Convert(value *map[string]interface{}, col collection.Collection) (map[string]interface{}, error) {

	mvalue := *value
	// mFields := make(map[string]*collection.CollectionField)
	// for _, field := range col.Fields {
	// 	f := field
	// 	mFields[field.Name] = &f
	// }

	// for key, v := range mvalue {
	// 	if mFields[key] != nil {
	// 		v1, err := ConvertField(v, *mFields[key])
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		result[key] = v1
	// 	}
	// }

	result := make(map[string]interface{})

	for _, field := range col.Fields {
		// 关联字段也忽略，没有值也忽略
		if mvalue[field.Name] == nil || field.RefField {
			continue
		}

		val, err := ConvertField(mvalue[field.Name], field)
		if err != nil {
			return nil, err
		}
		result[field.Name] = val
	}

	return result, nil
}

func ConvertField(value interface{}, field collection.CollectionField) (interface{}, error) {

	if field.Type == "" {
		return value, nil
	}

	rval := reflect.ValueOf(value)
	kind := rval.Kind()

	if strings.HasSuffix(field.Type, "[]") {
		return handerArray(kind, value, field)
	}
	if kind == reflect.Map {
		return handerMap(value, field)
	}
	if kind == reflect.String {
		return handerString(value, field)
	}
	return value, nil
}

func handerArray(kind reflect.Kind, value interface{}, field collection.CollectionField) (interface{}, error) {

	if kind != reflect.Slice && kind != reflect.Array {
		ejson, _ := json.Marshal(value)
		return nil, errors.New("value :" + string(ejson) + " is not array element")
	} else {
		m := value.([]interface{})
		val := make([]interface{}, len(m))
		for i, v := range m {
			v1, err := ConvertField(v, collection.CollectionField{
				Type: strings.TrimSuffix(field.Type, "[]"),
			})
			if err != nil {
				return nil, err
			}
			val[i] = v1
		}
		return val, nil
	}
}

func handerMap(value interface{}, field collection.CollectionField) (interface{}, error) {
	val := make(map[string]interface{})
	m := value.(map[string]interface{})
	if field.Fields == nil {
		ejson, _ := json.Marshal(m)
		return nil, errors.New("value :" + string(ejson) + " has no fields config")
	}
	for key, v := range m {
		for _, f := range field.Fields {
			if key == f.Name {
				v1, err := ConvertField(v, f)
				if err != nil {
					return nil, err
				}
				val[key] = v1
			}
		}
	}

	return val, nil
}

func handerString(value interface{}, field collection.CollectionField) (interface{}, error) {
	switch field.Type {
	case "objectId":
		{
			val, err := primitive.ObjectIDFromHex(value.(string))
			if err != nil {
				ejson, _ := json.Marshal(field)
				return nil, errors.New("value :" + value.(string) + " " + err.Error() + " ; the field config is  " + string(ejson))
			}
			return val, nil
		}
	case "int":
		{
			val, err := strconv.Atoi(value.(string))
			if err != nil {
				ejson, _ := json.Marshal(field)
				return nil, errors.New("value :" + value.(string) + " " + err.Error() + " ; the field config is  " + string(ejson))
			}
			return val, nil
		}
	case "long":
		{
			val, err := strconv.ParseInt(value.(string), 10, 64)
			if err != nil {
				ejson, _ := json.Marshal(field)
				return nil, errors.New("value :" + value.(string) + " " + err.Error() + " ; the field config is  " + string(ejson))
			}
			return val, nil
		}
	case "float":
		{
			// float 64
			val, err := strconv.ParseFloat(value.(string), 64)
			if err != nil {
				ejson, _ := json.Marshal(field)
				return nil, errors.New("value :" + value.(string) + " " + err.Error() + " ; the field config is " + string(ejson))
			}
			return val, nil
		}
	case "bool":
		{
			val, err := strconv.ParseBool(value.(string))
			if err != nil {
				ejson, _ := json.Marshal(field)
				return nil, errors.New("value :" + value.(string) + " " + err.Error() + " ; the field config is  " + string(ejson))
			}
			return val, nil
		}
	case "time":
		{
			str := value.(string)
			layout := "2006-01-02 15:04:05"
			if strings.IndexAny(str, ":") == -1 {
				layout = "2006-01-02"
			}
			val, err := time.Parse(layout, str)
			if err != nil {
				ejson, _ := json.Marshal(field)
				return nil, errors.New("value :" + value.(string) + " " + err.Error() + " ; the field config is " + string(ejson))
			}
			return val, nil
		}
	}
	return value, nil
}
