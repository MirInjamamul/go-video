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

	// Serve static files from the "uploads" directory
	r.Static("/uploads", "./uploads")

	// Register the route

	view.RegisterRoutes(r)

	r.Run("51.159.111.59:8090")
	// r.Run(":8090")
}
