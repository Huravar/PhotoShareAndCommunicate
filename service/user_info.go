package service

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"photo_service/crypt"
	"photo_service/model"
	"photo_service/picture_handle"
	"photo_service/user_static_info"
	"strconv"
	"time"
)

// UploadAvatars
// @Summary 上传用户头像
// @Description 上传用户头像接口（支持JPEG/PNG格式，最大30MB）
// @Tags User Home Message
// @Param avatar formData file true "头像文件（支持JPEG/PNG格式）"
// @Param id     formData string true "用户id"
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json{"code","message"}
// @Failure 400 {string} json{"code","message"}
// @Failure 413 {string} json{"code","message"}
// @Failure 500 {string} json{"code","message"}
// @Router /user/upload-avatar [post]
func UploadAvatars(c *gin.Context) {
	ItemUserTokenBasicInfo, err := VerifyToken(c)
	if err != nil {
		return
	}
	const MaxSIze = 30 << 20
	if err := c.Request.ParseMultipartForm(MaxSIze); err != nil {
		log.Println("图片尺寸过大！", err)
		c.JSON(400, gin.H{"code": -1, //图片尺寸过大！
			"message": "图片尺寸过大！"})
	}
	ItemAvatar, err := c.FormFile("avatar")
	if err != nil {
		log.Println("未找到头像文件", err)
		c.JSON(http.StatusBadRequest, gin.H{"code": 0, "message": "未找到头像文件"})
		return
	}

	if !picture_handle.IsImage(ItemAvatar) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "message": "仅支持JPEG/PNG图片"})
		return
	}
	err1 := user_static_info.AddFileDir("./user_static_info/" + ItemUserTokenBasicInfo.UserId + "/avatars")
	IfilePath := "./user_static_info/" + ItemUserTokenBasicInfo.UserId + "/avatars/" + generateFilename(ItemAvatar.Filename)
	UpdataOrCreateAvatar(ItemUserTokenBasicInfo, IfilePath)
	if err := user_static_info.AddFileByGin(c, ItemAvatar, IfilePath); err != nil || err1 != nil {
		log.Println("保存头像失败", err)
		c.JSON(500, gin.H{"code": 2, "message": "保存头像失败"})
	}
	c.JSON(200, gin.H{"code": 3, "message": "设置头像成功！"})

}
func generateFilename(original string) string {
	ext := filepath.Ext(original)
	return fmt.Sprintf("%d_%x%s",
		time.Now().UnixNano(),
		md5.Sum([]byte(original)),
		ext)
}

func UpdataOrCreateAvatar(PUserTokenBasicInfo crypt.UserTokenBasicInfo, PAvaDir string) {
	IUserId, _ := strconv.Atoi(PUserTokenBasicInfo.UserId)
	IUserHomePageInfo, Num := model.FindUserHomePageInfoByUserId(uint(IUserId))
	if Num == 0 {
		model.CreateUserHomePageInfo(model.UserHomePageInfo{
			UserID:     uint(IUserId),
			UserName:   PUserTokenBasicInfo.UserName,
			AvatarPath: PAvaDir,
		})
		return
	}
	if Num == 1 {
		err := user_static_info.DeleteFile(IUserHomePageInfo.AvatarPath)
		if err != nil {
			log.Println("删除旧头像失败", err)
			return
		}
		err = model.UpdateAvaPathById(uint(IUserId), PAvaDir)
		if err != nil {
			log.Println("头像路径更新错误", err)
			return
		}

	}

}

// UploadSelfIntroduce
// @Summary 更新用户自我介绍
// @Description 用户更新自我介绍信息接口（需Token认证）
// @Tags User Home Message
// @Param selfIntroduce formData string true "用户自我介绍内容"
// @Param id     formData string true "用户id"
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {string} json{"code", "id","message"}
// @Failure 400 {string} json{"code", "id","message"}
// @Failure 500 {string} json{"code", "id","message"}
// @Router /user/upload_self-introduce [post]
func UploadSelfIntroduce(c *gin.Context) {
	ItemUserTokenBasicInfo, err := VerifyToken(c)
	if err != nil {
		return
	}
	ISelfIntroduce, err1 := c.GetPostForm("selfIntroduce")
	if !err1 {
		log.Println("缺少请求参数")
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "id": ItemUserTokenBasicInfo.UserId, "message": "缺少请求参数"})
		return
	}
	err = UpdateOrCreateSelfIntroduce(ItemUserTokenBasicInfo, ISelfIntroduce)
	if err != nil {
		log.Println("<UNK>", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 0, "id": ItemUserTokenBasicInfo.UserId, "message": "服务器错误！"})
		return
	}
	c.JSON(200, gin.H{"code": 1, "id": ItemUserTokenBasicInfo.UserId, "message": "更新自我介绍成功！"})
	return
}

func UpdateOrCreateSelfIntroduce(PUserTokenBasicInfo crypt.UserTokenBasicInfo, PSelfIntroduce string) error {
	IUserId, _ := strconv.Atoi(PUserTokenBasicInfo.UserId)
	_, Num := model.FindUserHomePageInfoByUserId(uint(IUserId))
	if Num == 0 {
		model.CreateUserHomePageInfo(model.UserHomePageInfo{
			UserID:        uint(IUserId),
			UserName:      PUserTokenBasicInfo.UserName,
			SelfIntroduce: PSelfIntroduce,
		})
		return nil
	}
	if Num == 1 {
		err := model.UpdateSlfIntroduceById(uint(IUserId), PSelfIntroduce)
		if err != nil {
			log.Println("更新数据库自我介绍失败", err)
			return err
		}
	}
	return nil
}
