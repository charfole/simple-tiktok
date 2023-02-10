package model

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID  uint `json:"user_id"`
	VideoID uint `json:"video_id"`
	// state equals 1 means favorite, 0 means soft delete
	State uint
}
