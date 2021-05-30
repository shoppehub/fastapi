package collection

import (
	"strings"

	"github.com/shoppehub/commons"
	"github.com/shoppehub/fastapi/base"
)

// 集合
type Collection struct {
	base.BaseId `bson,inline`
	Name        string              `bson:"name,omitempty" json:"name,omitempty" update:"setOnInsert"`
	Description string              `bson:"description,omitempty" json:"description,omitempty"`
	Version     int64               `bson:"version,omitempty" json:"version,omitempty"`
	Extend      string              `bson:"extend,omitempty" json:"extend,omitempty"`
	Owner       string              `bson:"owner,omitempty" json:"owner,omitempty"`
	Fields      []CollectionField   `bson:"fields,omitempty" json:"fields,omitempty"`
	Developers  []Developer         `bson:"developers,omitempty" json:"developers,omitempty"`
	Functions   map[string]Function `bson:"functions,omitempty" json:"functions,omitempty"`
}

// 开发者
type Developer struct {
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	Time string `bson:"time,omitempty" json:"time,omitempty"`
	Desc string `bson:"desc,omitempty" json:"desc,omitempty"`
}

// 集合字段类型
type CollectionField struct {
	// 字段名称
	Name string `bson:"name" json:"name"`
	// 字段中文名称
	Title string `bson:"title" json:"title"`
	// 字段类型
	Type string `bson:"type" json:"type"`
	// 字段描述
	Desc string `bson:"desc,omitempty" json:"desc,omitempty"`

	RefField bool `bson:"refField,omitempty" json:"refField,omitempty"`

	SetOnInsert bool `bson:"setOnInsert,omitempty" json:"setOnInsert,omitempty"`
	// 如果是内置对象模型，具体的字段是定义
	Fields []CollectionField `bson:"fields,omitempty" json:"fields,omitempty"`

	SelectOptions []SelectOptions `bson:"selectOptions,omitempty" json:"selectOptions,omitempty"`
	// 值
	Value interface{} `bson:"value,omitempty" json:"value,omitempty"`
	// 默认值
	DefaultValue interface{} `bson:"defaultValue,omitempty" json:"defaultValue,omitempty"`
	// 验证规则
	Validate string `bson:"validate,omitempty" json:"validate,omitempty"`
}

type SelectOptions struct {
	Label string `bson:"label" json:"label"`
	Value string `bson:"value" json:"value"`
}

// 开发者
type Function struct {
	Params   []CollectionField `bson:"params,omitempty" json:"params,omitempty"`
	Template string            `bson:"template,omitempty" json:"template,omitempty"`
}

type CollectionRequest struct {
	commons.PagingRequest
	Name string `bson:"name,omitempty" json:"name,omitempty" update:"setOnInsert"`
}

func (collection *Collection) GetCollectionName() *string {
	collectionName := strings.ReplaceAll(collection.Name, "/", "_")
	return &collectionName
}
