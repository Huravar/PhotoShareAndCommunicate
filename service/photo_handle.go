package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"photo_service/gadgets"
	"photo_service/model"
	"photo_service/picture_handle"
	"photo_service/user_static_info"
)

// UploadPhoto
// @Summary      用户上传图片接口
// @Description  接收用户上传的图片文件并存储到指定用户目录，同时记录文件元信息到数据库
// @Tags         Photo control
// @Param        Authorization header string true "Bearer 用户令牌"
// @Param        photo formData file true "要上传的图片文件(支持JPG/PNG)"
// @Param 		 id    formData string true "用户id"
// @Param        description formData string false "图片描述信息"
// @Success      200  {string} json{"code","message"}
// @Success      400  {string} json{"code","message"}
// @Failure      401  {string} json{"code","message"}
// @Failure      500  {string} json{"code","message"}
// @Router       /api/photo/upload [post]
func UploadPhoto(c *gin.Context) {
	IUserTokenBasicInfo, err := VerifyToken(c)
	if err != nil {
		return
	}
	IPhoto, IPhotoType, err := picture_handle.CommonPhotoDeal(c, "photo", "./user_static_info/"+IUserTokenBasicInfo.UserId+"/pictures")
	if err != nil {
		return
	}
	IfilePath := "./user_static_info/" + IUserTokenBasicInfo.UserId + "/pictures/" + generateFilename(IPhoto.Filename)
	_, IPhotoNum := model.FindUserPhotoInfoByPhotoName(gadgets.StringToUint(IUserTokenBasicInfo.UserId), IPhoto.Filename)
	if IPhotoNum == 1 {
		log.Println(IPhoto.Filename, "图片已创建不可重复创建！")
		c.JSON(http.StatusBadRequest, gin.H{"code": 3, "message": "图片已创建不可重复创建！"})
		return
	}
	model.CreatUserPhotoInfo(model.UserPhotoInfo{UserID: gadgets.StringToUint(IUserTokenBasicInfo.UserId),
		UserName: IUserTokenBasicInfo.UserName, OriginalName: IPhoto.Filename, StoragePath: IfilePath, FileSize: IPhoto.Size,
		FileType: IPhotoType, Description: c.PostForm("description")})
	if err = user_static_info.AddFileByGin(c, IPhoto, IfilePath); err != nil {
		log.Println("保存图片失败", err)
		c.JSON(500, gin.H{"code": 2, "message": "保存图片失败"})
		err := model.DeleteUserPhotoInfoByPhotoName(IPhoto.Filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	c.JSON(200, gin.H{"code": 4, "message": "上传文件成功！"})
	return
}
