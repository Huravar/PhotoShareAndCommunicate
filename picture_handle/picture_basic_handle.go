package picture_handle

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"photo_service/user_static_info"
	"strings"
)

func IsImage(PImageFile *multipart.FileHeader) (string, bool) {
	IImage, err := PImageFile.Open()
	if err != nil {
		return "", false
	}
	defer IImage.Close()
	buffer := make([]byte, 512)
	if _, err := IImage.Read(buffer); err != nil {
		return "", false
	}
	mimeType := http.DetectContentType(buffer)
	return mimeType, strings.HasPrefix(mimeType, "image/")
}

func CommonPhotoDeal(c *gin.Context, PRequestParam, PFilePath string) (*multipart.FileHeader, string, error) {
	Iphoto, err := c.FormFile(PRequestParam)
	if err != nil {
		log.Println("图片文件为空！", err)
		c.JSON(400, gin.H{"code": -1, "message": "图片文件为空！"})
		return Iphoto, "", fmt.Errorf("图片文件为空！")
	}
	IPhotoType, err1 := IsImage(Iphoto)
	if !err1 {
		log.Println("仅支持JPEG/PNG图片")
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "仅支持JPEG/PNG图片"})
		return Iphoto, IPhotoType, fmt.Errorf("仅支持JPEG/PNG图片")
	}
	err = user_static_info.AddFileDir(PFilePath)
	if err != nil {
		log.Println("创建文件夹失败！", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "message": "创建文件夹失败！"})
		return Iphoto, IPhotoType, fmt.Errorf("创建文件夹失败！")
	}
	return Iphoto, IPhotoType, nil
}
