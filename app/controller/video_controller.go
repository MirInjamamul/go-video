package controller

import (
	"log"
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
		// Log for Form Values
		// log.Printf("Request Headers: %v \n", c.Request.Header)
		// Log for Form Values
		// formValues := c.Request.PostForm
		// log.Printf("Request Form Values: %v \n", formValues)
		log.Printf(err.Error())

		c.JSON(400, gin.H{"status": false, "error": "Bad Request - No Video file uploaded"})
		return
	}

	videoPaths, err, filename := vc.videoModel.SaveVideo(c, file)

	if err != nil {
		log.Printf(err.Error())
		c.JSON(500, gin.H{"status": false, "error": "Internal Server Error - Failed to save Video"})
	} else {
		c.JSON(200, gin.H{"status": true, "message": "Video File Upload Successfully", "url": videoPaths})
	}

	// Process the uploaded video

	go func() {
		videoPaths, err = vc.videoModel.ProcessVideo(videoPaths["original"], filename)
		if err != nil {
			log.Printf((err.Error()))
		} else {
			log.Printf("Process Done")
		}
	}()

}
