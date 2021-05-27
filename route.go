package fastapi

import "github.com/gin-gonic/gin"

func InitApi(resource *Resource, r *gin.Engine) {
	dBResource = resource

	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/collection", CreateCollection)
		apiv1.GET("/collection/:id", GetCollection)
		apiv1.POST("/collections", QueryCollection)
	}

	apicol := r.Group("/api/collection")
	{
		apicol.GET("/:group/:collection/:id", GetWithId)
	}

}
