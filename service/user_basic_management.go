package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"photo_service/crypt"
	"photo_service/model"
	"photo_service/utils"
	"strconv"
	"time"
)

// CreateUser
// @Summary add user
// @Tags user management
// @param UserName formData string false "UserName"
// @param Password formData string flase "Password"
// @param RePassword formData string flase "RePassword"
// @Success 200 {string} json{"code","message"}
// @Failure 400 {string} json{"code","message"}
// @Router /user/CreateUser [post]
func CreateUser(c *gin.Context) {

	Name, NameErr := c.GetPostForm("UserName")
	Password, PasswordErr := c.GetPostForm("Password")
	RePassword, RePasswordErr := c.GetPostForm("RePassword")
	if !NameErr || !PasswordErr || !RePasswordErr {
		c.JSON(200, gin.H{
			"code":    -1, // -1表示元素缺少
			"message": "用户名或密码或重复密码不能为空！",
		})
		return
	}
	if Password != RePassword {
		c.JSON(200, gin.H{
			"code":    0, //0表示密码与重复密码不同
			"message": "密码与重复密码不同！",
		})
		return
	}
	_, EffNum := model.FindUserByName(Name)
	fmt.Println(EffNum)
	if EffNum == 1 {
		c.JSON(200, gin.H{
			"code":    -2, //-2表示用户名以存在
			"message": "用户名重复请更换用户名！",
		})
		return
	}
	CryPassword, _ := crypt.HashPassword(Password)
	user := model.BasicUserInformation{
		UserName: Name,
		PassWord: CryPassword,
		Identity: uint(0),
	}
	model.CreatUserBasicInfo(&user)
	c.JSON(200, gin.H{
		"code":    1, // 表示用户创建成功
		"message": "创建用户成功",
	})
	return
}

// LoginInUser
// @Summary login in user
// @Tags user management
// @param UserName formData string false "UserName"
// @param Password formData string flase "Password"
// @Success 200 {string} json{"code","message","id","token"}
// @Failure 400 {string} json{"code","message","id","token"}
// @Router /user/LoginInUser [post]
func LoginInUser(c *gin.Context) {
	UserName, NameErr := c.GetPostForm("UserName")
	Password, PasswordErr := c.GetPostForm("Password")
	if !NameErr || !PasswordErr {
		c.JSON(200, gin.H{
			"code":    -1, //-1 用户名或密码为空
			"message": "用户名或密码为空！",
			"id":      "",
			"token":   "",
		})
		return
	}
	user, EffNum := model.FindUserByName(UserName)
	if EffNum < 1 {
		c.JSON(200, gin.H{
			"code":    -1, //用户不存在
			"message": "用户不存在请注册！",
			"id":      "",
			"token":   "",
		})
		return
	}
	if !crypt.CheckPasswordHash(Password, user.PassWord) {
		c.JSON(200, gin.H{
			"code":    0, //密码错误
			"message": "密码错误！",
			"id":      "",
			"token":   "",
		})
		return
	}
	UserToken, TokenKey, err := crypt.GenerateToken(user)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -5, //服务错误
			"message": "服务错误！",
			"id":      "",
			"token":   "",
		})
		return
	}
	err = CrateOrUpdateNetworkInfoForLogin(user, TokenKey, c)
	if err != nil {
		return
	}
	utils.Red.Set(context.Background(), strconv.Itoa(int(user.ID)), TokenKey, 6*time.Hour)
	c.JSON(200, gin.H{
		"code":    1, //登录成功
		"message": "登录成功！",
		"id":      user.ID,
		"token":   UserToken,
	})
	return
}

func CrateOrUpdateNetworkInfoForLogin(User model.BasicUserInformation, TokenKey string, c *gin.Context) error {
	ItemUserNetworkInfo, ItemEffNum := model.FindUserNetworkById(User.ID)
	fmt.Println(c.ClientIP())
	remoteAddr := c.Request.RemoteAddr
	ip, port, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Printf("Failed to parse RemoteAddr: %v", err)
		return err
	}
	if ItemEffNum < 1 {
		Item2UserNetworkInfo := model.UserNetwork{
			UserID:        User.ID,
			UserName:      User.UserName,
			ClientIp:      ip,
			ClientPort:    port,
			UserTK:        TokenKey,
			LoginTime:     model.TimePointer(),
			HeartbeatTime: model.TimePointer(),
			IsLogin:       true,
			DeviceInfo:    c.GetHeader("User-Agent"),
		}
		model.CreateUserNetwork(Item2UserNetworkInfo)
		return nil
	} else {
		model.UpdateNetworkForDeviceInfo(ItemUserNetworkInfo, c.GetHeader("User-Agent"))
		model.UpdateNetworkForClintNAdress(ItemUserNetworkInfo, ip, port)
		model.UpdateNetworkForTk(ItemUserNetworkInfo, TokenKey)
		model.UpdateNetworkForIsLogout(ItemUserNetworkInfo, true)
		model.UpdateNetworkForLoginTime(ItemUserNetworkInfo)
		model.UpdateNetworkForHeartbeatTime(ItemUserNetworkInfo)
		return nil
	}
}

func VerifyToken(c *gin.Context) (crypt.UserTokenBasicInfo, error) {
	ItemUserId, err := c.GetPostForm("id")
	if !err {
		log.Println("获取用户id失败", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -2, "message": "获取用户id失败",
		})
		return crypt.UserTokenBasicInfo{}, errors.New("获取用户id失败")
	}
	ItemTokenKey, err2 := utils.Red.Get(context.Background(), ItemUserId).Result()
	if err2 != nil {
		log.Println(ItemUserId, "登录时间过期请重新登录", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -3, "message": "登录时间过期请重新登录",
		})
		return crypt.UserTokenBasicInfo{}, errors.New("登录时间过期请重新登录")
	}
	ItemUserTokenBasicInfo, err3 := crypt.ParasedAndVerify(c.GetHeader("Authorization"), ItemTokenKey)
	if err3 != nil {
		log.Println(ItemUserId, "token验证失败请重新登录或检查token值", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": -4, "message": "token验证失败请重新登录或检查token值",
		})
		return crypt.UserTokenBasicInfo{}, errors.New("token验证失败请重新登录或检查token值")
	}
	return ItemUserTokenBasicInfo, nil
}
