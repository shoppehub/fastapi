package fastapi

import (
	"github.com/gin-gonic/gin"
	"github.com/shoppehub/fastapi/collection"
	"github.com/shoppehub/fastapi/crud"
	"github.com/shoppehub/fastapi/engine"
)

func InitApi(resource *crud.Resource, r *gin.Engine) {

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/collection", func(c *gin.Context) {
			collection.CreateCollection(resource, c)
		})
		apiv1.GET("/collection/:id", func(c *gin.Context) {
			collection.GetCollection(resource, c)
		})
		apiv1.POST("/collections", func(c *gin.Context) {
			collection.QueryCollection(resource, c)
		})
	}

	apicol := r.Group("/api/collection")
	{
		apicol.GET("/:group/:collection/:id", func(c *gin.Context) {
			engine.GetWithId(resource, c)
		})
		apicol.POST("/:group/:collection/findone", func(c *gin.Context) {
			engine.FindOne(resource, c)
		})
		apicol.POST("/:group/:collection/", func(c *gin.Context) {
			engine.Post(resource, c)
		})

		apicol.POST("/delete/:group/:collection/:id", func(c *gin.Context) {
			engine.DeleteId(resource, c)
		})

	}

}
