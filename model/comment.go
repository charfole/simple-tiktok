package model

import "gorm.io/gorm"

type Comment struct { // 评论
	gorm.Model
	VideoID uint   `json:"video_id" gorm:"index"`
	UserID  uint   `json:"user_id" gorm:"index"`
	Content string `json:"content"`
}
