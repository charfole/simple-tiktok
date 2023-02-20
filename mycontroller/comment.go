package mycontroller

import (
	"net/http"
	"strconv"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
	"github.com/charfole/simple-tiktok/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// CommentListResponse 评论表的响应结构体
type CommentListResponse struct {
	common.Response
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

// CommentActionResponse 评论操作的响应结构体
type CommentActionResponse struct {
	common.Response
	Comment CommentResponse `json:"comment,omitempty"`
}

// UserResponse 用户信息的响应结构体
type UserResponse struct {
	ID            uint   `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint   `json:"follow_count,omitempty"`
	FollowerCount uint   `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// CommentResponse 评论信息的响应结构体
type CommentResponse struct {
	ID         uint         `json:"id,omitempty"`
	Content    string       `json:"content,omitempty"`
	CreateDate string       `json:"create_date,omitempty"`
	User       UserResponse `json:"user,omitempty"`
}

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	//1 数据处理
	getUserID, _ := c.Get("user_id")
	var userID uint
	if v, ok := getUserID.(uint); ok {
		userID = v
	}
	actionType := c.Query("action_type")
	videoIDStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 10)

	// 2 判断评论操作类型：1代表发布评论，2代表删除评论
	//2.1 非合法操作类型
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 405,
			StatusMsg:  "Unsupported actionType",
		})
		c.Abort()
		return
	}
	//2.2 合法操作类型
	if actionType == "1" { // 发布评论
		text := c.Query("comment_text")
		PostComment(c, userID, text, uint(videoID))
	} else if actionType == "2" { //删除评论
		commentIDStr := c.Query("comment_id")
		commentID, _ := strconv.ParseInt(commentIDStr, 10, 10)
		DeleteComment(c, uint(videoID), uint(commentID))
	}

}

// PostComment 发布评论
func PostComment(c *gin.Context, userID uint, text string, videoID uint) {
	//1 准备数据模型
	newComment := model.Comment{
		VideoID: videoID,
		UserID:  userID,
		Content: text,
	}

	//2 调用service层发布评论并改变评论数量，获取video作者信息
	err1 := mysql.DB.Transaction(func(db *gorm.DB) error {
		if err := mysql.PostComment(newComment); err != nil {
			return err
		}
		if err := mysql.AddCommentCount(videoID); err != nil {
			return err
		}
		return nil
	})
	// getUser, err2 := service.GetUser(userID)
	var getUser model.User
	err2 := mysql.GetAUserByID(userID, &getUser)
	videoAuthor, err3 := mysql.GetVideoAuthorID(videoID)

	//3 响应处理
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to post comment",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, CommentActionResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "post the comment successfully",
		},
		Comment: CommentResponse{
			ID:         newComment.ID,
			Content:    newComment.Content,
			CreateDate: newComment.CreatedAt.Format("01-02"),
			User: UserResponse{
				ID:            getUser.ID,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      service.IsFollowing(userID, videoAuthor),
			},
		},
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context, videoID uint, commentID uint) {

	//1 调用service层删除评论并改变评论数量，获取video作者信息
	err := mysql.DB.Transaction(func(db *gorm.DB) error {
		if err := mysql.DeleteComment(commentID); err != nil {
			return err
		}
		if err := mysql.ReduceCommentCount(videoID); err != nil {
			return err
		}
		return nil
	})
	//2 响应处理
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to delete comment",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  "delete the comment successfully",
	})
}

// CommentList 获取评论表
func CommentList(c *gin.Context) {
	//1 数据处理
	getUserID, _ := c.Get("user_id")
	var userID uint
	if v, ok := getUserID.(uint); ok {
		userID = v
	}
	videoIDStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 10)

	//2.调用service层获取指定videoid的评论表
	commentList, err := mysql.GetCommentList(uint(videoID))

	//2.1 评论表不存在
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 403,
			StatusMsg:  "Failed to get commentList",
		})
		c.Abort()
		return
	}

	//2.2 评论表存在
	var responseCommentList []CommentResponse
	for i := 0; i < len(commentList); i++ {
		// getUser, err1 := service.GetUser(commentList[i].UserID)
		var getUser model.User
		err1 := mysql.GetAUserByID(commentList[i].UserID, &getUser)

		if err1 != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 403,
				StatusMsg:  "Failed to get commentList.",
			})
			c.Abort()
			return
		}
		responseComment := CommentResponse{
			ID:         commentList[i].ID,
			Content:    commentList[i].Content,
			CreateDate: commentList[i].CreatedAt.Format("01-02"), // mm-dd
			User: UserResponse{
				ID:            getUser.ID,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      service.IsFollowing(userID, commentList[i].ID),
			},
		}
		responseCommentList = append(responseCommentList, responseComment)

	}

	//响应返回
	c.JSON(http.StatusOK, CommentListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "Successfully obtained the comment list.",
		},
		CommentList: responseCommentList,
	})

}
