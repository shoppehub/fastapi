package crud

import (
	"context"

	"github.com/shoppehub/fastapi/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 集合的定义
type DeleteOption struct {
	// 集合名称
	CollectionName string
	// 过滤字段，如果设置则以字段为查询条件进行修改，如果都没有值则失败
	Filter []Filter
}

// 删除的filter

//根据主键删除数据
func (resource *Resource) DeleteById(tableName string, id primitive.ObjectID) error {
	return resource.DeleteAny(&DeleteOption{
		CollectionName: tableName,
		Filter: []Filter{
			0: {
				Key:   base.ID,
				Value: id,
			},
		},
	})
}

//根据条件删除数据
func (resource *Resource) DeleteAny(deleteOption *DeleteOption) error {
	filter := make(map[string]interface{})
	for _, v := range deleteOption.Filter {
		if v.Operator == "" {
			filter[v.Key] = v.Value
		} else {
			//  {"k":{"$in":[]}}

			filter[v.Key] = bson.M{v.Operator: v.Value}
		}
	}
	_, err := resource.DB.Collection(deleteOption.CollectionName).DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
