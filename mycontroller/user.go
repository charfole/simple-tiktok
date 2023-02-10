package mycontroller

import (
	"net/http"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/middleware"
	"github.com/charfole/simple-tiktok/service"
	"github.com/gin-gonic/gin"
)

type UserInfoResponse struct {
	common.Response
	UserInfoList service.UserInfoQueryResponse `json:"user"`
}

// UserInfo controller
func UserInfo(c *gin.Context) {
	// 1. query the id of current guest user
	rawGuestID := c.Query("user_id")
	userInfoResponse, err := service.UserInfoService(rawGuestID)

	// 2. user not found, return error
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	// 3. query the id of the login user
	token := c.Query("token")
	tokenStruct, check := middleware.CheckToken(token)

	// 4. if token of the login user not checked, return error
	if !check {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	// 5. if checked, get the login userID
	hostID := tokenStruct.UserID

	// 6. check the login user follows the guest user or not, then update the "IsFollow" filed
	userInfoResponse.IsFollow = service.CheckIsFollow(hostID, rawGuestID)

	// 7. all done, return
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		UserInfoList: userInfoResponse,
	})

}
