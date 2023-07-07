package controller

import (
	"log"
	"net/http"
	"video-server/app/model"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
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

func GetStackStatus(c *gin.Context) {
	stackName := c.Param("name")

	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, err)
		return
	}

	stack, _, err := cli.StackInspectWithRaw(c.Request.Context(), stackName)

	if err != nil {
		if client.IsErrNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack Not Found"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}

		return
	}

	status := "Stopped"
	if stack.State.status == types.ActiveStatus {
		status = "Running"
	}

	c.JSON(http.StatusOK, gin.H{"status": status})

}
