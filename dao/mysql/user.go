package mysql

import (
	"errors"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/model"
	"gorm.io/gorm"
)

// create a user
func CreateAUser(user *model.User) (err error) {
	err = DB.Model(&model.User{}).Create(&user).Error
	if err != nil {
		return common.ErrorCreateUserFalse
	}
	return
}

func IsUserExist(username string) (flag bool, err error) {
	var userExist = &model.User{}
	err = DB.Model(&model.User{}).Where("name=?", username).First(&userExist).Error
	
	// can't find the user in the "user" tables means user doesn't exist
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// user doesn't exist
		return false, nil
	}
	// other unpredicted error
	if err != nil {
		return false, common.ErrorSQLFalse
	}
	// user exists
	return true, common.ErrorUserExist
}
