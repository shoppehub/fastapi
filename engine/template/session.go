package template

import (
	"reflect"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
	"github.com/shoppehub/fastapi/crud"
	"github.com/shoppehub/fastapi/session"
	"github.com/shoppehub/sjet"
)

func initSession(resource *crud.Resource) {
	sjet.RegCustomFunc("login", func(c *gin.Context) jet.Func {
		return loginFunc(resource, c)
	})

	sjet.RegCustomFunc("logout", func(c *gin.Context) jet.Func {
		return logoutFunc(resource, c)
	})

	sjet.RegCustomFunc("getUserId", func(c *gin.Context) jet.Func {
		return getUserIdFunc(resource, c)
	})

	sjet.RegCustomFunc("getUserSession", func(c *gin.Context) jet.Func {
		return getUserSessionFunc(resource, c)
	})
}

func loginFunc(resource *crud.Resource, c *gin.Context) jet.Func {

	return func(a jet.Arguments) reflect.Value {
		mm := a.Get(0).Interface().(map[string]interface{})

		userSession := session.UserSession{
			Uid: mm["uid"].(string),
		}
		if mm["avatar"] != nil {
			userSession.Avatar = mm["avatar"].(string)
		}
		if mm["nickName"] != nil {
			userSession.NickName = mm["nickName"].(string)
		}
		if mm["maxAge"] != nil {
			userSession.MaxAge = int64(mm["maxAge"].(int64))
		}

		s, _ := session.NewUserSession(resource, userSession, c.Request, c.Writer)

		return reflect.ValueOf(s)
	}
}

func logoutFunc(resource *crud.Resource, c *gin.Context) jet.Func {

	return func(a jet.Arguments) reflect.Value {
		session.RemoveUserSession(resource, c.Request)
		return reflect.ValueOf("")
	}
}

func getUserSessionFunc(resource *crud.Resource, c *gin.Context) jet.Func {

	return func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(session.GetUserSession(c.Request))
	}
}

func getUserIdFunc(resource *crud.Resource, c *gin.Context) jet.Func {

	return func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(session.GetUserId(c.Request))
	}
}
