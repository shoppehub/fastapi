package engine

import (
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shoppehub/commons"
	"github.com/shoppehub/fastapi/crud"
	"github.com/shoppehub/fastapi/engine/template"
	"github.com/shoppehub/fastapi/engine/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 根据id查询数据
func GetWithId(resource *crud.Resource, c *gin.Context) {
	var query CollectionQuery
	c.ShouldBindUri(&query)

	id := c.Params.ByName("id")
	//result := NewStruct(collection)

	dbCollection := query.GetDbCollection(resource, c)

	result := make(map[string]interface{})

	resource.FindById(id, result, crud.FindOneOptions{CollectionName: dbCollection.GetCollectionName()})
	c.JSON(http.StatusOK, commons.ActionResponse{
		Success: true,
		Data:    result,
	})
}

// 查询单个数据
func FindOne(resource *crud.Resource, c *gin.Context) {
	var query CollectionQuery
	c.ShouldBindUri(&query)

	var body CollectionBody
	c.ShouldBindJSON(&body)
	dbCollection := query.GetDbCollection(resource, c)
	result := make(map[string]interface{})
	fmt.Println(body)
	resource.FindOne(body.Filter, result,
		crud.FindOneOptions{
			CollectionName: dbCollection.GetCollectionName(),
			Sort:           body.Sort,
		})
	c.JSON(http.StatusOK, commons.ActionResponse{
		Success: true,
		Data:    result,
	})
}

// 保存数据
func Post(resource *crud.Resource, c *gin.Context) {
	var query CollectionQuery

	if err := c.ShouldBindUri(&query); err != nil {
		c.JSON(400, commons.ActionResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	var body CollectionBody
	c.ShouldBindJSON(&body)

	if body.Body == nil || len(body.Body) == 0 {
		c.JSON(400, commons.ActionResponse{Success: false, ErrMessage: "no body data"})
		return
	}

	dbCollection := query.GetDbCollection(resource, c)

	obj, err := types.Convert(&body.Body, *dbCollection)
	if err != nil {
		c.JSON(http.StatusOK, commons.ActionResponse{
			Success:    false,
			ErrMessage: err.Error(),
		})
		return
	}
	body.Body = obj

	result := Save(resource, *dbCollection, body)

	c.JSON(http.StatusOK, commons.ActionResponse{
		Success: true,
		Data:    result,
	})
}

// 根据id查询数据
func DeleteId(resource *crud.Resource, c *gin.Context) {
	var query CollectionQuery
	c.ShouldBindUri(&query)

	id := c.Params.ByName("id")
	//result := NewStruct(collection)

	dbCollection := query.GetDbCollection(resource, c)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(500, commons.ActionResponse{Success: false, ErrMessage: id + " is not objectId"})
		return
	}
	resource.DeleteById(*dbCollection.GetCollectionName(), oid)
	c.JSON(http.StatusOK, commons.ActionResponse{
		Success: true,
		Data:    id,
	})
}

// 保存数据
func Query(resource *crud.Resource, c *gin.Context) {
	var query CollectionQuery

	if err := c.ShouldBindUri(&query); err != nil {
		c.JSON(400, commons.ActionResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	var body CollectionBody
	c.ShouldBindJSON(&body)

	dbCollection := query.GetDbCollection(resource, c)

	var options crud.FindOptions
	options.CollectionName = dbCollection.GetCollectionName()

	result := resource.QueryWithBson(body.Aggregate, options)

	c.JSON(http.StatusOK, commons.ActionResponse{
		Success: true,
		Data:    result,
	})
}

// 保存数据
func Func(resource *crud.Resource, c *gin.Context) {
	var query CollectionQuery

	if err := c.ShouldBindUri(&query); err != nil {
		c.JSON(400, commons.ActionResponse{Success: false, ErrMessage: err.Error()})
		return
	}

	var body CollectionBody
	c.ShouldBindJSON(&body)

	dbCollection := query.GetDbCollection(resource, c)

	result, err := template.Render(resource, *dbCollection, query.Func, c)

	if err != nil {
		c.JSON(http.StatusOK, commons.ActionResponse{
			Success:    false,
			ErrMessage: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, commons.ActionResponse{
		Success: true,
		Data:    result,
	})
}
