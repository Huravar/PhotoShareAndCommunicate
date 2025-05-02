package user_static_info

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
)

func AddFileDir(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}
func DeleteFile(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

func AddFileByByte(path string, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func AddFileByGin(c *gin.Context, file *multipart.FileHeader, path string) error {
	err := c.SaveUploadedFile(file, path)
	if err != nil {
		return err
	}
	return nil
}
