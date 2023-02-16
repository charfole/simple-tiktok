package mysql

import "github.com/charfole/simple-tiktok/model"

// CreateVideo 添加一条视频信息
func CreateVideo(video *model.Video) {
	DB.Model(&model.Video{}).Create(&video)
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
