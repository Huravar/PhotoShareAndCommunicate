package picture_handle

import (
	"mime/multipart"
	"net/http"
	"strings"
)

func IsImage(PImageFile *multipart.FileHeader) bool {
	IImage, err := PImageFile.Open()
	if err != nil {
		return false
	}
	defer IImage.Close()
	buffer := make([]byte, 512)
	if _, err := IImage.Read(buffer); err != nil {
		return false
	}
	mimeType := http.DetectContentType(buffer)
	return strings.HasPrefix(mimeType, "image/")
}
