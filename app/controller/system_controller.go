package controller

import (
	"context"
	"fmt"

	"video-server/app/model"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type SystemController struct {
	systemModel *model.System
}

func NewSystemController() *SystemController {
	return &SystemController{
		systemModel: &model.System{},
	}
}

func (sc *SystemController) ContainerStatus(c *gin.Context) {

	// Create a new Docker Client
	cli, err := client.NewClientWithOpts(client.WithHost("http://169.150.255.33:2375"))
	if err != nil {
		c.String(400, err.Error())
	}

	// Specify the Container ID or Name you want to check
	containerID := "469433b71b77"

	// Get the container's Details
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		fmt.Printf("Error Inspecting Container: %s \n", err)
		c.String(400, err.Error())
		return
	}

	// Check the container State
	if containerInfo.State.Running {
		c.String(200, "Container is Running")
	} else {
		c.String(400, "Container is not Running")
	}
}
