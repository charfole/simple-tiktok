package model

import "gorm.io/gorm"

type Comment struct { // 评论
	gorm.Model
	VideoID uint   `json:"video_id,omitempty"`
	UserID  uint   `json:"user_id,omitempty"`
	Content string `json:"content,omitempty"`
}
