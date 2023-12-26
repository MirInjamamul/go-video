package controller

import (
	"log"
	"strconv"
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
	userId := c.PostForm("userId")

	if err != nil {
		// Log for Form Values
		// log.Printf("Request Headers: %v \n", c.Request.Header)
		// Log for Form Values
		// formValues := c.Request.PostForm
		// log.Printf("Request Form Values: %v \n", formValues)
		log.Printf(err.Error())

		c.JSON(400,
			gin.H{
				"status": false,
				"error":  "Bad Request - No Video file uploaded"})
		return
	}

	log.Printf("User Id %s", userId)

	videoPaths, err, filename := vc.videoModel.SaveVideo(c, file, userId)

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

			// Hit the Update API with VideoPath
			if err1 := vc.videoModel.UpdateVideoPaths(videoPaths); err1 != nil {
				log.Printf("Failed to update video Paths: %v", err)
			} else {
				log.Printf("video Path Updated Successfully")
			}
		}
	}()

}

func (vc *VideoController) UploadSDKFile(c *gin.Context) {
	var videoPath string
	var err error
	var filename string

	file, err := c.FormFile("file")
	fileTypeStr := c.PostForm("fileType")

	if err != nil {
		// Log for Form Values
		// log.Printf("Request Headers: %v \n", c.Request.Header)
		// Log for Form Values
		// formValues := c.Request.PostForm
		// log.Printf("Request Form Values: %v \n", formValues)

		c.JSON(400,
			gin.H{
				"status": false,
				"error":  "Bad Request - No Video file uploaded"})
		return
	}

	fileType, err1 := strconv.Atoi(fileTypeStr)
	if err1 != nil {
		c.JSON(400,
			gin.H{
				"status": false,
				"error":  "Invalid file type provided"})
	}

	switch fileType {
	case 1:
		/// Video File
		videoPath, err, filename = vc.videoModel.SaveSDKVideo(c, file)

	}

	if err != nil {
		log.Printf(err.Error())
		c.JSON(500, gin.H{"status": false, "error": "Internal Server Error - Failed to save Video"})
	} else {
		c.JSON(200, gin.H{"status": true, "message": "Video File Upload Successfully", "url": videoPath})
	}

	// Process the uploaded video

	go func() {
		err = vc.videoModel.ProcessSDKVideo(videoPath, filename)
		if err != nil {
			log.Printf((err.Error()))
		} else {
			log.Printf("Process Done")
			// Hit the Update API with VideoPath
			// if err1 := vc.videoModel.UpdateVideoPaths(videoPaths); err1 != nil {
			// 	log.Printf("Failed to update video Paths: %v", err)
			// } else {
			// 	log.Printf("video Path Updated Successfully")
			// }
		}

	}()

}

func (vc *VideoController) UploadChatFile(c *gin.Context) {
	file, err := c.FormFile("file")
	userId := c.PostForm("userId")

	if err != nil {
		// Log for Form Values
		// log.Printf("Request Headers: %v \n", c.Request.Header)
		// Log for Form Values
		// formValues := c.Request.PostForm
		// log.Printf("Request Form Values: %v \n", formValues)
		log.Printf(err.Error())

		c.JSON(400,
			gin.H{
				"status": false,
				"error":  "Bad Request - No file uploaded"})
		return
	}

	videoPaths, err := vc.videoModel.SaveChatFile(c, file, userId)

	if err != nil {
		log.Printf(err.Error())
		c.JSON(500, gin.H{"status": false, "error": "Internal Server Error - Failed to save Video"})
	} else {
		c.JSON(200, gin.H{
			"status":  true,
			"message": "Video File Upload Successfully",
			"url":     videoPaths})
	}
}
