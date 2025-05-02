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
	
	return r
}
