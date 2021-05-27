package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	filename "github.com/keepeye/logrus-filename"
	"github.com/shoppehub/conf"
	"github.com/shoppehub/fastapi"
	"github.com/sirupsen/logrus"
)

func main() {
	filenameHook := filename.NewHook()
	logrus.AddHook(filenameHook)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Debug("init server")

	conf.Init("")

	instance, err := fastapi.NewDB(conf.GetString("mongodb.url"), conf.GetString("mongodb.dbname"))

	if err != nil {
		logrus.Error(err)
		return
	}
	r := New()
	fastapi.InitApi(instance, r)

	port := conf.GetString("port")
	if port == "" {
		port = "4000"
	}
	logrus.Info("start server on " + port)

	r.Run("0.0.0.0:" + port)
}

func New() *gin.Engine {

	//log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	// 使用默认中间件（logger和recovery）
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	return r
}
