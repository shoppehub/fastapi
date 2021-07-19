package session

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/shoppehub/conf"
	"github.com/shoppehub/fastapi/crud"
	"github.com/sirupsen/logrus"
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

	initLogin(r)

}

func initLogin(r *gin.Engine) {
	loginUrl := conf.GetString("auth.login.url")
	matchs := conf.GetStrings("auth.login.matchs")
	if loginUrl == "" || matchs == nil || len(matchs) == 0 {
		return
	}
	var regs []*regexp.Regexp
	for _, regStr := range matchs {
		reg := regexp.MustCompile(regStr)

		if reg == nil { //解释失败，返回nil
			logrus.Error("the auth.login.matchs err:", regStr)
			continue
		}
		regs = append(regs, reg)
	}

	r.Use(func(c *gin.Context) {
		if regs == nil {
			return
		}

		for _, reg := range regs {
			if reg.Match([]byte(c.Request.URL.Path)) {
				// 匹配到需要登录
				if GetUserId(c.Request) == "" {
					c.Redirect(http.StatusFound, loginUrl+"?target="+c.Request.URL.String())
					return
				}
			}
		}
	})

}
