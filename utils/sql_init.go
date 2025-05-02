package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func DatabaseLoginString() string {
	return fmt.Sprintf("%v:%v@%v(%v:%v)/%v?charset=%v&parseTime=True&loc=Local",
		viper.Get("mysql.username"), viper.Get("mysql.password"), viper.Get("mysql.network_protocal"),
		viper.Get("mysql.server_address"), viper.Get("mysql.server_port"), viper.Get("mysql.database_name"),
		viper.Get("mysql.character"))
}

func OpenDatabase() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		},
	)
	itemDB, err := gorm.Open(mysql.Open(DatabaseLoginString()), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	DB = itemDB
}
