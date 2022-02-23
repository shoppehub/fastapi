package crud

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/shoppehub/fastapi/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Upsert = true

const (
	SetOnInsert = "setOnInsert"
)

// 集合的定义
type UpdateOption struct {
	// 集合名称
	CollectionName *string
	// 过滤字段，如果设置则以字段为查询条件进行修改，如果都没有值则失败
	Filter []string
	// 自增字段
	Inc []Inc
}

type Inc struct {
	Key   string
	Value int64
}

type updateField struct {
	setElements         bson.D
	setOnInsertElements bson.D
	unsetElements       bson.D
	fieldMap            map[string]interface{}
}

// 修改一个数据
func (instance *Resource) SaveOrUpdateOne(u interface{}, updateOptions ...*UpdateOption) (interface{}, error) {

	var collectionName string
	var filterField []string

	var option *UpdateOption

	if len(updateOptions) > 0 {
		option = updateOptions[0]
		if option.CollectionName != nil {
			collectionName = *option.CollectionName
		}
		filterField = option.Filter
	}

	if collectionName == "" {
		collectionName = strcase.ToSnake(reflect.TypeOf(u).Elem().Name())
	}

	updateField := initUpdateField()

	updateField.eachField(u, nil)

	update := bson.D{}
	incKeys := make(map[string]string)
	if option != nil && len(option.Inc) > 0 {
		var incElements bson.D
		for _, v := range option.Inc {
			incElements = append(incElements, bson.E{v.Key, v.Value})
			incKeys[v.Key] = "true"
		}
		update = append(update, bson.E{"$inc", incElements})
	}
	update = append(update, bson.E{"$setOnInsert", updateField.setOnInsertElements})

	setElements := bson.D{}
	for _, element := range updateField.setElements {
		if incKeys[element.Key] != "" {
			continue
		}
		setElements = append(setElements, bson.E{
			element.Key, element.Value,
		})
	}
	update = append(update, bson.E{"$set", setElements})

	var filter bson.D

	if len(filterField) == 0 || base.In(base.ID, filterField) {
		if updateField.fieldMap[base.ID] == nil {
			updateField.fieldMap[base.ID] = primitive.NewObjectID()
		}
		filter = append(filter, bson.E{base.ID, updateField.fieldMap[base.ID]})
	} else {
		for _, f := range filterField {
			obj := updateField.fieldMap[f]
			if obj == nil {
				return nil, errors.New("the filter field " + f + " has no value")
			}
			filter = append(filter, bson.E{f, obj})
		}
	}
	session, startSessionErr := instance.Client.StartSession()
	if startSessionErr != nil {
		logrus.Error(startSessionErr)
		if session != nil {
			session.EndSession(context.TODO())
		}
		return nil, startSessionErr
	}
	defer session.EndSession(context.TODO())

	_, err := session.Client().Database(instance.DB.Name()).Collection(collectionName).UpdateOne(context.Background(), filter, update, &options.UpdateOptions{Upsert: &Upsert})
	if err != nil {
		return nil, err
	}
	result := reflect.New(reflect.TypeOf(u))
	err = session.Client().Database(instance.DB.Name()).Collection(collectionName).FindOne(context.Background(), filter).Decode(result.Interface())
	return result.Elem().Interface(), err
}

// 遍历结构体的每个字段
func (updateField *updateField) eachField(o interface{}, parent map[string]interface{}) {

	kind := reflect.ValueOf(o).Kind()
	var keys reflect.Type
	var obj reflect.Value

	if kind == reflect.Ptr || kind == reflect.Interface {
		obj = reflect.ValueOf(o).Elem()
		keys = reflect.TypeOf(o).Elem()
	} else {
		obj = reflect.ValueOf(o)
		keys = reflect.TypeOf(o)
	}

	for i := 0; i < keys.NumField(); i++ {
		field := keys.Field(i)
		// 匿名结构体
		if field.Anonymous {
			for fi := 0; fi < field.Type.NumField(); fi++ {
				efield := field.Type.Field(fi)
				val := obj.Field(i).FieldByName(efield.Name)
				updateField.wrapField(efield, val, parent)
			}
		} else {
			val := obj.FieldByName(field.Name)
			updateField.wrapField(field, val, parent)
		}
	}
}

// 包装每个字段，完成mongo修改的需要
func (updateField *updateField) wrapField(field reflect.StructField, value reflect.Value, parent map[string]interface{}) {

	if value.IsZero() {
		return
	}
	name := base.GetFieldName(field)
	// 忽略该字段
	if name == "" {
		return
	}

	var val interface{}
	// 说明是嵌套字段，需要继续遍历属性
	if field.Type.Kind() == reflect.Struct && field.Type.String() != "time.Time" {
		innerValue := make(map[string]interface{})

		updateField.eachField(value.Interface(), innerValue)
		val = innerValue
	} else {
		val = base.GetFieldValue(value)
	}

	if val == nil {
		return
	}

	if parent == nil {
		if updateField.fieldMap[name] != nil {
			return
		}
		updateField.fieldMap[name] = val

		if reflect.ValueOf(val).Kind() == reflect.String && val.(string) == "null" {
			updateField.unsetElements = append(updateField.unsetElements, bson.E{name, ""})
			return
		}

		tag := field.Tag.Get("update")
		if tag == "setOnInsert" {
			updateField.setOnInsertElements = append(updateField.setOnInsertElements, bson.E{name, val})
		} else {
			updateField.setElements = append(updateField.setElements, bson.E{name, val})
		}
		return
	} else {
		parent[name] = val
	}
}

// 初始化 baseId 的字段值
func initBaseIdField(updateField *updateField) {
	t := time.Now()
	updateField.setOnInsertElements = append(updateField.setOnInsertElements, bson.E{"createdAt", t})
	updateField.fieldMap["createdAt"] = t

	updateField.setElements = append(updateField.setElements, bson.E{"updatedAt", t})
	updateField.fieldMap["updatedAt"] = t
}

// 初始化字段
func initUpdateField() *updateField {
	var setElements bson.D
	var setOnInsertElements bson.D
	fieldMap := make(map[string]interface{})

	updateField := &updateField{
		setElements:         setElements,
		setOnInsertElements: setOnInsertElements,
		fieldMap:            fieldMap,
	}
	initBaseIdField(updateField)

	return updateField
}
