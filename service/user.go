package service

import (
	"strconv"

	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
)

type UserInfoQueryResponse struct {
	UserID          uint   `json:"user_id"`
	UserName        string `json:"name"`
	FollowCount     uint   `json:"follow_count"`
	FollowerCount   uint   `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	TotalFavorited  uint   `json:"total_favorited"`
	FavoriteCount   uint   `json:"favorite_count"`
}

// UserInfoService
func UserInfoService(rawID string) (UserInfoQueryResponse, error) {
	// 1. prepare for the response
	var userInfoQueryResponse = UserInfoQueryResponse{}

	// 2. convert the string to int
	userID, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		return userInfoQueryResponse, err
	}

	// 3. query the user info
	var user model.User
	err = mysql.GetAUserByID(uint(userID), &user)
	// error found and error
	if err != nil {
		return userInfoQueryResponse, err
	}

	// 4. return the user info
	userInfoQueryResponse = UserInfoQueryResponse{
		UserID:          user.Model.ID,
		UserName:        user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		TotalFavorited:  user.TotalFavorited,
		FavoriteCount:   user.FavoriteCount,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		IsFollow:        false,
	}
	return userInfoQueryResponse, nil
}

// CheckIsFollow 检验已登录用户是否关注目标用户
func CheckIsFollow(hostID uint, rawGuestID string) bool {
	// 1. convert the string to int
	guestID, err := strconv.ParseUint(rawGuestID, 10, 64)
	if err != nil {
		return false
	}

	// 2. if the guestID equals to the hostID, means not follow
	if uint(guestID) == hostID {
		return false
	}

	// 3. check the host follows the guest or not
	return IsFollowing(hostID, uint(guestID))
}
