package crud

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 生成id
func GenerateId(resource *Resource, key string, initVal int64) int64 {

	if key == "" {
		key = "common"
	}

	filter := bson.M{"key": key}

	update := bson.M{
		"$set": bson.M{"time": time.Now()},
		"$inc": bson.M{"value": int64(1)},
		"$setOnInsert": bson.M{
			"key": key,
		},
	}

	// 7) Create an instance of an options and set the desired options
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var val bson.M
	resource.DB.Collection("golbal_id").FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(&val)

	return initVal + (val["value"].(int64))
}
