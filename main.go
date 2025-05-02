package main

import (
	"photo_service/model"
	"photo_service/router"
	"photo_service/utils"
)

// @title          APIphoto_server
// @version        V1.0

func main() {
	utils.ViperInitialization() //初始化viper
	utils.OpenDatabase()
	utils.OpenRedis()
	model.TableInit()
	r := router.Router()
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
