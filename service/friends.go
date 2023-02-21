package service

import (
	"fmt"

	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
)

func FriendListService(userId uint) ([]model.FriendUser, error) {
	var FriendUsers []model.FriendUser
	var err error
	fmt.Printf("FriendListService")
	FriendUsers, err = PackageFriendLists(userId)
	if err != nil {
		return nil, err
	}
	return FriendUsers, nil
}

func FriendList(userId uint) ([]model.User, error) {
	//2.查userId的关注表
	userList, err := mysql.FriendList(userId)
	if err != nil {
		return userList, err
	}
	return userList, nil
}

// PackageFriendLists 通过userId查询粉丝的数据，再包装加入消息
func PackageFriendLists(userId uint) ([]model.FriendUser, error) {
	var FriendLists []model.FriendUser
	var message string
	var msgType int8
	userInfos, err := FriendList(userId)
	if err != nil {
		return nil, err
	}
	fmt.Printf("获取了好友表,userID是", userId)
	fmt.Printf("获取了好友表,内容是是", userInfos)
	//for _, userInfo := range userInfos {
	//	FriendLists = append(FriendLists, model.FriendUser{
	//		User:    userInfo,
	//		Message: message,
	//		MsgType: msgType,
	//	})
	//}
	for _, userInfo := range userInfos {
		// 查询Message和MsgType
		message, msgType, err = mysql.QueryNewestMessageByUserIdAndToUserID(userId, userInfo.ID)
		if err != nil {
			return nil, err
		}
		fmt.Printf("ID1", userId)
		fmt.Printf("ID2", userInfo.ID)
		fmt.Printf("最新消息是", message)
		//userInfo.AvatarUrl = "https://xingqiu-tuchuang-1256524210.cos.ap-shanghai.myqcloud.com/12640/20230206133334.png"
		FriendLists = append(FriendLists, model.FriendUser{
			ID:              userInfo.Model.ID,
			Name:            userInfo.Name,
			FollowCount:     userInfo.FollowCount,
			FollowerCount:   userInfo.FollowerCount,
			TotalFavorited:  userInfo.TotalFavorited,
			FavoriteCount:   userInfo.FavoriteCount,
			Avatar:          userInfo.Avatar,
			BackgroundImage: userInfo.BackgroundImage,
			IsFollow:        true, //好友都是互相关注的
			Message:         message,
			MsgType:         msgType,
		})
	}

	return FriendLists, nil
}
