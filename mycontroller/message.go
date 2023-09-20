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

type MessageListResponse struct {
	common.Response
	MessageList []model.Message `json:"message_list"`
}

// MessageChat 聊天记录 点进去才能看到
func MessageChat(c *gin.Context) {
	pre_msg_time := c.Query("pre_msg_time")
	to_user_id := c.Query("to_user_id")
	toUserId, _ := strconv.ParseUint(to_user_id, 10, 64)
	fmt.Printf("收到pre_msg time", pre_msg_time)
	// 提取用户Id
	//userid, exists := c.Get("Id")
	//if !exists {
	//	c.JSON(http.StatusNotFound, model.BaseResponse{
	//		StatusCode: -1,
	//		StatusMsg:  config.TokenIsNotExist,
	//	})
	//	return
	//}
	// 3. query the id of the login user
	PreMsgTime, _ := strconv.ParseInt(pre_msg_time, 10, 64)
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
	userId := tokenStruct.UserID
	fmt.Printf("")
	//userId := int64(userid.(float64))

	messages, err := service.MessageChatService(userId, uint(toUserId), PreMsgTime)
	//这里是为了解决左右显示的问题

	var MyMessageList = make([]model.Message, len(messages))
	for i, m := range messages {
		if m.UserId == userId { //如果消息的发送者是登录的本人，则0->1
			MyMessageList[i].UserId = 0
			MyMessageList[i].ToUserId = 1
		}
		if m.ToUserId == userId { //如果消息的接收者是登录的本人，则1->0
			MyMessageList[i].UserId = 1
			MyMessageList[i].ToUserId = 0
		}

		MyMessageList[i].Content = m.Content
		MyMessageList[i].CreateTime = m.CreateTime
		// MyMessageList[i].IsWithdraw = m.IsWithdraw
		MyMessageList[i].ID = m.ID
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: -1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, MessageListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "Success",
		},
		MessageList: MyMessageList,
	})
	return
}

// MessageAction 发送消息
func MessageAction(c *gin.Context) {
	// 提取用户Id
	//userid, exists := c.Get("Id")
	//fmt.Printf("", exists)
	//if !exists {
	//	c.JSON(http.StatusNotFound, common.BaseResponseInstance.FailMsg(config.TokenIsNotExist))
	//	return
	//}
	//userId := int64(userid.(float64))
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
	userId := tokenStruct.UserID

	// 获取toUser
	to_user_id := c.Query("to_user_id")
	toUserId, err := strconv.ParseUint(to_user_id, 10, 64)
	content := c.Query("content")
	actionType := c.Query("action_type")
	// 参数错误
	if actionType != "1" {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  common.ErrorTokenFalse.Error(),
			},
		})
		return
	}
	pass, err := service.MessageActionService(userId, uint(toUserId), content)
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: -1,
				StatusMsg:  "failedmsg",
			},
		})
		return
	}
	if !pass {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: -1,
				StatusMsg:  "fail",
			},
		})
		return
	} else {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
		})
		return
	}

}
