package mycontroller

import (
	"net/http"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/service"
	"github.com/gin-gonic/gin"
)

type UserLoginResponse struct {
	common.Response
	service.TokenResponse
}

// UserLogin controller
func UserLogin(c *gin.Context) {
	// get the username and password from app
	username := c.Query("username")
	password := c.Query("password")

	// call the service and get response
	userLoginResponse, err := service.UserLoginService(username, password)

	// if user not found, return error
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	// user found, login successfully, return the response
	c.JSON(http.StatusOK, UserLoginResponse{
		Response:      common.Response{StatusCode: 0},
		TokenResponse: userLoginResponse,
	})
}
