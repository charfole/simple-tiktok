package mysql

import (
	"time"

	"github.com/charfole/simple-tiktok/model"
)

// GetCommentList get the comment list of video
func GetCommentList(videoID uint) ([]model.Comment, error) {
	var commentList []model.Comment
	if err := DB.Model(model.Comment{}).Where("video_id=?", videoID).Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}

// PostComment post a new comment
func PostComment(comment *model.Comment) error {
	if err := DB.Model(model.Comment{}).Create(comment).Error; err != nil {
		return err
	}
	return nil
}

// DeleteComment delete a comment
func DeleteComment(commentID uint) error {
	if err := DB.Model(model.Comment{}).Where("id = ?", commentID).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}
