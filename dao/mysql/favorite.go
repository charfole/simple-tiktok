package mysql

import (
	"errors"

	"github.com/charfole/simple-tiktok/model"
	"gorm.io/gorm"
)

// Check the user likes the video or not
func IsFavorite(uid uint, vid uint) bool {
	var total int64
	if err := DB.Model(model.Favorite{}).
		Where("user_id = ? AND video_id = ? AND state = 1", uid, vid).Count(&total).
		Error; errors.Is(err, gorm.ErrRecordNotFound) { //没有该条记录
		return false
	}

	if total == 0 {
		return false
	}

	return true
}

func CreateAFavorite(favoriteAction *model.Favorite) (err error) {
	err = DB.Model(model.Favorite{}).Create(favoriteAction).Error
	return err
}

func UpdateFavoriteState(userID, videoID, state uint) (err error) {
	err = DB.Model(model.Favorite{}).Where("user_id = ? AND video_id = ?", userID, videoID).
		Update("state", state).Error
	return
}

func IsFavoriteRecordExist(userID, videoID uint, favoriteStruct *model.Favorite) (err error) {
	err = DB.Model(model.Favorite{}).Where("user_id = ? AND video_id = ?", userID, videoID).
		First(favoriteStruct).Error
	return
}

func GetFavoriteList(userID uint) (favoriteList []model.Favorite, err error) {
	err = DB.Model(model.Favorite{}).Where("user_id=? AND state=?", userID, 1).Find(&favoriteList).Error
	return favoriteList, err
}
