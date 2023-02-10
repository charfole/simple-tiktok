package mycontroller

import (
	"net/http"
	"strconv"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/service"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	common.Response
	VideoList []service.FeedVideo `json:"video_list,omitempty"`
	NextTime  uint                `json:"next_time,omitempty"`
}

func Feed(c *gin.Context) {
	// 1. get the token of current user
	strToken := c.Query("token")

	// 2. get the latest_time
	var strLatestTime = c.Query("latest_time")
	// convert to integer
	latestTime, err := strconv.ParseInt(strLatestTime, 10, 32)
	if err != nil {
		latestTime = 0
	}

	// 3. get the feed video before latestTime
	videoList, _ := service.FeedGet(latestTime)

	// 4. pack the response by token and videoList
	feedVideoList, newTime := service.PackFeedResponse(strToken, videoList)

	// 5. get new feed, return the feed video
	if len(feedVideoList) > 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  common.Response{StatusCode: 0},
			VideoList: feedVideoList,
			NextTime:  uint(newTime),
		})
	} else {
		// 6. no new feed, reset the latest time to get the latest video again
		c.JSON(http.StatusOK, FeedResponse{
			Response:  common.Response{StatusCode: 0},
			VideoList: nil,
			NextTime:  0, // when NextTime equals to 0 equals, the latest time will be reset
		})
	}
}
