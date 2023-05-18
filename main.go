package main

import (
	"log"
	"os"

	"video-server/app/view"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Fatal("Failed to Create Upload Directory")
	}
	// Register the route

	view.RegisterRoutes(r)

	r.Run(":8090")
}
