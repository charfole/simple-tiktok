package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/charfole/simple-tiktok/config"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
	"github.com/tencentyun/cos-go-sdk-v5"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/gorm"
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

// COSUpload upload the file to the COS
func COSUpload(fileName string, reader io.Reader) (string, error) {
	// bucketURL := fmt.Sprintf(objectstorage.COS_URL_FORMAT, objectstorage.COS_BUCKET_NAME, objectstorage.COS_APP_ID, objectstorage.COS_REGION)
	bucketURL := fmt.Sprintf(config.Info.COS.URLFormat, config.Info.COS.BucketName, config.Info.COS.AppID, config.Info.COS.Region)
	fmt.Println("bucketURL: ", bucketURL)
	u, _ := url.Parse(bucketURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Info.COS.SecretID,
			SecretKey: config.Info.COS.SecretKey,
		},
	})
	// put(upload) the file to the COS
	_, err := client.Object.Put(context.Background(), fileName, reader, nil)
	if err != nil {
		panic(err)
	}
	// return "https://charfolebase-1301984140.cos.ap-guangzhou.myqcloud.com/" + fileName, nil
	return bucketURL + "/" + fileName, nil
}

// GetCoverFrame call the ffmpeg-go to get the cover
func GetCoverFrame(inFileName string, frameNum int) io.Reader {
	// create and write the cover image into buf
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}

func CreateVideo(userID uint, playURL, coverURL, title string) (err error) {
	video := model.Video{
		Model:         gorm.Model{},
		AuthorID:      userID,
		PlayURL:       playURL,
		CoverURL:      coverURL,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	err = mysql.CreateVideo(&video)
	return err
}
