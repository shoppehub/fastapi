package crud

import (
	"log"

	"github.com/iancoleman/strcase"
	"github.com/shoppehub/commons"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @param filterJSON mongo查询语句
func (instance *Resource) FindOneWithBson(filterJSON string, result interface{}, opts FindOneOptions) {

	if filterJSON == "" {
		return
	}
	var filter bson.M
	err := bson.UnmarshalExtJSON([]byte(filterJSON), true, &filter)
	if err != nil {
		log.Println(err, filterJSON)
		return
	}

	instance.FindOne(filter, result, opts)
}

func getTableName(opts FindOneOptions) string {
	return strcase.ToSnake(*opts.CollectionName)
}

// @param filterJSON mongo查询语句
func (instance *Resource) FindWithBson(filterJSON string, opts FindOptions) *commons.PagingResponse {

	var filter bson.M
	if filterJSON == "" {
		filter = bson.M{}
	} else {
		err := bson.UnmarshalExtJSON([]byte(filterJSON), true, &filter)
		if err != nil {
			logrus.Error(err, filterJSON)
			return nil
		}
	}

	return instance.Find(filter, opts)
}

func (instance *Resource) QueryWithBson(filterJSON string, opts FindOptions) *commons.PagingResponse {

	var filter mongo.Pipeline
	if filterJSON != "" {
		err := bson.UnmarshalExtJSON([]byte(filterJSON), true, &filter)
		if err != nil {
			logrus.Error(err, filterJSON)
			return nil
		}
	}

	return instance.Query(filter, opts)

}
