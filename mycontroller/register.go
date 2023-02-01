package mycontroller

import (
	"net/http"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/service"
	"github.com/gin-gonic/gin"
)

// response to /douyin/user/register/ api
type UserRegisterResponse struct {
	common.Response
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

// user register controller
func UserRegister(c *gin.Context) {
	// 1. get the username and password
	username := c.Query("username")
	password := c.Query("password")

	// 2. call the service of user register
	tokenResponse, err := service.UserRegisterService(username, password)

	// 3. register error, return the error message
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{
				StatusCode: 1,
				// the value of StatusMsg depends on the type of return error
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	// 4. register successed, return the register response
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response: common.Response{StatusCode: 0},
		UserID:   tokenResponse.UserID,
		Token:    tokenResponse.Token,
	})
}
