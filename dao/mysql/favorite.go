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
