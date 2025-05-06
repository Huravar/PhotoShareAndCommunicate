package model

import (
	"log"
	"photo_service/utils"

	"gorm.io/gorm"
)

type BasicUserInformation struct {
	gorm.Model        //DBName: id,created_at,updated_at,deleted_at
	UserName   string //DBName: user_name
	PassWord   string //DBName: pass_word
	Phone      string //DBName: phone
	Email      string //DBName: email
	Identity   uint   //DBName: identity 0：普通用户，1：管理员，2：网站拥有者

}

func (table *BasicUserInformation) TableName() string {
	return "BasicUserInformation"
}

func GetUserList(User *BasicUserInformation) ([]*BasicUserInformation, int64) {
	var nums int64
	utils.DB.Model(&BasicUserInformation{}).Count(&nums)
	ans := make([]*BasicUserInformation, nums)
	utils.DB.Model(&BasicUserInformation{}).Find(&ans)
	return ans, nums
}

func CreatUserBasicInfo(user *BasicUserInformation) {
	utils.DB.Model(&BasicUserInformation{}).Create(user)
}

func FindUserByPhone(phone string) (BasicUserInformation, int64) {
	ans := BasicUserInformation{}
	result := utils.DB.Model(&BasicUserInformation{}).Where("phone = ?", phone).Find(&ans)
	return ans, result.RowsAffected
}

func FindUserById(id uint) (BasicUserInformation, int64) {
	ans := BasicUserInformation{}
	result := utils.DB.Model(&BasicUserInformation{}).Where("id = ?", id).Find(&ans)
	return ans, result.RowsAffected
}
func FindUserByName(name string) (BasicUserInformation, int64) {
	ans := BasicUserInformation{}
	var UserNum int64
	utils.DB.Model(&BasicUserInformation{}).Where("user_name = ?", name).Find(&ans).Count(&UserNum)
	return ans, UserNum
}

func AddUserRecord(User BasicUserInformation) {
	utils.DB.Model(&BasicUserInformation{}).Create(&User)
}

func UpdatePhoneById(id uint, phone string) error {
	IdbResult := utils.DB.Model(&BasicUserInformation{}).Where("id = ?", id).Update("phone", phone)
	if IdbResult.RowsAffected == 0 {
		log.Println("更新用户电话号失败")
		return IdbResult.Error
	}
	return nil
}

func UpdateEmailById(id uint, email string) error {
	IdbResult := utils.DB.Model(&BasicUserInformation{}).Where("id = ?", id).Update("email", email)
	if IdbResult.RowsAffected == 0 {
		log.Println("更新邮箱号号失败")
		return IdbResult.Error
	}
	return nil
}
