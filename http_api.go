package fastapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shoppehub/commons"
	"go.mongodb.org/mongo-driver/bson"
)

var dBResource *Resource

type CollectionQuery struct {
	group      string
	collection string
}

func (q CollectionQuery) toString() string {
	return q.group + "/" + q.collection
}

func (q CollectionQuery) getCollectionName() *string {
	name := q.group + "_" + q.collection
	return &name
}

func GetWithId(c *gin.Context) {
	var query CollectionQuery
	c.ShouldBindUri(&query)

	id := c.Params.ByName("id")
	//result := NewStruct(collection)
	result := make(map[string]interface{})

	var dbCollection Collection
	filter := bson.M{"name": query.toString()}
	dBResource.FindOne(filter, &dbCollection, FindOneOptions{})

	if &dbCollection != nil {
		c.JSON(http.StatusNotFound, commons.ActionResponse{
			Success:    false,
			ErrMessage: query.toString() + " not found",
		})
		return
	}

	dBResource.FindById(id, result, FindOneOptions{CollectionName: query.getCollectionName()})
	c.JSON(http.StatusOK, result)
}

func filterObject(collection Collection, object map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	return result

}
