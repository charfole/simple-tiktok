package service

import (
	"errors"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
	"gorm.io/gorm"
)

// 粉丝表
var followers = "followers"

// 用户表
var users = "users"

// IsFollower 判断HostID是否有GuestID这个粉丝
func IsFollower(HostID uint, GuestID uint) bool {
	//2.查询粉丝表中粉丝是否存在
	err := mysql.IsFollower(HostID, GuestID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// follower not found
		return false
	}
	// follower found
	return true
}

// FollowerList  获取粉丝表
func FollowerList(HostID uint) ([]model.User, error) {
	//2.查HostID的关注表
	userList, err := mysql.FollowerList(HostID)
	if err != nil {
		return userList, err
	}
	return userList, nil
}
