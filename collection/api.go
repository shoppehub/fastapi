package collection

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shoppehub/commons"
	"github.com/shoppehub/fastapi/crud"
	"go.mongodb.org/mongo-driver/bson"
)

// 创建集合
func CreateCollection(dBResource *crud.Resource, c *gin.Context) {

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

	result, saveErr := dBResource.SaveOrUpdateOne(&col, &crud.UpdateOption{
		Filter: []string{0: "name"},
		Inc: []crud.Inc{
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

// 查看集合
func GetCollection(dBResource *crud.Resource, c *gin.Context) {

	id := c.Params.ByName("id")

	var dbCollection Collection
	dBResource.FindById(id, &dbCollection, crud.FindOneOptions{})

	c.JSON(http.StatusOK, commons.ActionResponse{
		Success: true,
		Data:    dbCollection,
	})
}

// 查看集合
func FindOneCollection(dBResource *crud.Resource, name string) *map[string]CollectionField {
	var dbCollection Collection

	filter := bson.M{
		"name": name,
	}

	dBResource.FindOne(filter, &dbCollection, crud.FindOneOptions{})

	fieldMap := make(map[string]CollectionField)

	var eachField func(prefix string, fields *[]CollectionField)

	eachField = func(prefix string, fields *[]CollectionField) {
		for _, field := range *fields {

			name := field.Name

			if prefix != "" {
				name += "." + field.Name
			}

			var fi CollectionField

			fi.Type = field.Type
			fi.DefaultValue = field.DefaultValue
			fi.SelectOptions = field.SelectOptions
			if len(field.Fields) > 0 {
				eachField(name, &field.Fields)
			}
			fieldMap[name] = fi
		}
	}
	fields := dbCollection.Fields
	eachField("", &fields)

	return &fieldMap
}

// 获取集合
func QueryCollection(dBResource *crud.Resource, c *gin.Context) {

	var req CollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, commons.ActionResponse{
			Success:    false,
			ErrMessage: err.Error(),
		})
		return
	}

	filter := bson.M{}

	var option crud.FindOptions
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
