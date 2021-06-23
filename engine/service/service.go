package service

import (
	"github.com/shoppehub/fastapi/collection"
	"github.com/shoppehub/fastapi/crud"
	"go.mongodb.org/mongo-driver/bson"
)

func (query *CollectionQuery) GetDbCollection(resource *crud.Resource) *collection.Collection {
	var dbCollection collection.Collection
	filter := bson.M{"name": query.ToString()}
	resource.FindOne(filter, &dbCollection, crud.FindOneOptions{CollectionName: GetCollectionName()})

	return &dbCollection
}
