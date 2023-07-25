package model

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Video struct {
	Filename string
}

func (v *Video) SaveVideo(c *gin.Context, file *multipart.FileHeader, userId string) (map[string]string, error, string) {
	// Generating new FileName
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	v.Filename = fmt.Sprintf("%s_%d_%s", userId, timestamp, filepath.Ext(file.Filename))

	videopaths := make(map[string]string)

	// Extract the filename

	filename := strings.TrimSuffix(v.Filename, filepath.Ext(v.Filename))

	// Create a Directory with the same name as the uploaded file
	uploadDir := filepath.Join("uploads", filename)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err, ""
	}

	destination := filepath.Join(uploadDir, v.Filename)

	log.Printf("Video Uploaded Started")
	if err := c.SaveUploadedFile(file, destination); err != nil {
		return nil, err, ""
	}

	log.Printf("Video Uploaded Finished")

	// Save the original video path
	videopaths["original"] = destination
	videopaths["360p"] = ""

	return videopaths, nil, filename
}

func (v *Video) ProcessVideo(destination string, filename string) (map[string]string, error) {
	log.Printf("Video Process Started")
	// Create a Directory with the same name as the uploaded file
	uploadDir := filepath.Join("uploads", filename)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	// Convert and save the video in different resolutions using FFmpeg
	// resolutions := []string{"240p", "360p", "480p"}
	resolutions := []string{"360p"}
	videopaths := make(map[string]string)

	for _, res := range resolutions {
		resolutionDir := filepath.Join(uploadDir, res)
		if err := os.MkdirAll(resolutionDir, os.ModePerm); err != nil {
			return nil, err
		}

		outputFilename := fmt.Sprintf("%s_%s.mp4", filename, res)
		outputPath := filepath.Join(resolutionDir, outputFilename)

		// Run FFmpeg command to convert video to the desired resolution
		cmd := exec.Command("ffmpeg", "-i", destination, "-vf", fmt.Sprintf("scale=-2:%s", res), outputPath)
		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("FFmpeg conversion error: %s", err.Error())
		}

		// Save the resolutions path in the videoPaths map
		videopaths[res] = outputPath
	}

	log.Printf("Video Process Finished")

	return videopaths, nil
}
