package fastapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shoppehub/commons"
	"go.mongodb.org/mongo-driver/bson"
)

// 集合
type Collection struct {
	BaseId      `bson,inline`
	Name        string            `bson:"name,omitempty" json:"name,omitempty" update:"setOnInsert"`
	Description string            `bson:"description,omitempty" json:"description,omitempty"`
	Version     int64             `bson:"version,omitempty" json:"version,omitempty"`
	Extend      string            `bson:"extend,omitempty" json:"extend,omitempty"`
	Owner       string            `bson:"owner,omitempty" json:"owner,omitempty"`
	Fields      []CollectionField `bson:"fields,omitempty" json:"fields,omitempty"`
	Developers  []Developer       `bson:"developers,omitempty" json:"developers,omitempty"`
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
	Desc string `bson:"desc" json:"desc"`
	// 如果是内置对象模型，具体的字段是定义
	Fields []CollectionField `bson:"fields,omitempty" json:"fields,omitempty"`

	SelectOptions []SelectOptions `bson:"selectOptions,omitempty" json:"selectOptions,omitempty"`
	// 默认值
	DefaultValue string `bson:"defaultValue,omitempty" json:"defaultValue,omitempty"`
	// 验证规则
	Validate string `bson:"validate,omitempty" json:"validate,omitempty"`
}

type SelectOptions struct {
	Label string `bson:"label" json:"label"`
	Value string `bson:"value" json:"value"`
}

type CollectionRequest struct {
	commons.PagingRequest
	Name string `bson:"name,omitempty" json:"name,omitempty" update:"setOnInsert"`
}

// 创建集合
func CreateCollection(c *gin.Context) {

	var col Collection
	if err := c.ShouldBindJSON(&col); err != nil {
		c.JSON(http.StatusBadRequest, commons.ActionResponse{
			Success:    false,
			ErrMessage: err.Error(),
		})
		return
	}

	if col.Name == "" {
		c.JSON(http.StatusBadRequest, commons.ActionResponse{
			Success:    false,
			ErrCode:    10001,
			ErrMessage: "collectionName is empty",
		})
		return
	}
	// var dbCollection Collection
	// filter := bson.M{"name": col.Name}
	// dBResource.FindOne(filter, &dbCollection)

	// if &dbCollection != nil {
	// 	col.Id = dbCollection.Id
	// }
	wrapField(&col)

	result, saveErr := dBResource.SaveOrUpdateOne(&col, &UpdateOption{
		Filter: []string{0: "name"},
		Inc: []Inc{
			0: {
				Key:   "version",
				Value: 1,
			},
		},
	})
	if saveErr != nil {
		c.JSON(http.StatusBadRequest, commons.ActionResponse{
			Success:    false,
			ErrMessage: saveErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, commons.ActionResponse{
		Success: true,
		Data:    result,
	})
}

func wrapField(col *Collection) {
	var fields []CollectionField

	fields = append(fields, CollectionField{
		Name: "_id",
		Type: "ObjectId",
	})
	fields = append(fields, CollectionField{
		Name: "createdAt",
		Type: "Time",
	})
	fields = append(fields, CollectionField{
		Name: "updatedAt",
		Type: "Time",
	})

	if col.Fields != nil {
		for _, v := range col.Fields {
			if v.Name == "_id" || v.Name == "createdAt" || v.Name == "updatedAt" {
				continue
			}
			fields = append(fields, v)
		}
		col.Fields = fields
	}
}

// 查看集合
func GetCollection(c *gin.Context) {

	id := c.Params.ByName("id")

	var dbCollection Collection
	dBResource.FindById(id, &dbCollection, FindOneOptions{})

	c.JSON(http.StatusOK, commons.ActionResponse{
		Success: true,
		Data:    dbCollection,
	})
}

// 获取集合
func QueryCollection(c *gin.Context) {

	var req CollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, commons.ActionResponse{
			Success:    false,
			ErrMessage: err.Error(),
		})
		return
	}

	filter := bson.M{}

	var option FindOptions
	option.CurPage = req.CurPage
	option.PageSize = req.PageSize
	collection := "collection"
	option.CollectionName = &collection
	option.Results = &[]Collection{}

	if req.Name != "" {
		filter["name"] = req.Name
	}

	res := dBResource.Find(filter, option)

	c.JSON(http.StatusOK, res)
}
