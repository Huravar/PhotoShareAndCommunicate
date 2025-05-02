package model

import (
	"gorm.io/gorm"
	"photo_service/utils"
	"time"
)

type UserNetwork struct {
	gorm.Model               //DBName: id,created_at,updated_at,deleted_at
	UserID        uint       //DBName:user_id
	UserName      string     //DBName:user_name
	ClientIp      string     //DBName:client_ip
	ClientPort    string     //DBName:clint_port
	UserTK        string     //DBName: user_tk  user token key
	LoginTime     *time.Time //DBName:login_time
	HeartbeatTime *time.Time //DBName:heartbeat_time
	LoginOutTime  *time.Time //DBName:login_out_time
	IsLogin       bool       //DBName:is_login
	DeviceInfo    string     //DBName:device_info
}

func TimePointer() *time.Time {
	ItemTime := time.Now()
	return &ItemTime
}

func (t UserNetwork) TableName() string {
	return "UserNetwork"
}

func CreateUserNetwork(PUserNetWork UserNetwork) {
	utils.DB.Model(&UserNetwork{}).Model(&UserNetwork{}).Create(&PUserNetWork)

}

func FindUserNetworkById(UserId uint) (UserNetwork, int64) {
	var UserNum int64
	var RUserNetwork UserNetwork
	utils.DB.Model(&UserNetwork{}).Model(&UserNetwork{}).Where("user_id = ?", UserId).Find(&RUserNetwork).Count(&UserNum)
	return RUserNetwork, UserNum
}

func UpdateNetworkForTk(PUserNetWork UserNetwork, Token string) int64 {
	EffRow := utils.DB.Model(&UserNetwork{}).Where("user_id = ?", PUserNetWork.UserID).Update("user_tk", Token)
	return EffRow.RowsAffected
}
func UpdateNetworkForClintNAdress(PUserNetWork UserNetwork, PClientIp, PClientPort string) int64 {
	EffRow := utils.DB.Model(&UserNetwork{}).Where("user_id = ?", PUserNetWork.UserID).Update("client_ip", PClientIp).
		Update("clint_port", PClientPort)
	return EffRow.RowsAffected
}

func UpdateNetworkForLoginTime(PUserNetWork UserNetwork) int64 {
	EffRow := utils.DB.Model(&UserNetwork{}).Where("user_id = ?", PUserNetWork.UserID).Update("login_time", TimePointer())
	return EffRow.RowsAffected
}

func UpdateNetworkForHeartbeatTime(PUserNetWork UserNetwork) int64 {
	EffRow := utils.DB.Model(&UserNetwork{}).Where("user_id = ?", PUserNetWork.UserID).Update("heartbeat_time", TimePointer())
	return EffRow.RowsAffected
}
func UpdateNetworkForLoginOutTime(PUserNetWork UserNetwork) int64 {
	EffRow := utils.DB.Model(&UserNetwork{}).Where("user_id = ?", PUserNetWork.UserID).Update("login_out_time", TimePointer())
	return EffRow.RowsAffected
}
func UpdateNetworkForIsLogout(PUserNetWork UserNetwork, PIsLogout bool) int64 {
	EffRow := utils.DB.Model(&UserNetwork{}).Where("user_id = ?", PUserNetWork.UserID).Update("is_logout", PIsLogout)
	return EffRow.RowsAffected
}
func UpdateNetworkForDeviceInfo(PUserNetWork UserNetwork, PDeviceInfo string) int64 {
	EffRow := utils.DB.Model(&UserNetwork{}).Where("user_id = ?", PUserNetWork.UserID).Update("device_info", PDeviceInfo)
	return EffRow.RowsAffected
}
