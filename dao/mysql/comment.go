package mysql

import (
	"time"

	"github.com/charfole/simple-tiktok/model"
)

// GetCommentList 获取指定videoID的评论表
func GetCommentList(videoID uint) ([]model.Comment, error) {
	var commentList []model.Comment
	if err := DB.Table("comments").Where("video_id=?", videoID).Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}

// PostComment 发布评论
func PostComment(comment model.Comment) error {
	if err := DB.Table("comments").Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

// DeleteComment 删除指定commentID的评论
func DeleteComment(commentID uint) error {
	if err := DB.Table("comments").Where("id = ?", commentID).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}
