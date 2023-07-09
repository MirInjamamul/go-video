package controller

import (
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
		c.String(400, "Bad Request - No Video file uploaded")
		return
	}

	err = vc.videoModel.SaveVideo(c, file)

	if err != nil {
		c.String(500, err.Error())
	} else {
		c.String(200, "Video File Upload Successfully")
	}

}
