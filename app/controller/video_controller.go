package controller

import (
	"fmt"
	"video-server/app/model"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	videoModel *model.Video
}

func NewVideoController() *VideoController {
	return &VideoController{
		videoModel: &model.Video{},
	}
}

func (vc *VideoController) UploadVideo(c *gin.Context) {
	file, err := c.FormFile("video")

	if err != nil {
		c.JSON(400, gin.H{"status": false, "error": "Bad Request - No Video file uploaded"})
		return
	}

	err = vc.videoModel.SaveVideo(c, file)

	if err != nil {
		fmt.Sprintf(err.Error())
		c.JSON(500, gin.H{"status": false, "error": "Internal Server Error - Failed to save Video"})
	} else {
		c.JSON(200, gin.H{"status": true, "message": "Video File Upload Successfully"})
	}

}
