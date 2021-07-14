package session

import (
	"time"

	"github.com/shoppehub/fastapi/base"
)

// 用户session模型
type UserSession struct {
	base.BaseId `bson,inline`
	Uid         string     `bson:"uid" json:"uid"`
	Expires     *time.Time `bson:"expires,omitempty" json:"expires,omitempty"`
	Agent       string     `bson:"agent,omitempty" json:"agent,omitempty"`
	Ip          string     `bson:"ip,omitempty" json:"ip,omitempty"`
	// login、logout
	Status string `bson:"status" json:"status"`
	// 头像
	Avatar   string `bson:"avatar" json:"avatar"`
	NickName string `bson:"nickName" json:"nickName"`

	MaxAge int64 `bson:"-"`
}
