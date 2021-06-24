package types

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shoppehub/fastapi/collection"
	"github.com/shoppehub/fastapi/crud"
	"github.com/shoppehub/fastapi/engine/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 转换数据类型
func Convert(resource *crud.Resource, value *map[string]interface{}, col collection.Collection) (map[string]interface{}, error) {

	mvalue := *value
	result := make(map[string]interface{})
	for _, field := range col.Fields {

		// logrus.Error(field.Name, mvalue[field.Name], field.RefField)

		// 关联字段也忽略，没有值也忽略
		// if mvalue[field.Name] == nil || field.RefField {
		// 	continue
		// }
		if mvalue[field.Name] == nil {
			continue
		}

		val, err := ConvertField(resource, mvalue[field.Name], field)
		if err != nil {
			return nil, err
		}
		result[field.Name] = val
	}

	return result, nil
}

func ConvertField(resource *crud.Resource, value interface{}, field collection.CollectionField) (interface{}, error) {

	if field.Type == "" {
		return value, nil
	}

	rval := reflect.ValueOf(value)
	kind := rval.Kind()
	if strings.Contains(field.Type, "/") {
		// 外部内容
		strs := strings.Split(field.Type, "/")
		q := service.CollectionQuery{
			Group:      strs[0],
			Collection: strs[1],
		}
		dbCollection := q.GetDbCollection(resource)
		if &dbCollection == nil {
			return value, nil
		}
		field.Fields = dbCollection.Fields
	}

	if field.Type == "object" {
		return handerMap(resource, value, field)
	}

	if strings.HasSuffix(field.Type, "[]") {
		return handerArray(resource, kind, value, field)
	}
	if kind == reflect.Map {
		return handerMap(resource, value, field)
	}
	if kind == reflect.String {
		return handerString(value, field)
	}
	return value, nil
}

func handerArray(resource *crud.Resource, kind reflect.Kind, value interface{}, field collection.CollectionField) (interface{}, error) {

	if kind != reflect.Slice && kind != reflect.Array {
		ejson, _ := json.Marshal(value)
		return nil, errors.New("value :" + string(ejson) + " is not array element")
	} else {
		m := value.([]interface{})
		val := make([]interface{}, len(m))
		for i, v := range m {

			otype := strings.TrimSuffix(field.Type, "[]")
			if otype == "object" {
				objectVal, err := handerMap(resource, v, field)
				if err != nil {
					ejson, _ := json.Marshal(v)
					return nil, errors.New("value :" + string(ejson) + " convert to object err")
				}
				val[i] = objectVal
			} else {
				v1, err := ConvertField(resource, v, collection.CollectionField{
					Type: strings.TrimSuffix(field.Type, "[]"),
				})
				if err != nil {
					return nil, err
				}
				val[i] = v1
			}
		}
		return val, nil
	}
}

func handerMap(resource *crud.Resource, value interface{}, field collection.CollectionField) (interface{}, error) {
	val := make(map[string]interface{})
	m := value.(map[string]interface{})
	if field.Fields == nil {
		ejson, _ := json.Marshal(m)
		return nil, errors.New("value :" + string(ejson) + " has no fields config")
	}
	for key, v := range m {
		for _, f := range field.Fields {
			if key == f.Name {
				v1, err := ConvertField(resource, v, f)
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
