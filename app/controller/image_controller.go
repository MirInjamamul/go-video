package controller

import (
	"net/http"
	"video-server/app/model"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	imageModel *model.Images
}

func NewImageController() *ImageController {
	return &ImageController{
		imageModel: &model.Images{},
	}
}

func (ic *ImageController) Authenticate(c *gin.Context) {
	image1, err1 := c.FormFile("image1")
	image2, err2 := c.FormFile("image2")

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"status": false,
				"error":  "ImageFile Missing"})
		return
	}

	if image1.Filename == image2.Filename {
		c.JSON(http.StatusOK,
			gin.H{
				"status": true,
				"result": "Authenticated",
			})
	} else {
		c.JSON(http.StatusUnauthorized,
			gin.H{
				"status": true,
				"result": "Authenticated",
			})
	}
}
