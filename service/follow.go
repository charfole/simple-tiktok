package service

import (
	"errors"
	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
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

// FollowingList 获取关注表
func FollowingList(HostID uint) ([]model.User, error) {
	//2.查HostID的关注表
	userList, err := mysql.FollowingList(HostID)
	if err != nil {
		return userList, err
	}
	return userList, nil
}

// FollowAction 关注操作
func FollowAction(HostID uint, GuestID uint, actionType uint) error {
	//创建关注操作
	if actionType == 1 {
		//判断关注是否存在
		if IsFollowing(HostID, GuestID) {
			//关注存在
			return common.ErrorRelationExit
		} else {
			//关注不存在,创建关注(启用事务Transaction)
			err1 := mysql.DB.Transaction(func(db *gorm.DB) error {
				err := mysql.CreateFollowing(HostID, GuestID)
				if err != nil {
					return err
				}
				err = mysql.CreateFollower(GuestID, HostID)
				if err != nil {
					return err
				}
				//增加host_id的关注数
				err = mysql.IncreaseFollowCount(HostID)
				if err != nil {
					return err
				}
				//增加guest_id的粉丝数
				err = mysql.IncreaseFollowerCount(GuestID)
				if err != nil {
					return err
				}
				return nil
			})
			if err1 != nil {
				return err1
			}
		}
	}
	if actionType == 2 {
		//判断关注是否存在
		if IsFollowing(HostID, GuestID) {
			//关注存在,删除关注(启用事务Transaction)
			if err1 := mysql.DB.Transaction(func(db *gorm.DB) error {
				err := mysql.DeleteFollowing(HostID, GuestID)
				if err != nil {
					return err
				}
				err = mysql.DeleteFollower(GuestID, HostID)
				if err != nil {
					return err
				}
				//减少host_id的关注数
				err = mysql.DecreaseFollowCount(HostID)
				if err != nil {
					return err
				}
				//减少guest_id的粉丝数
				err = mysql.DecreaseFollowerCount(GuestID)
				if err != nil {
					return err
				}
				return nil
			}); err1 != nil {
				return err1
			}

		} else {
			//关注不存在
			return common.ErrorRelationNull
		}
	}
	return nil
}
