package engine

import (
	"github.com/shoppehub/commons"
	"github.com/shoppehub/fastapi/crud"
)

type CollectionBody struct {
	Filter    map[string]interface{} `json:"filter"`
	Body      map[string]interface{} `json:"body" `
	Sort      []crud.Sort            `json:"sort" `
	Page      commons.PagingRequest  `json:"page" `
	Aggregate string                 `json:"aggregate" `
}

func GetCollectionName() *string {
	name := "collection"
	return &name
}

type CollectionQuery struct {
	Group      string `uri:"group" binding:"required"`
	Collection string `uri:"collection" binding:"required"`

	Func string `uri:"func"`
}

func (q *CollectionQuery) toString() string {
	str := q.Group + "/" + q.Collection
	return str
}

func (q *CollectionQuery) getCollectionName() *string {
	name := q.Group + "_" + q.Collection
	return &name
}
