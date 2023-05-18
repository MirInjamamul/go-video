package model

import (
	"mime/multipart"
	"os"
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

	return nil
}
