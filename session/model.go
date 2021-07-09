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
	Status      string     `bson:"status" json:"status"`
}
