package mycontroller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/config"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
	"github.com/charfole/simple-tiktok/service"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VideoListResponse struct {
	common.Response
	VideoList []service.ReturnVideo `json:"video_list"`
}

func Publish(c *gin.Context) { //上传视频方法
	//1.中间件验证token后，获取userID
	getUserID, _ := c.Get("user_id")
	var userID uint
	userID = getUserID.(uint)

	//2.接收请求参数信息
	title := c.PostForm("title")
	data, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorVideoUpload.Error(),
		})
		return
	}

	//3.返回至前端页面的展示信息
	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userID, fileName)
	fmt.Println(fileName)
	fmt.Println(finalName)

	//先存储到本地文件夹，再保存到云端，获取封面后最后删除
	saveFile := filepath.Join(config.Info.Path.StaticSourcePath, "/", finalName)
	fmt.Println(saveFile)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorVideoUpload.Error(),
		})
		return
	}
	// 打开本地文件并读入
	f, err := data.Open()
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorVideoUpload.Error(),
		})
		return
	}
	//从本地上传到云端，并获取云端地址
	playUrl, err := service.COSUpload(finalName, f)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  finalName + "--uploaded successfully",
	})

	//直接传至云端，不用存储到本地
	coverName := strings.Replace(finalName, ".mp4", ".jpeg", 1)
	img := service.ExampleReadFrameAsJpeg(saveFile, 2) //获取第2帧封面

	coverUrl, err := service.COSUpload(coverName, img)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//删除保存在本地中的视频
	err = os.Remove(saveFile) // ignore_security_alert
	if err != nil {
		logging.Info(err)
	}

	//4.保存发布信息至数据库,刚开始发布，喜爱和评论默认为0
	video := model.Video{
		Model:         gorm.Model{},
		AuthorID:      userID,
		PlayURL:       playUrl,
		CoverURL:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	mysql.CreateVideo(&video)
}

func PublishList(c *gin.Context) { //获取列表的方法
	// 1. get the current user id by JWT
	rawHostID, _ := c.Get("user_id")
	hostID := rawHostID.(uint)

	// 2. get the guest id
	rawGuestID := c.Query("user_id")
	id, _ := strconv.Atoi(rawGuestID)
	guestID := uint(id)

	// if guestID equals to 0, it means the guest is current user
	if guestID == 0 {
		guestID = hostID
	}
	fmt.Println("hostID: ", hostID)
	fmt.Println("guestID: ", guestID)

	// 3. get the guest
	var user model.User
	err := mysql.GetAUserByID(guestID, &user)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorUserNotFound.Error(),
		})
		c.Abort()
		return
	}

	// 4. pack the guest
	returnAuthor := service.PackAuthor(user, hostID, guestID)

	// 5. get the video list of this guest(author), pack it and return
	videoList := mysql.GetVideoList(guestID)
	if len(videoList) > 0 {
		// videolist found, pack the videos and return
		var returnVideoList []service.ReturnVideo
		returnVideoList = service.PackVideo(videoList, returnAuthor, hostID)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: common.Response{
				StatusCode: 0,
				StatusMsg:  "success",
			},
			VideoList: returnVideoList,
		})
	} else {
		// video list not found, return error
		c.JSON(http.StatusOK, VideoListResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  common.ErrorVideoList.Error(),
			},
			VideoList: nil,
		})
	}
}
