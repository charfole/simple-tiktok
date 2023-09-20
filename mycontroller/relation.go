package mycontroller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/middleware"
	"github.com/charfole/simple-tiktok/model"
	"github.com/charfole/simple-tiktok/service"
	"github.com/gin-gonic/gin"
)

// ReturnFollower 关注表与粉丝表共用的用户数据模型
type ReturnFollower struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	FollowCount     uint   `json:"follow_count"`
	FollowerCount   uint   `json:"follower_count"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	IsFollow        bool   `json:"is_follow"`
}

// FollowingListResponse 关注表相应结构体
type FollowingListResponse struct {
	common.Response
	UserList []ReturnFollower `json:"user_list"`
}

// FollowerListResponse 粉丝表相应结构体
type FollowerListResponse struct {
	common.Response
	UserList []ReturnFollower `json:"user_list"`
}

// RelationAction 关注/取消关注操作
func RelationAction(c *gin.Context) {
	//1.取数据
	//1.1 从token中获取用户id
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	hostID := tokenStruct.UserID
	//1.2 获取待关注的用户id
	getToUserID, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	guestID := uint(getToUserID)
	//1.3 获取关注操作（关注1，取消关注2）
	getActionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	actionType := uint(getActionType)

	//2.自己关注/取消关注自己不合法
	if hostID == guestID {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 405,
			StatusMsg:  "无法关注自己",
		})
		c.Abort()
		return
	}

	//3.service层进行关注/取消关注处理
	err := service.FollowAction(hostID, guestID, actionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  "关注/取消关注成功！",
		})
	}
}

// FollowList 获取用户关注列表
func FollowList(c *gin.Context) {

	//1.数据预处理
	//1.1获取用户本人id
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	hostID := tokenStruct.UserID
	//1.2获取其他用户id
	getGuestID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	guestID := uint(getGuestID)

	//2.判断查询类型，从数据库取用户列表
	var err error
	var userList []model.User
	if guestID == 0 {
		//若其他用户id为0，代表查本人的关注表
		userList, err = service.FollowingList(hostID)
	} else {
		//若其他用户id不为0，代表查对方的关注表
		userList, err = service.FollowingList(guestID)
	}

	//构造返回的数据
	var ReturnFollowerList = make([]ReturnFollower, len(userList))
	for i, m := range userList {
		ReturnFollowerList[i].ID = m.ID
		ReturnFollowerList[i].Name = m.Name
		ReturnFollowerList[i].FollowCount = m.FollowCount
		ReturnFollowerList[i].FollowerCount = m.FollowerCount
		ReturnFollowerList[i].Avatar = m.Avatar
		ReturnFollowerList[i].BackgroundImage = m.BackgroundImage
		ReturnFollowerList[i].IsFollow = service.IsFollowing(hostID, m.ID)
	}
	fmt.Printf("找到关注表", ReturnFollowerList)

	//3.响应返回
	if err != nil {
		c.JSON(http.StatusBadRequest, FollowingListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, FollowingListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			UserList: ReturnFollowerList,
		})
	}
}

// FollowerList 获取用户粉丝列表
func FollowerList(c *gin.Context) {

	//1.数据预处理
	//1.1获取用户本人id
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	hostID := tokenStruct.UserID
	//1.2获取其他用户id
	getGuestID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	guestID := uint(getGuestID)

	//2.判断查询类型
	var err error
	var userList []model.User
	if guestID == 0 {
		//查本人的粉丝表
		userList, err = service.FollowerList(hostID)
	} else {
		//查对方的粉丝表
		userList, err = service.FollowerList(guestID)
	}

	//3.判断查询类型，从数据库取用户列表
	var ReturnFollowerList = make([]ReturnFollower, len(userList))
	for i, m := range userList {
		ReturnFollowerList[i].ID = m.ID
		ReturnFollowerList[i].Name = m.Name
		ReturnFollowerList[i].FollowCount = m.FollowCount
		ReturnFollowerList[i].FollowerCount = m.FollowerCount
		ReturnFollowerList[i].Avatar = m.Avatar
		ReturnFollowerList[i].BackgroundImage = m.BackgroundImage
		ReturnFollowerList[i].IsFollow = service.IsFollowing(hostID, m.ID)
	}
	fmt.Printf("找到粉丝表", ReturnFollowerList)

	//3.处理
	if err != nil {
		c.JSON(http.StatusBadRequest, FollowerListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, FollowerListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			UserList: ReturnFollowerList,
		})
	}
}

type FriendListResponse struct {
	common.Response
	UserList []model.FriendUser `json:"user_list"` // 用户列表
}

// FriendList 好友列表
func FriendList(c *gin.Context) {
	//user_id := c.Query("user_id")
	token := c.Query("token")
	tokenStruct, check := middleware.CheckToken(token)
	// 4. if token of the login user not checked, return error
	if !check {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  common.ErrorTokenFalse.Error(),
			},
		})
		return
	}

	// 5. if checked, get the login userID
	//hostID := tokenStruct.UserID
	userID := tokenStruct.UserID
	//userID, err := strconv.ParseUint(user_id, 10, 64)
	fmt.Printf("点击了消息列表,user_id是", userID)

	friendUsers, err := service.FriendListService(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, FriendListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "获取好友成功",
		},
		UserList: friendUsers,
	})
	return
}
