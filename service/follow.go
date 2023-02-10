package service

import (
	"errors"

	"github.com/charfole/simple-tiktok/dao/mysql"

	"gorm.io/gorm"
)

// IsFollowing check if HostID follows GuestID
func IsFollowing(HostID uint, GuestID uint) bool {
	// err := mysql.DB.Model(&model.Following{}).Where("host_id=? AND guest_id=?", HostID, GuestID).
	// 	First(&relationExist).Error

	// couldn't follow myself
	if HostID == GuestID {
		return false
	}
	err := mysql.IsFollowing(HostID, GuestID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// follow not found
		return false
	}
	// follow found
	return true
}
