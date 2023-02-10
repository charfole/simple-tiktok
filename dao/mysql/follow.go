package mysql

import (
	"github.com/charfole/simple-tiktok/model"
)

func IsFollowing(HostID uint, GuestID uint) (err error) {
	var relationExist = &model.Following{}
	//判断关注是否存在
	err = DB.Model(&model.Following{}).
		Where("host_id=? AND guest_id=?", HostID, GuestID).
		First(&relationExist).Error
	return
}
