package mycontroller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
	"github.com/charfole/simple-tiktok/service"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	common.Response
	VideoList []service.ReturnVideo `json:"video_list"`
}

func PublishList(c *gin.Context) { //获取列表的方法
	// 1. get the current user id by JWT
	rawHostID, _ := c.Get("user_id")
	hostID := rawHostID.(uint)

	// 2. get the guest id
	rawGuestID := c.Query("user_id")
	id, _ := strconv.Atoi(rawGuestID)
	guestID := uint(id)

	// if guestID equals to 0, it means the guest is current user
	if guestID == 0 {
		guestID = hostID
	}
	fmt.Println("hostID: ", hostID)
	fmt.Println("guestID: ", guestID)

	// 3. get the guest
	var user model.User
	err := mysql.GetAUserByID(guestID, &user)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorUserNotFound.Error(),
		})
		c.Abort()
		return
	}

	// 4. pack the guest
	returnAuthor := service.PackAuthor(user, hostID, guestID)

	// 5. get the video list of this guest(author), pack it and return
	videoList := mysql.GetVideoList(guestID)
	if len(videoList) > 0 {
		// videolist found, pack the videos and return
		var returnVideoList []service.ReturnVideo
		returnVideoList = service.PackVideo(videoList, returnAuthor, hostID)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: returnVideoList,
		})
	} else {
		// video list not found, return error
		c.JSON(http.StatusOK, VideoListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  common.ErrorVideoList.Error(),
			},
			VideoList: nil,
		})
	}
}
