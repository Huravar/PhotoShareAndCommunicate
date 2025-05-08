package router

import (
	"photo_service/docs"
	"photo_service/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.POST("/user/CreateUser", service.CreateUser)
	r.POST("/user/LoginInUser", service.LoginInUser)
	r.POST("/user/upload-avatar", service.UploadAvatars)
	r.POST("/user/upload_self-introduce", service.UploadSelfIntroduce)
	r.POST("/user/upload_user-email", service.UploadUserEmail)
	r.POST("/user/upload_user-phone", service.UploadUserPhone)
	r.POST("/api/photo/upload", service.UploadPhoto)
	r.POST("/user/download_user-homepage-message", service.SearchUserHomePageInfo)
	r.POST("/user/download_user-basic-message", service.SearchUserBasicInfo)
	r.POST("/api/photo/delete", service.DeleteUserPhoto)
	r.GET("/ws", service.UserBasicCommunicate)
	return r
}
