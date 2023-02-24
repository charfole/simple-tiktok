package mycontroller

import (
	"fmt"
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
	CommentList []CommentResponse `json:"comment_list"`
}

// CommentActionResponse 评论操作的响应结构体
type CommentActionResponse struct {
	common.Response
	Comment CommentResponse `json:"comment"`
}

// UserResponse 用户信息的响应结构体
type UserResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	FollowCount     uint   `json:"follow_count"`
	FollowerCount   uint   `json:"follower_count"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	IsFollow        bool   `json:"is_follow"`
}

// CommentResponse 评论信息的响应结构体
type CommentResponse struct {
	ID         uint         `json:"id"`
	Content    string       `json:"content"`
	CreateDate string       `json:"create_date"`
	User       UserResponse `json:"user"`
}

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	// 1. get the login user id
	getUserID, _ := c.Get("user_id")
	var userID uint
	if v, ok := getUserID.(uint); ok {
		userID = v
	}

	// 2. get the action_type and video_id from app
	actionType := c.Query("action_type")
	videoIDStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 10)

	// 3. check the actionType is valid or not
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  "非法操作",
		})
		c.Abort()
		return
	}

	// 4. actionType == 1 means post comment, actionType == 2 means delete comment
	if actionType == "1" { // 发布评论
		text := c.Query("comment_text")
		PostComment(c, userID, text, uint(videoID))
	} else if actionType == "2" { //删除评论
		commentIDStr := c.Query("comment_id")
		commentID, _ := strconv.ParseInt(commentIDStr, 10, 10)
		DeleteComment(c, uint(videoID), uint(commentID))
	}
}

// PostComment post the comment
func PostComment(c *gin.Context, userID uint, text string, videoID uint) {
	//1. prepare the new comment
	newComment := model.Comment{
		VideoID: videoID,
		UserID:  userID,
		Content: text,
	}

	// 2. post a new comment and add the comment count for this video
	err1 := mysql.DB.Transaction(func(db *gorm.DB) error {
		if err := mysql.PostComment(&newComment); err != nil {
			return err
		}
		if err := mysql.AddCommentCount(videoID); err != nil {
			return err
		}
		return nil
	})
	// getUser, err2 := service.GetUser(userID)
	// 3. get the login user info and author id
	var getUser model.User
	err2 := mysql.GetAUserByID(userID, &getUser)
	authorID, err3 := mysql.GetVideoAuthorID(videoID)

	// 4. return error
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  "发布评论失败",
		})
		c.Abort()
		return
	}

	// 5. return the latest comment
	c.JSON(http.StatusOK, CommentActionResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "发布评论成功",
		},
		Comment: CommentResponse{
			ID:         newComment.ID,
			Content:    newComment.Content,
			CreateDate: newComment.CreatedAt.Format("01-02"),
			User: UserResponse{
				// ID:            getUser.ID,
				// 因为发布评论的人一定是当前登录用户，因此设成0，该用户可以删除刚发布的评论
				ID:              0,
				Name:            getUser.Name,
				FollowCount:     getUser.FollowCount,
				FollowerCount:   getUser.FollowerCount,
				Avatar:          getUser.Avatar,
				BackgroundImage: getUser.BackgroundImage,
				IsFollow:        service.IsFollowing(userID, authorID),
			},
		},
	})
}

// DeleteComment delete the comment
func DeleteComment(c *gin.Context, videoID uint, commentID uint) {
	// 1. delete the comment and reduce the comment count of video
	err := mysql.DB.Transaction(func(db *gorm.DB) error {
		if err := mysql.DeleteComment(commentID); err != nil {
			return err
		}
		if err := mysql.ReduceCommentCount(videoID); err != nil {
			return err
		}
		return nil
	})

	// 2. return error
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  "删除评论失败",
		})
		c.Abort()
		return
	}

	// 3. delete successfully
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  "成功删除评论",
	})
}

// CommentList get the comment list
func CommentList(c *gin.Context) {
	// 1. get the user id and video id
	getUserID, _ := c.Get("user_id")
	var userID uint
	if v, ok := getUserID.(uint); ok {
		userID = v
	}
	videoIDStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 10)

	// 2. get the comment list
	commentList, err := mysql.GetCommentList(uint(videoID))

	// 3. fail to get the comment list
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorCommentListGet.Error(),
		})
		c.Abort()
		return
	}

	// 4. pack the response
	var responseCommentList []CommentResponse
	for i := 0; i < len(commentList); i++ {
		// getUser, err1 := service.GetUser(commentList[i].UserID)
		var getUser model.User
		err1 := mysql.GetAUserByID(commentList[i].UserID, &getUser)

		if err1 != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  common.ErrorCommentListGet.Error(),
			})
			c.Abort()
			return
		}
		var returnID uint
		fmt.Printf("\n当前登录的用户id: %d\n\n", userID)
		fmt.Printf("\n评论者的id: %d\n\n", commentList[i].UserID)

		// if the comment user equals to the login user, change the id to
		// to trigger the delete action in app
		if userID == commentList[i].UserID {
			returnID = 0
		} else {
			returnID = commentList[i].UserID
		}

		responseComment := CommentResponse{
			ID:         commentList[i].ID,
			Content:    commentList[i].Content,
			CreateDate: commentList[i].CreatedAt.Format("01-02"), // mm-dd
			User: UserResponse{
				// ID:            getUser.ID,
				ID:              returnID,
				Name:            getUser.Name,
				FollowCount:     getUser.FollowCount,
				FollowerCount:   getUser.FollowerCount,
				Avatar:          getUser.Avatar,
				BackgroundImage: getUser.BackgroundImage,
				IsFollow:        service.IsFollowing(userID, commentList[i].ID),
			},
		}
		responseCommentList = append(responseCommentList, responseComment)

	}

	// 5. return the response
	c.JSON(http.StatusOK, CommentListResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "成功获取评论列表",
		},
		CommentList: responseCommentList,
	})
}
