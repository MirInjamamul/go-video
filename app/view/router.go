package view

import (
	"video-server/app/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	videoController := controller.NewVideoController()
	systemController := controller.NewSystemController()

	r.POST("/uploadVideo", videoController.UploadVideo)
	r.GET("/containerStatus", systemController.ContainerStatus)
}
