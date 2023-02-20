package mysql

import (
	"errors"
	"fmt"

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

// AddTotalFavorited 增加total_favorited
func AddTotalFavorited(HostID uint) error {
	fmt.Printf("AddTotalFavorited: 增加%d的获赞数目\n", HostID)
	if err := DB.Model(&model.User{}).
		Where("id=?", HostID).
		Update("total_favorited", gorm.Expr("total_favorited+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// ReduceTotalFavorited 减少total_favorited
func ReduceTotalFavorited(HostID uint) error {
	fmt.Printf("ReduceTotalFavorited: 减少%d的获赞数目\n", HostID)
	if err := DB.Model(&model.User{}).
		Where("id=?", HostID).
		Update("total_favorited", gorm.Expr("total_favorited-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// AddFavoriteCount 增加favorite_count
func AddFavoriteCount(HostID uint) error {
	fmt.Printf("AddFavoriteCount: 增加%d的喜欢数目\n", HostID)
	if err := DB.Model(&model.User{}).
		Where("id=?", HostID).
		Update("favorite_count", gorm.Expr("favorite_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// ReduceFavoriteCount 减少favorite_count
func ReduceFavoriteCount(HostID uint) error {
	fmt.Printf("ReduceFavoriteCount: 减少%d的喜欢数目\n", HostID)
	if err := DB.Model(&model.User{}).
		Where("id=?", HostID).
		Update("favorite_count", gorm.Expr("favorite_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}
