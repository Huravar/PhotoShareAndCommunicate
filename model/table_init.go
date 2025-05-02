package model

import (
	"fmt"
	"photo_service/utils"
)

func TableInit() error {

	err := utils.DB.AutoMigrate(&BasicUserInformation{})
	if err != nil {
		return fmt.Errorf("create BasicUserInformation table failed: %v", err)
	}
	err = utils.DB.AutoMigrate(&UserNetwork{})
	if err != nil {
		return fmt.Errorf("create UserNetwork table failed: %v", err)
	}
	err = utils.DB.AutoMigrate(&UserHomePageInfo{})
	if err != nil {
		return fmt.Errorf("create UserHomePageInfo table failed: %v", err)
	}
	return nil
}
