package mycontroller

import (
	"net/http"
	"strconv"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
	"github.com/charfole/simple-tiktok/service"

	"github.com/gin-gonic/gin"
)

type FavoriteAuthor struct { //从user中获取,getUser函数
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"` //从following或follower中获取
}

type FavoriteVideo struct { //从video中获取
	ID            uint           `json:"id,omitempty"`
	Author        FavoriteAuthor `json:"author,omitempty"`
	PlayURL       string         `json:"play_url" json:"play_url,omitempty"`
	CoverURL      string         `json:"cover_url,omitempty"`
	FavoriteCount uint           `json:"favorite_count,omitempty"`
	CommentCount  uint           `json:"comment_count,omitempty"`
	IsFavorite    bool           `json:"is_favorite,omitempty"` //true
	Title         string         `json:"title,omitempty"`
}

type FavoriteListResponse struct {
	common.Response
	VideoList []FavoriteVideo `json:"video_list,omitempty"`
}

// Favorite 点赞视频方法
func Favorite(c *gin.Context) {
	// 1. get the userID through middleware
	getUserID, _ := c.Get("user_id")

	// 2. convert the data type from any to uint
	var userID uint
	if v, ok := getUserID.(uint); ok {
		userID = v
	}

	// 3. get the action_type and video_id from app
	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.ParseUint(actionTypeStr, 10, 64)
	videoIDStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 64)

	// 4. call the FavoriteAction function and return the message
	err := service.FavoriteAction(userID, uint(videoID), uint(actionType))
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  "操作成功！",
		})
	}
}

// FavoriteList the controller to get the favorite list
func FavoriteList(c *gin.Context) {
	// 1. get the id of login user
	getHostID, _ := c.Get("user_id")
	var hostID uint
	if v, ok := getHostID.(uint); ok {
		hostID = v
	}

	// 2. get the id of query user
	userIDStr := c.Query("user_id")
	userIDRaw, _ := strconv.ParseUint(userIDStr, 10, 10)
	userID := uint(userIDRaw)
	if userID == 0 {
		userID = hostID
	}

	// 3. get the raw favorite video list and pack the response
	rawVideoList, err := service.FavoriteList(userID)
	favoriteVideoList := make([]FavoriteVideo, 0)
	for _, v := range rawVideoList {
		var getAuthor model.User
		var author FavoriteAuthor
		var video FavoriteVideo

		// 4. get the video author by id
		err := mysql.GetAUserByID(v.AuthorID, &getAuthor)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "找不到作者！",
			})
			c.Abort()
			return
		}

		// 5. check the host follows the video author or not
		isfollowing := service.IsFollowing(hostID, v.AuthorID)

		// 6. pack the author struct in response
		author.ID = getAuthor.ID
		author.Name = getAuthor.Name
		author.FollowCount = getAuthor.FollowCount
		author.FollowerCount = getAuthor.FollowerCount
		author.IsFollow = isfollowing

		// 7. check the host favors the video or not
		isfavorite := mysql.IsFavorite(hostID, v.ID)

		// 8. pack the video struct in response
		video.ID = v.ID //类型转换
		video.Author = author
		video.PlayURL = v.PlayURL
		video.CoverURL = v.CoverURL
		video.FavoriteCount = v.FavoriteCount
		video.CommentCount = v.CommentCount
		video.IsFavorite = isfavorite
		video.Title = v.Title

		favoriteVideoList = append(favoriteVideoList, video)
	}

	// 9. error found, return error
	if err != nil {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "获取喜欢列表失败！",
			},
			VideoList: nil,
		})
	} else {
		// 10. return the favorite list
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "已找到喜欢列表！",
			},
			VideoList: favoriteVideoList,
		})
	}
}
