package crud

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/shoppehub/commons"
	"github.com/shoppehub/fastapi/base"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 根据id查找元素
//
// @param id 数据主键
//
// @param result 返回的数据对象，比如 &user
func (instance *Resource) FindById(id string, result interface{}, opts FindOneOptions) {

	if id == "" {
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err, id+" is error")
		return
	}

	if opts.CollectionName == nil {
		collectionName := reflect.TypeOf(result).Elem().Name()
		opts.CollectionName = &collectionName
	}

	instance.FindOne(bson.M{"_id": objectId}, result, opts)
}

// @param filterJSON mongo查询语句
func (instance *Resource) FindOne(filter bson.M, result interface{}, opts FindOneOptions) {

	if filter == nil || len(filter) == 0 {
		return
	}

	tableName := getTableName(opts)
	findOptions := options.FindOne()
	var sort bson.D
	if opts.Sort != nil {
		for _, s := range opts.Sort {
			if s.Key != "" {
				sort = append(sort, bson.E{Key: s.Key, Value: s.Sort})
			}
		}
	}
	if sort == nil {
		sort = bson.D{{base.ID, -1}}
	}
	findOptions.SetSort(sort)

	singerResult := instance.DB.Collection(tableName).FindOne(context.Background(), filter, findOptions)
	singerResult.Decode(result)
}

// @param filterJSON mongo查询语句
func (instance *Resource) FindWithoutPaging(filter bson.M, opts FindOptions) []bson.M {

	tableName := getTableName(opts.FindOneOptions)

	cursor, err := instance.DB.Collection(tableName).Find(context.Background(), filter)

	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil {
		logrus.Error(err, filter)
	}

	return result
}

// @param filterJSON mongo查询语句
func (instance *Resource) Find(filter bson.M, opts FindOptions) *commons.PagingResponse {

	tableName := getTableName(opts.FindOneOptions)

	var response commons.PagingResponse

	total, err := instance.DB.Collection(tableName).CountDocuments(context.Background(), filter)
	if err != nil {
		str, _ := json.Marshal(filter)

		logrus.Error(err, str)
		return &response
	}

	response.TotalCount = total

	option := options.Find()

	if opts.PageSize == 0 {
		opts.PageSize = 15
	}
	response.CurPage = opts.CurPage
	response.PageSize = opts.PageSize
	option.Limit = &opts.PageSize

	if opts.CurPage > 0 {
		curPage := opts.CurPage
		pageSize := opts.PageSize
		skip := (curPage - 1) * pageSize
		option.Skip = &skip
	} else {
		skip := int64(0)
		option.Skip = &skip
	}

	cursor, err := instance.DB.Collection(tableName).Find(context.Background(), filter, option)
	if opts.Results != nil {
		if err = cursor.All(context.Background(), opts.Results); err != nil {
			logrus.Error(err, filter)
		}
		response.Data = opts.Results
	} else {
		var result []bson.M
		if err = cursor.All(context.Background(), &result); err != nil {
			logrus.Error(err, filter)
		}
		response.Data = result
	}

	response.Compute()

	return &response
}

// Example usage:
//
//		mongo.Pipeline{
//			{{"$group", bson.D{{"_id", "$state"}, {"totalPop", bson.D{{"$sum", "$pop"}}}}}},
//			{{"$match", bson.D{{"totalPop", bson.D{{"$gte", 10*1000*1000}}}}}},
//		}
func (instance *Resource) Query(pipeline mongo.Pipeline, opts FindOptions) *commons.PagingResponse {

	tableName := getTableName(opts.FindOneOptions)
	var response commons.PagingResponse

	var countPipeline mongo.Pipeline
	for _, p := range pipeline {
		countPipeline = append(countPipeline, p)
	}
	countPipeline = append(countPipeline, bson.D{{
		"$count", "totalCount",
	}})

	countCursor, countErr := instance.DB.Collection(tableName).Aggregate(context.Background(), countPipeline)
	if countErr != nil {
		str, _ := json.Marshal(countPipeline)
		logrus.Error("Aggregate Error of "+tableName, string(str))
		return &response
	}
	var response2 []commons.PagingResponse
	countCursor.All(context.Background(), &response2)
	response.TotalCount = response2[0].TotalCount

	if opts.PageSize == 0 {
		opts.PageSize = 15
	}
	if opts.CurPage == 0 {
		opts.CurPage = 1
	}
	response.CurPage = opts.CurPage
	response.PageSize = opts.PageSize

	pipeline = append(pipeline, bson.D{{
		"$limit", &opts.PageSize,
	}})
	skip := int64(0)
	if opts.CurPage > 0 {
		curPage := opts.CurPage
		pageSize := opts.PageSize
		skip = (curPage - 1) * pageSize

		pipeline = append(pipeline, bson.D{{
			"$skip", skip,
		}})
	}
	pipeline = append(pipeline, bson.D{{
		"$skip", skip,
	}})

	cursor, err := instance.DB.Collection(tableName).Aggregate(context.Background(), pipeline)
	if cursor == nil {
		str, _ := json.Marshal(pipeline)
		logrus.Error("Aggregate Error of "+tableName, string(str))
		return &response
	}

	if opts.Results != nil {
		if err = cursor.All(context.Background(), opts.Results); err != nil {
			logrus.Error(err, pipeline)
		}
		response.Data = opts.Results
	} else {
		var result []bson.M
		if err = cursor.All(context.Background(), &result); err != nil {
			logrus.Error(err, pipeline)
		}
		response.Data = result
	}
	response.Compute()
	return &response
}

func (instance *Resource) QueryWithoutPaging(pipeline mongo.Pipeline, opts FindOptions) ([]bson.M, error) {

	tableName := getTableName(opts.FindOneOptions)

	var result []bson.M

	cursor, err := instance.DB.Collection(tableName).Aggregate(context.Background(), pipeline)

	if err = cursor.All(context.Background(), &result); err != nil {
		logrus.Error(err, pipeline)
		return nil, err
	}

	return result, nil
}
