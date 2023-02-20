package mysql

import (
	"github.com/charfole/simple-tiktok/model"
	"gorm.io/gorm"
)

// CreateVideo 添加一条视频信息
func CreateVideo(video *model.Video) (err error) {
	err = DB.Model(&model.Video{}).Create(&video).Error
	return err
}

func GetVideoByTime(time string, videoNum int, videoList *[]model.Video) (err error) {
	// return the video list by desc
	err = DB.Model(&model.Video{}).Where("created_at < ?", time).
		Order("created_at desc").Limit(videoNum).Find(videoList).Error
	return
}

// GetVideoList 根据用户id查找 所有与该用户相关视频信息
func GetVideoList(userID uint) []model.Video {
	var videoList []model.Video
	DB.Model(model.Video{}).Where("author_id=?", userID).Find(&videoList)
	return videoList
}

// GetVideoAuthor get video author
func GetVideoAuthorID(videoID uint) (uint, error) {
	var video model.Video
	if err := DB.Table("videos").Where("id = ?", videoID).Find(&video).Error; err != nil {
		return video.ID, err
	}
	return video.AuthorID, nil
}

func AddVideoFavoriteCount(videoID uint) (err error) {
	err = DB.Model(model.Video{}).Where("id = ?", videoID).
		Update("favorite_count", gorm.Expr("favorite_count + 1")).
		Error

	return
}

func ReduceVideoFavoriteCount(videoID uint) (err error) {
	err = DB.Model(model.Video{}).Where("id = ?", videoID).
		Update("favorite_count", gorm.Expr("favorite_count - 1")).
		Error

	return
}
