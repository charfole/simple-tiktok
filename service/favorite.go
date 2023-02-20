package service

import (
	"fmt"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
)

// FavoriteAction handle the favorite and disfavor action
func FavoriteAction(userID uint, videoID uint, actionType uint) (err error) {
	// 1. if the action is favor
	if actionType == 1 {
		favoriteAction := model.Favorite{
			UserID:  userID,
			VideoID: videoID,
			State:   1,
		}
		favoriteStruct := model.Favorite{}
		// 2. check the user favors the video or not
		err := mysql.IsFavoriteRecordExist(userID, videoID, &favoriteStruct)

		// 3. if not found
		if err != nil {
			// 4. create a new favor record from "favoriteAction" struct
			if err := mysql.CreateAFavorite(&favoriteAction); err != nil { //创建记录
				return err
			}

			// 5. add the favorite count of video
			if err := mysql.AddVideoFavoriteCount(videoID); err != nil {
				return err
			}

			// 6. add the favorite count of user
			if err := mysql.AddFavoriteCount(userID); err != nil {
				return err
			}

			// 7. get the authorID and add the total favorited
			AuthorID, err := mysql.GetVideoAuthorID(videoID)
			if err != nil {
				return err
			}
			if err := mysql.AddTotalFavorited(AuthorID); err != nil {
				return err
			}
		} else {
			// 1. if the action is disfavor and the favor record exists with "state" tag equals to 0
			if favoriteStruct.State == 0 {
				// 2. add the favorite count of video
				if err := mysql.AddVideoFavoriteCount(videoID); err != nil {
					return err
				}

				// 3. update the state of this favorite record from 0 to 1
				if err := mysql.UpdateFavoriteState(userID, videoID, 1); err != nil {
					return err
				}

				// 4. add the favorite count of user
				if err := mysql.AddFavoriteCount(userID); err != nil {
					return err
				}

				// 5. get the authorID and add the total favorited
				AuthorID, err := mysql.GetVideoAuthorID(videoID)
				if err != nil {
					return err
				}
				if err := mysql.AddTotalFavorited(AuthorID); err != nil {
					return err
				}
			} else if favoriteStruct.State == 1 {
				// 1. the user already favors the video, return error
				fmt.Printf("\n喜欢已存在\n")
				return common.ErrorFavoriteExist
			}
		}
	} else if actionType == 2 {
		// 1. if the action is disfavor
		var favoriteStruct model.Favorite
		// 2. if the favorite record not found, return error
		if err := mysql.IsFavoriteRecordExist(userID, videoID, &favoriteStruct); err != nil {
			return common.ErrorFavoriteNotFound
		}
		// 3. record found, handle the disfavor action
		if favoriteStruct.State == 1 {
			// 4. reduce the favorite count of video
			if err := mysql.ReduceVideoFavoriteCount(videoID); err != nil {
				return err
			}
			// 5. update the state of this favorite record from 1 to 0
			if err := mysql.UpdateFavoriteState(userID, videoID, 0); err != nil {
				return err
			}
			// 6. reduce the favorite count of user
			if err := mysql.ReduceFavoriteCount(userID); err != nil {
				return err
			}
			// 7. get the authorID and reduce the total favorited
			AuthorID, err := mysql.GetVideoAuthorID(videoID)
			if err != nil {
				return err
			}
			if err := mysql.ReduceTotalFavorited(AuthorID); err != nil {
				return err
			}
			return err
		} else if favoriteStruct.State == 0 {
			// 1. the user already disfavors the video, return error
			fmt.Printf("\n取消喜欢已存在\n")
			return common.ErrorDisfavorExist
		}
	} else {
		return common.ErrorUnknownAction
	}
	return nil
}

// FavoriteList return the favorite list of user
func FavoriteList(userID uint) ([]model.Video, error) {
	// var favoriteList []model.Favorite
	videoList := make([]model.Video, 0)
	// if err = mysql.DB.Table("favorites").Where("user_id=? AND state=?", userID, 1).Find(&favoriteList).Error; err != nil {
	// 	return videoList, nil
	// }
	favoriteList, err := mysql.GetFavoriteList(userID)
	if err != nil {
		// empty favorite list
		return videoList, nil
	}
	// get the id of favorited video and get the video
	for _, m := range favoriteList {
		var video model.Video
		// if err := mysql.DB.Table("videos").Where("id=?", m.VideoID).Find(&video).Error; err != nil {
		// 	return nil, err
		// }
		if err := mysql.GetVideoByID(m.VideoID, &video); err != nil {
			// return nil, err
			// video not found, maybe deleted, continue to find other favorited videos
			fmt.Printf("%d号视频查找失败\n", m.VideoID)
			continue
		} else {
			videoList = append(videoList, video)
		}
	}
	return videoList, nil
}
