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

// FavoriteList 获取列表方法
func FavoriteList(c *gin.Context) {
	//user_id获取
	getUserID, _ := c.Get("user_id")
	var userIDHost uint
	if v, ok := getUserID.(uint); ok {
		userIDHost = v
	}
	userIDStr := c.Query("user_id") //自己id或别人id
	userID, _ := strconv.ParseUint(userIDStr, 10, 10)
	userIDNew := uint(userID)
	if userIDNew == 0 {
		userIDNew = userIDHost
	}

	//函数调用及响应
	videoList, err := service.FavoriteList(userIDNew)
	videoListNew := make([]FavoriteVideo, 0)
	for _, m := range videoList {
		var author = FavoriteAuthor{}
		var getAuthor = model.User{}
		err := mysql.GetAUserByID(m.AuthorID, &getAuthor) //参数类型、错误处理
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "找不到作者！",
			})
			c.Abort()
			return
		}
		//isfollowing
		isfollowing := service.IsFollowing(userIDHost, m.AuthorID) //参数类型、错误处理
		//isfavorite
		isfavorite := mysql.IsFavorite(userIDHost, m.ID)
		//作者信息
		author.ID = getAuthor.ID
		author.Name = getAuthor.Name
		author.FollowCount = getAuthor.FollowCount
		author.FollowerCount = getAuthor.FollowerCount
		author.IsFollow = isfollowing
		//组装
		var video = FavoriteVideo{}
		video.ID = m.ID //类型转换
		video.Author = author
		video.PlayURL = m.PlayURL
		video.CoverURL = m.CoverURL
		video.FavoriteCount = m.FavoriteCount
		video.CommentCount = m.CommentCount
		video.IsFavorite = isfavorite
		video.Title = m.Title

		videoListNew = append(videoListNew, video)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, FavoriteListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			VideoList: nil,
		})
	} else {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			VideoList: videoListNew,
		})
	}
}
