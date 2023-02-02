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

// func IsUserExist(username string) (flag bool, err error) {
// 	var userExist = &model.User{}
// 	err = DB.Model(&model.User{}).Where("name=?", username).First(&userExist).Error

// 	// user not found
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		// user not found
// 		return false, nil
// 	}
// 	// other unpredicted error
// 	if err != nil {
// 		return false, common.ErrorSQLFalse
// 	}
// 	// user found
// 	return true, common.ErrorUserExist
// }

func GetAUser(username string, login *model.User) error {
	err := DB.Where("name=?", username).First(login).Error
	// fmt.Printf("%+v", login)
	// user not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.ErrorUserNotFound
	}
	// other unpredicted error
	if err != nil {
		// fmt.Println("获取user错误")
		return common.ErrorSQLFalse
	}
	// user found
	return nil
}
