package mysql

import (
	"errors"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/model"
	"gorm.io/gorm"
)

// Create a user
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

// Get user info by name
func GetAUserByName(username string, login *model.User) error {
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

// Get user info by ID
func GetAUserByID(userID uint, user *model.User) error {
	// DB.Where("id=?", userID).First(user)
	if err := DB.Model(&model.User{}).Where("id = ?", userID).Find(user).Error; err != nil {
		return err
	}
	return nil
}

func GetAUser(userID uint) (model.User, error) {
	var user model.User
	if err := DB.Model(&model.User{}).Where("id = ?", userID).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
