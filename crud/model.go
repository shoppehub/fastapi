package crud

import (
	"fmt"
	"reflect"

	"github.com/iancoleman/strcase"
	"github.com/shoppehub/commons"
)

var typeRegistry = make(map[string]reflect.Type)

// 注册类型
func RegisterType(elem interface{}) {
	t := reflect.TypeOf(elem).Elem()
	name := strcase.ToSnake(t.Name())
	typeRegistry[name] = t
}

// 根据类型名称初始化类型
func NewStruct(name string) interface{} {
	name = strcase.ToSnake(name)
	elem, ok := typeRegistry[name]
	if !ok {
		panic(name + " has no registerType")
	}
	return reflect.New(elem).Elem().Interface()
}

type Operator string

const (
	OP_IN  Operator = "$in"
	OP_NIN Operator = "$nin"

	OP_EQ Operator = "$eq"
	OP_NE Operator = "$ne"

	OP_GT  Operator = "$gt"
	OP_GTE Operator = "$gte"

	OP_LT  Operator = "$lt"
	OP_LTE Operator = "$lte"

	OP_L_AND Operator = "$and"
	OP_L_NOT Operator = "$not"
	OP_L_NOR Operator = "$nor"
	OP_L_OR  Operator = "$or"

	OP_EL_EXISTS Operator = "$exists"
	OP_EL_TYPE   Operator = "$type"

	OP_EVA_EXPR       Operator = "$expr"
	OP_EVA_JSONSCHEMA Operator = "$jsonSchema"
	OP_EVA_MOD        Operator = "$mod"
	OP_EVA_REGEX      Operator = "$regex"
	OP_EVA_TEXT       Operator = "$text"
	OP_EVA_WHERE      Operator = "$where"

	OP_GEO_INTERSECTS  Operator = "$geoIntersects"
	OP_GEO_WITHIN      Operator = "$geoWithin"
	OP_GEO_NEAR        Operator = "$near"
	OP_GEO_NEAR_SPHERE Operator = "$nearSphere"

	OP_ARRAY_ELEMMATCH Operator = "$elemMatch"
)

func (op Operator) toString() string {
	return string(op)
}

// 过滤类型
type Filter struct {
	Key      string
	Operator string
	Value    interface{}
}

type FindOneOptions struct {
	// 集合名称
	CollectionName *string
	// 排序规则，默认是按照 _id 降序排序
	Sort []Sort
}

type FindOptions struct {
	commons.PagingRequest

	FindOneOptions

	// 返回的数据
	Results interface{}
}

type Sort struct {
	// 排序字段
	Key string `json:"key"`
	// -1 是降序，1是升序
	Sort int64 `json:"sort"`
}

func (op *FindOptions) SetCollectionName(colName string) *FindOptions {
	op.CollectionName = &colName
	return op
}

func (op *FindOptions) SetSort(sort ...Sort) *FindOptions {
	op.Sort = sort
	return op
}

func (op *FindOneOptions) SetSort(sort ...Sort) *FindOneOptions {
	op.Sort = sort
	return op
}

func CreateFindOneOptions(collectionName string) *FindOneOptions {
	op := &FindOneOptions{}
	op.CollectionName = &collectionName
	return op
}

func ToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}
