package engine

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shoppehub/commons"
	"github.com/shoppehub/fastapi/collection"
	"github.com/shoppehub/fastapi/crud"
	"go.mongodb.org/mongo-driver/bson"
)

func (query *CollectionQuery) GetDbCollection(resource *crud.Resource, c *gin.Context) *collection.Collection {
	var dbCollection collection.Collection
	filter := bson.M{"name": query.toString()}
	resource.FindOne(filter, &dbCollection, crud.FindOneOptions{CollectionName: GetCollectionName()})

	if &dbCollection == nil {
		c.JSON(http.StatusNotFound, commons.ActionResponse{
			Success:    false,
			ErrMessage: query.toString() + " not found",
		})
		return nil
	}
	return &dbCollection
}
