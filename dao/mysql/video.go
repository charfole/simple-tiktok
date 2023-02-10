package mysql

import "github.com/charfole/simple-tiktok/model"

func GetVideoByTime(time string, videoNum int, videoList *[]model.Video) (err error) {
	// return the video list by desc
	err = DB.Model(&model.Video{}).Where("created_at < ?", time).
		Order("created_at desc").Limit(videoNum).Find(videoList).Error
	return
}
