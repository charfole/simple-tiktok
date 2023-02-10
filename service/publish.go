package service

import (
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
)

type ReturnAuthor struct {
	AuthorID      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type ReturnVideo struct {
	VideoID       uint         `json:"video_id"`
	Author        ReturnAuthor `json:"author"`
	PlayURL       string       `json:"play_url"`
	CoverURL      string       `json:"cover_url"`
	FavoriteCount uint         `json:"favorite_count"`
	CommentCount  uint         `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

func PackAuthor(user model.User, hostID uint, guestID uint) (returnAuthor ReturnAuthor) {
	returnAuthor = ReturnAuthor{
		AuthorID:      user.ID,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      IsFollowing(hostID, guestID),
	}
	return
}

func PackVideo(videoList []model.Video, author ReturnAuthor, hostID uint) (returnVideoList []ReturnVideo) {
	for i := 0; i < len(videoList); i++ {
		returnVideo := ReturnVideo{
			VideoID:       videoList[i].ID,
			Author:        author,
			PlayURL:       videoList[i].PlayURL,
			CoverURL:      videoList[i].CoverURL,
			FavoriteCount: videoList[i].FavoriteCount,
			CommentCount:  videoList[i].CommentCount,
			IsFavorite:    mysql.IsFavorite(hostID, videoList[i].ID),
			Title:         videoList[i].Title,
		}
		returnVideoList = append(returnVideoList, returnVideo)
	}
	return
}
