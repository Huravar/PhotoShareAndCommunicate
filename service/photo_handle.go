package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
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
		err := model.DeleteUserPhotoInfoByPhotoName(gadgets.StringToUint(IUserTokenBasicInfo.UserId), IPhoto.Filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	c.JSON(200, gin.H{"code": 4, "message": "上传文件成功！"})
	return
}

// DeleteUserPhoto
// @Summary 删除用户上传的图片
// @Description 删除用户上传的图片
// @Tags Photo control
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id     formData string true "用户id"
// @Param PhotoName formData string true "图片名称"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/photo/delete [post]
func DeleteUserPhoto(c *gin.Context) {
	IUserTokenBasicInfo, err := VerifyToken(c)
	if err != nil {
		return
	}
	IPhotoName, err1 := c.GetPostForm("PhotoName")
	if !err1 {
		log.Println("没有上传必要参数！")
		c.JSON(400, gin.H{"code": -1, "message": "没有上传必要参数！"})
		return
	}
	IUserPhotoInfo, INum := model.FindUserPhotoInfoByPhotoName(gadgets.StringToUint(IUserTokenBasicInfo.UserId), IPhotoName)
	if INum == 0 {
		log.Println("没有找到文件", IPhotoName)
		c.JSON(400, gin.H{"code": 0, "message": "没有找到文件！"})
		return
	}
	_, err = os.Stat(IUserPhotoInfo.StoragePath)
	if err != nil {
		log.Println("文件不存在！", IUserPhotoInfo.StoragePath)
		c.JSON(400, gin.H{"code": -2, "message": "文件不存在！"})
		return
	}
	err = model.DeleteUserPhotoInfoByPhotoName(gadgets.StringToUint(IUserTokenBasicInfo.UserId), IPhotoName)
	if err != nil {
		log.Println("数据库图片删除失败！", IPhotoName)
		c.JSON(http.StatusInternalServerError, gin.H{"code": -3, "message": "数据库图片删除失败！"})
		return
	}
	err = os.Remove(IUserPhotoInfo.StoragePath)
	if err != nil {
		log.Println("图片删除失败！", IUserPhotoInfo.StoragePath)
		c.JSON(http.StatusInternalServerError, gin.H{"code": -4, "message": "图片删除失败！"})
		return
	}
	c.JSON(200, gin.H{"code": 1, "message": "图片删除成功！"})
	return
}
