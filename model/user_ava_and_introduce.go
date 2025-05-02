package model

import (
	"fmt"
	"gorm.io/gorm"
	"photo_service/utils"
)

type UserHomePageInfo struct {
	gorm.Model
	UserID        uint   //DBName:user_id
	UserName      string //DBName:user_name
	AvatarPath    string //DBName:avatar_path
	SelfIntroduce string //DBName:self_introduce
}

func (UserHomePageInfo) TableName() string {
	return "user_home_page_info"
}

func CreateUserHomePageInfo(PUserHomePageInfo UserHomePageInfo) {
	utils.DB.Create(&PUserHomePageInfo)
}

func FindUserHomePageInfoByUserId(userId uint) (UserHomePageInfo, int64) {
	var PUserHomePageInfo UserHomePageInfo
	IdbResult := utils.DB.Where("user_id = ?", userId).Find(&PUserHomePageInfo)
	return PUserHomePageInfo, IdbResult.RowsAffected
}

func UpdateAvaPathById(PUserId uint, AvaPath string) error {
	IdbResult := utils.DB.Model(&UserHomePageInfo{}).Where("user_id=?", PUserId).Update("avatar_path", AvaPath)
	if IdbResult.RowsAffected == 0 {
		return fmt.Errorf("更新头像路径失败！")
	}
	return nil
}

func UpdateSlfIntroduceById(PUserId uint, PSlfIntroduce string) error {
	IdbResult := utils.DB.Model(&UserHomePageInfo{}).Where("user_id=?", PUserId).Update("self_introduce", PSlfIntroduce)
	if IdbResult.RowsAffected == 0 {
		return fmt.Errorf("更新头像路径失败！")
	}
	return nil
}
