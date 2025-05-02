package sql

import (
	"fmt"
	"photo_service/model"
	"photo_service/utils"
	"reflect"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SqlText() {
	utils.ViperInitialization()
	fmt.Println(viper.Get("mysql"))

	db, err := gorm.Open(mysql.Open(DatabaseLoginString()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.BasicUserInformation{})
	// var ItemExample []model.BasicUserInformation
	// result := db.Find(&ItemExample)
	// a := result.RowsAffected
	// fmt.Printf("1111111%v\n", a)
	// fmt.Println("调试输出完成")
	// for _, value := range ItemExample {
	// 	fmt.Println(value)
	// }
	// db.Model(&model.BasicUserInformation{}).Where("Name=?", "huravar").Updates(map[string]interface{}{"Phone": "14567890", "Salt": "09876"})

	// db.Unscoped().Delete(&model.BasicUserInformation{}, 2)
	db.Create(&model.BasicUserInformation{UserName: "huravar"})
	ans := model.BasicUserInformation{}
	db.Find(&ans)
	fmt.Println(reflect.TypeOf(ans.ID))

}

func DatabaseLoginString() string {
	return fmt.Sprintf("%v:%v@%v(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
		viper.Get("mysql.username"), viper.Get("mysql.password"), viper.Get("mysql.network_protocal"),
		viper.Get("mysql.server_address"), viper.Get("mysql.server_port"), viper.Get("mysql.database_name"),
		viper.Get("mysql.character"))
}

func TimePointer() *time.Time {
	ItemTime := time.Now()
	return &ItemTime
}
