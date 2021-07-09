package session

import (
	"github.com/gin-gonic/gin"
	"github.com/shoppehub/conf"
	"github.com/shoppehub/fastapi/crud"
)

func Init(resource *crud.Resource, r *gin.Engine) {
	sid := conf.GetString("http.sid")
	if sid != "" {
		SidKey = sid
	}

	maxAge := conf.GetInt("http.maxAge")
	if maxAge != 0 {
		MaxAge = int(maxAge)
	}

	r.Use(func(c *gin.Context) {
		wrapUserSession(resource, c.Request)
	})

}
