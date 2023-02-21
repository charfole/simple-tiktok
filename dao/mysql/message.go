package mysql

import (
	"time"

	"github.com/charfole/simple-tiktok/model"
)

// InsertMessage 插入数据
func InsertMessage(userId uint, toUserId uint, content string) (bool, error) {
	messageInfo := model.Message{
		ToUserId: toUserId,
		UserId:   userId,
		Content:  content,
		// CreateTime: time.Now().UnixNano() / 1e6,
		CreateTime: time.Now().Unix(),
	}
	// INSERT INTO `messages` (`user_id`,`to_user_id`,`content`,`is_withdraw`,`createTime`) VALUES (5,1,'111',0,'2023-02-08 19:21:15.017')
	result := DB.Debug().Create(&messageInfo)
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// QueryMessageByUserId 根据用户Id查询聊天记录
func QueryMessageByUserId(userId uint) ([]model.Message, error) {
	var messages []model.Message
	//SELECT * FROM `messages` WHERE user_id = 1 AND is_withdraw = 0 LIMIT 10
	result := DB.Where("user_id = ?", userId).Limit(5).Find(&messages) //limit2条消息
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

// QueryNewestMessageByUserId 通过用户Id查询最新的聊天记录 0-接受 1-发送 有点问题
func QueryNewestMessageByUserId(userId uint) (string, int8, error) {
	message := model.Message{}
	// SELECT * FROM `messages` WHERE (user_id = 1 Or to_user_id = 1) AND is_withdraw = 0 ORDER BY createTime desc LIMIT 1
	result := DB.Debug().Where("(user_id = ? Or to_user_id = ?) AND is_withdraw = ?", userId, userId, 0).Order("createTime desc").Limit(1).Find(&message)
	if result.Error != nil {
		return "", -1, result.Error
	}
	if userId == message.UserId {
		return message.Content, 1, nil
	} else {
		return message.Content, 0, nil
	}
}

// QueryNewestMessageByUserIdAndToUserID 通过两者的用户Id查询最新最新的两者之间的聊天记录 0-接受 1-发送 有点问题
func QueryNewestMessageByUserIdAndToUserID(userId uint, toUserId uint) (string, int8, error) {
	message := model.Message{}
	// SELECT `content`,`createTime`,`user_id`,`to_user_id` FROM `messages` WHERE (user_id = 2 AND to_user_id = 7 AND is_withdraw = 0) OR (user_id = 7 AND to_user_id = 2 AND is_withdraw = 0) ORDER BY createTime desc LIMIT 1
	result := DB.Debug().Where("user_id = ? AND to_user_id = ?", userId, toUserId).Or("user_id = ? AND to_user_id = ?", toUserId, userId).Order("createTime desc").Limit(1).Find(&message)
	if result.Error != nil {
		return "", -1, result.Error
	}
	if userId == message.UserId {
		return message.Content, 1, nil
	} else {
		return message.Content, 0, nil
	}
}

// QueryMessageByUserIdAndToUserId 查询两者的全部聊天记录
func QueryMessageByUserIdAndToUserId(userId uint, toUserId uint, pre_msg_time int64) ([]model.Message, error) {
	var messages []model.Message
	// SELECT * FROM `messages` WHERE (user_id = 1 AND to_user_id = 2 AND is_withdraw = 0) OR (user_id = 2 AND to_user_id = 1 AND is_withdraw = 0) ORDER BY createTime desc
	result := DB.Debug().Where("user_id = ? AND to_user_id = ? AND createTime > ? ", userId, toUserId, pre_msg_time).Or("user_id = ? AND to_user_id = ? AND createTime > ? ", toUserId, userId, pre_msg_time).Find(&messages)
	if result.Error != nil {
		return messages, result.Error
	}
	//for _, message := range messages {
	//	message.CreateTime = message.CreateTime.Unix()
	//}
	return messages, nil
}

// QueryMessageMaxCount 查询消息记录的最大值
func QueryMessageMaxCount(userId uint, toUserId uint) (int64, error) {
	var count int64
	// SELECT count(*) FROM `messages` WHERE (user_id = 1 AND to_user_id = 2 ) OR (user_id = 2 AND to_user_id = 1 )
	result := DB.Model(&model.Message{}).Where("user_id = ? AND to_user_id = ? ", userId, toUserId).Or("user_id = ? AND to_user_id = ? ", toUserId, userId).Count(&count)
	if result.Error != nil {
		return -1, result.Error
	}
	return count, nil
}
