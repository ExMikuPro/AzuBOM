package Api

import (
	"AzuBOM/Message"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	addApi := router.Group("/v1/")

	addApi.GET("/type", Message.GetType)
	addApi.POST("/type", Message.AddType)
	addApi.PUT("/type", Message.UpdateType)
	addApi.DELETE("/type", Message.DeleteType)

	addApi.GET("/component", Message.GetComponent)
	addApi.POST("/component", Message.AddComponent)
	addApi.PUT("/component", Message.UpdateComponent)
	addApi.DELETE("/component", Message.DeleteComponent)

}
