package model

import "gorm.io/gorm"

type Following struct {
	gorm.Model
	HostID  uint
	GuestID uint
}

type Followers struct {
	gorm.Model
	HostID  uint
	GuestID uint `gorm:"index"`
}
