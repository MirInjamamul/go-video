package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"video-server/app/config"

	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Video struct {
	Filename string
}

type VideoPaths struct {
	OldPath string `json:"oldPath"`
	NewPath string `json:"newPath"`
}

func (v *Video) SaveVideo(c *gin.Context, file *multipart.FileHeader, userId string) (map[string]string, error, string) {
	// Generating new FileName
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	v.Filename = fmt.Sprintf("%s_%d%s", userId, timestamp, filepath.Ext(file.Filename))

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

func (v *Video) SaveSDKVideo(c *gin.Context, file *multipart.FileHeader) (string, error, string) {
	// Generating new FileName
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	v.Filename = fmt.Sprintf("%d", timestamp)

	videopath := ""

	// Extract the filename

	filename := strings.TrimSuffix(v.Filename, filepath.Ext(v.Filename))

	// Create a Directory with the same name as the uploaded file
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err, ""
	}

	destination := filepath.Join(uploadDir, v.Filename)

	log.Printf("Video Uploaded Started")
	if err := c.SaveUploadedFile(file, destination); err != nil {
		return "", err, ""
	}

	log.Printf("Video Uploaded Finished")

	// Save the original video path
	videopath = destination

	return videopath, nil, filename
}

func (v *Video) ProcessVideo(sourceFile string, filename string) (map[string]string, error) {
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
		cmd := exec.Command("ffmpeg", "-i", sourceFile, "-vf", fmt.Sprintf("scale=-2:%s", res), outputPath)
		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("FFmpeg conversion error: %s", err.Error())
		}

		// Save the resolutions path in the videoPaths map
		videopaths[res] = outputPath
	}

	if err := os.Rename(videopaths["360p"], videopaths["original"]); err != nil {
		return nil, fmt.Errorf("error replacing input file with output file")
	}

	log.Printf("Video Process Finished")

	videopaths["original"] = videopaths["360p"]

	return videopaths, nil
}

func (v *Video) ProcessSDKVideo(sourceFile string, filename string) error {
	log.Printf("Video Process Started")
	// Create a Directory with the same name as the uploaded file
	uploadDir := "uploads"

	// Convert and save the video in different resolutions using FFmpeg
	// resolutions := []string{"240p", "360p", "480p"}
	resolutions := []string{"360p"}
	outputPath := ""

	for _, res := range resolutions {

		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return err
		}

		outputFilename := fmt.Sprintf("%s_%s.mp4", filename, res)
		outputPath = filepath.Join(uploadDir, outputFilename)

		// Run FFmpeg command to convert video to the desired resolution
		cmd := exec.Command("ffmpeg", "-i", sourceFile, "-vf", fmt.Sprintf("scale=-2:%s", res), outputPath)
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("FFmpeg conversion error: %s", err.Error())
		}
	}

	if err := os.Remove(sourceFile); err != nil {
		return fmt.Errorf("error removing original sourceFile")
	}

	if err := os.Rename(outputPath, sourceFile); err != nil {
		return fmt.Errorf("error removing original sourceFile")
	}

	log.Printf("Video Process Finished")

	return nil
}

func (v *Video) UpdateVideoPaths(videoPaths map[string]string) error {
	apiURL := config.APIBaseURL

	// Create the request body from video Paths data
	paths := VideoPaths{
		OldPath: videoPaths["original"],
		NewPath: videoPaths["360p"],
	}

	// Marshal the data structure to JSON
	payload, err := json.Marshal(paths)

	if err != nil {
		return err
	}

	// create a http request with post method
	req, err := http.NewRequest("POST", apiURL+config.PostUpdate, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Set the header
	req.Header.Set("Content-Type", "application/json")

	// // Perform Http request
	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	// defer response.Body.Close()

	// // Check the response status
	if response.StatusCode == 200 {
		return nil
	} else {

		responseByrtes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		errorMessage := string(responseByrtes)
		log.Printf(errorMessage)
	}

	return nil
}

func (v *Video) SaveChatFile(c *gin.Context, file *multipart.FileHeader, userId string) (map[string]string, error) {
	// Generating new FileName
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	v.Filename = fmt.Sprintf("%s_%d%s", userId, timestamp, filepath.Ext(file.Filename))

	videopaths := make(map[string]string)

	// Extract the filename

	filename := strings.TrimSuffix(v.Filename, filepath.Ext(v.Filename))

	// Create a Directory with the same name as the uploaded file
	uploadDir := filepath.Join("uploads", filename)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, err
	}

	destination := filepath.Join(uploadDir, v.Filename)

	log.Printf("File Uploaded Started")
	if err := c.SaveUploadedFile(file, destination); err != nil {
		return nil, err
	}

	log.Printf("File Uploaded Finished")

	// Save the original video path
	videopaths["original"] = destination

	return videopaths, nil
}
