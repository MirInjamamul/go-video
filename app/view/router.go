package view

import (
	"video-server/app/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	videoController := controller.NewVideoController()

	r.POST("/uploadVideo", videoController.UploadVideo)
	r.GET("/stack/:name", videoController.GetStackStatus)
}
