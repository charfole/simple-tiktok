package mysql

import (
	"github.com/charfole/simple-tiktok/model"
	"gorm.io/gorm"
)

// 粉丝表
var followers = "followers"

// 用户表
var users = "users"

func IsFollower(HostID uint, GuestID uint) (err error) {
	var relationExist = &model.Followers{}
	//判断关注是否存在
	err = DB.Model(&model.Followers{}).
		Where("host_id=? AND guest_id=?", HostID, GuestID).
		First(&relationExist).Error
	return
}

// FollowerList  获取粉丝表
func FollowerList(HostID uint) ([]model.User, error) {
	//1.userList数据模型准备
	var userList []model.User
	//2.查HostID的关注表
	err := DB.Model(&model.User{}).
		Joins("left join "+followers+" on "+users+".id = "+followers+".guest_id").
		Where(followers+".host_id=? AND "+followers+".deleted_at is null", HostID).
		Scan(&userList).Error
	return userList, err
}

// IncreaseFollowerCount 增加HostID的粉丝数（Host_id 的 follow_count+1）
func IncreaseFollowerCount(HostID uint) error {
	if err := DB.Model(&model.User{}).
		Where("id=?", HostID).
		Update("follower_count", gorm.Expr("follower_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// DecreaseFollowerCount 减少HostID的粉丝数（Host_id 的 follow_count-1）
func DecreaseFollowerCount(HostID uint) error {
	if err := DB.Model(&model.User{}).
		Where("id=?", HostID).
		Update("follower_count", gorm.Expr("follower_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// CreateFollower 创建粉丝
func CreateFollower(HostID uint, GuestID uint) error {

	//1.Following数据模型准备
	newFollower := model.Followers{
		HostID:  HostID,
		GuestID: GuestID,
	}

	//2.新建following
	if err := DB.Model(&model.Followers{}).
		Create(&newFollower).Error; err != nil {
		return err
	}
	return nil
}

// DeleteFollower 删除粉丝
func DeleteFollower(HostID uint, GuestID uint) error {
	//1.Following数据模型准备
	newFollower := model.Followers{
		HostID:  HostID,
		GuestID: GuestID,
	}

	//2.删除following
	if err := DB.Model(&model.Followers{}).
		Where("host_id=? AND guest_id=?", HostID, GuestID).
		Delete(&newFollower).Error; err != nil {
		return err
	}

	return nil
}

// 获取好友表，即互相都是对方粉丝
func FriendList(userID uint) ([]model.User, error) {
	var userList []model.User
	//查粉丝表
	err := DB.Model(&model.User{}).
		Where("users.ID IN (SELECT a.host_id FROM followers a JOIN followers b ON a.host_id  = b.guest_id AND a.guest_id = b.host_id  AND a.guest_id = ? AND a.deleted_at is null AND b.deleted_at is null)", userID).
		Scan(&userList).Error
	return userList, err
}
