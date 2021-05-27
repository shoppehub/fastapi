package base

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ID = "_id"

// 默认类型
type BaseId struct {
	//主键
	Id *primitive.ObjectID `bson:"_id" json:"_id"`
	//创建时间
	CreatedAt *time.Time `bson:"createdAt" json:"createdAt,omitempty" update:"setOnInsert"`
	//修改时间
	UpdatedAt *time.Time `bson:"updatedAt" json:"updatedAt,omitempty"`
}
