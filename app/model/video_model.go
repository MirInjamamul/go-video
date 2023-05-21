package model

import (
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Video struct {
	Filename string
}

func (v *Video) SaveVideo(c *gin.Context, file *multipart.FileHeader) error {
	v.Filename = filepath.Base(file.Filename)

	// Extract the filename

	filename := strings.TrimSuffix(v.Filename, filepath.Ext(v.Filename))

	// Create a Directory with the same name as the uploaded file

	uploadDir := filepath.Join("uploads", filename)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return err
	}

	destination := filepath.Join(uploadDir, v.Filename)

	if err := c.SaveUploadedFile(file, destination); err != nil {
		return err
	}

	// Convert and save the video in different resolutions using FFmpeg
	resolutions := []string{"240p", "360p", "480p"}
	for _, res := range resolutions {
		resolutionDir := filepath.Join(uploadDir, res)
		if err := os.MkdirAll(resolutionDir, os.ModePerm); err != nil {
			return err
		}

		outputFilename := fmt.Sprintf("%s_%s.mp4", filename, res)
		outputPath := filepath.Join(resolutionDir, outputFilename)

		// Run FFmpeg command to convert video to the desired resolution
		cmd := exec.Command("ffmpeg", "-i", destination, "-vf", fmt.Sprintf("scale=-2:%s", res), outputPath)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("FFmpeg conversion error: %s", err.Error())
		}
	}

	return nil
}
