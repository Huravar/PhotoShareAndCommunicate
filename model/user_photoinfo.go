package model

import (
	"fmt"
	"gorm.io/gorm"
	"photo_service/utils"
)

type UserPhotoInfo struct {
	gorm.Model
	UserID       uint   //DBName:user_id
	UserName     string //DBName:user_name
	OriginalName string //DBName:original_name
	StoragePath  string //DBName:storage_path
	FileSize     int64  //DBName:file_size
	FileType     string //DBName:file_type
	HitNum       uint   //DBName:hit_num
	Description  string //DBName:description
}

func (UserPhotoInfo) TableName() string {
	return "UserPhotoInfo"
}

func FindUserPhotoInfoByUserId(userID uint) ([]UserPhotoInfo, int64) {
	var IUserPhotoInfo []UserPhotoInfo
	IdbResult := utils.DB.Model(&UserPhotoInfo{}).Where("user_id = ?", userID).Find(&IUserPhotoInfo)
	return IUserPhotoInfo, IdbResult.RowsAffected
}

func CreatUserPhotoInfo(PUserPhotoInfo UserPhotoInfo) {
	utils.DB.Create(&PUserPhotoInfo)
}

func DeleteUserPhotoInfo(PPhotoName string) error {
	IdbResult := utils.DB.Model(&UserPhotoInfo{}).Where("original_name = ?", PPhotoName).Delete(&UserPhotoInfo{})
	if IdbResult.RowsAffected == 0 {
		return fmt.Errorf("删除%v失败！", PPhotoName)
	}
	return nil
}

func FindUserPhotoInfoByPhotoName(PUserId uint, PPhotoName string) (UserPhotoInfo, int64) {
	var IUserPhotoInfo UserPhotoInfo
	IdbResult := utils.DB.Model(&UserPhotoInfo{}).Where("original_name = ?", PPhotoName).
		Where("user_id=?", PUserId).Find(&IUserPhotoInfo)
	return IUserPhotoInfo, IdbResult.RowsAffected
}

func DeleteUserPhotoInfoByPhotoName(PUserId uint, POriginalName string) error {
	IdbResult := utils.DB.Model(&UserPhotoInfo{}).Where("original_name=?", POriginalName).
		Where("user_id=?", PUserId).Unscoped().Delete(&UserPhotoInfo{})
	if IdbResult.RowsAffected == 0 {
		return fmt.Errorf("%v删除失败！", POriginalName)
	}
	return nil
}
