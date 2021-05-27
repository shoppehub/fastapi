package engine

import (
	"github.com/shoppehub/fastapi/crud"
	"go.mongodb.org/mongo-driver/bson"
)

type CollectionBody struct {
	Filter bson.M                 `json:"filter"`
	Body   map[string]interface{} `json:"body" `
	Sort   []crud.Sort            `json:"sort" `
}

func GetCollectionName() *string {
	name := "collection"
	return &name
}

type CollectionQuery struct {
	Group      string `uri:"group" binding:"required"`
	Collection string `uri:"collection" binding:"required"`
}

func (q *CollectionQuery) toString() string {
	str := q.Group + "/" + q.Collection
	return str
}

func (q *CollectionQuery) getCollectionName() *string {
	name := q.Group + "_" + q.Collection
	return &name
}
