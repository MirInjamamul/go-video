package view

import (
	"video-server/app/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	videoController := controller.NewVideoController()
	imageController := controller.NewImageController()

	r.POST("/uploadVideo", videoController.UploadVideo)
	r.POST("/uploadSDKFile", videoController.UploadSDKFile)
	r.POST("/uploadChatFile", videoController.UploadChatFile)

	r.POST("/faceAuth", imageController.Authenticate)
}
