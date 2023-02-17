package mysql

import (
	"github.com/charfole/simple-tiktok/model"
	"gorm.io/gorm"
)

var followings = "followings"

func IsFollowing(HostID uint, GuestID uint) (err error) {
	var relationExist = &model.Following{}
	//判断关注是否存在
	err = DB.Model(&model.Following{}).
		Where("host_id=? AND guest_id=?", HostID, GuestID).
		First(&relationExist).Error
	return
}

// IncreaseFollowCount 增加HostID的关注数（Host_id 的 follow_count+1）
func IncreaseFollowCount(HostID uint) error {
	if err := DB.Model(&model.User{}).
		Where("id=?", HostID).
		Update("follow_count", gorm.Expr("follow_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// DecreaseFollowCount 减少HostID的关注数（Host_id 的 follow_count-1）
func DecreaseFollowCount(HostID uint) error {
	if err := DB.Model(&model.User{}).
		Where("id=?", HostID).
		Update("follow_count", gorm.Expr("follow_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// CreateFollowing 创建关注
func CreateFollowing(HostID uint, GuestID uint) error {

	//1.Following数据模型准备
	newFollowing := model.Following{
		HostID:  HostID,
		GuestID: GuestID,
	}

	//2.新建following
	if err := DB.Model(&model.Following{}).Create(&newFollowing).Error; err != nil {
		return err
	}
	return nil
}

// DeleteFollowing 删除关注
func DeleteFollowing(HostID uint, GuestID uint) error {
	//1.Following数据模型准备
	deleteFollowing := model.Following{
		HostID:  HostID,
		GuestID: GuestID,
	}

	//2.删除following
	if err := DB.Model(&model.Following{}).Where("host_id=? AND guest_id=?", HostID, GuestID).Delete(&deleteFollowing).Error; err != nil {
		return err
	}

	return nil
}

// FollowingList 获取关注表
func FollowingList(HostID uint) ([]model.User, error) {
	//1.userList数据模型准备
	var userList []model.User
	//2.查HostID的关注表
	err := DB.Model(&model.User{}).
		Joins("left join "+followings+" on "+users+".id = "+followings+".guest_id").
		Where(followings+".host_id=? AND "+followings+".deleted_at is null", HostID).
		Scan(&userList).Error
	return userList, err
}
