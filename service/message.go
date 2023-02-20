package service

import (
	"errors"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
)

// 定义键值对维护消息记录 用户的Id->用户目前的消息记录索引
var userCommentIndex = make(map[uint]int64)

// 定义键值对维护消息记录 用户的Id->用户目前的消息记录最大值
var userMessageMaxIndex = make(map[uint]int64)

func MessageChatService(userId uint, toUserId uint, pre_msg_time int64) ([]model.Message, error) {
	// 查询userid和toUserId的表
	messages, err := mysql.QueryMessageByUserIdAndToUserId(userId, toUserId, pre_msg_time)
	if err != nil {
		return nil, err
	}
	return messages, nil

}

func MessageActionService(userId uint, toUserId uint, content string) (bool, error) {
	pass, err := mysql.InsertMessage(userId, toUserId, content)
	if err != nil {
		return false, err
	}
	if !pass {
		return false, errors.New("发送消息错误")
	}
	return true, nil
}
