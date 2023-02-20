package mycontroller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/charfole/simple-tiktok/common"
	"github.com/charfole/simple-tiktok/config"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/model"
	"github.com/charfole/simple-tiktok/service"
	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

type VideoListResponse struct {
	common.Response
	VideoList []service.ReturnVideo `json:"video_list"`
}

func Publish(c *gin.Context) { //上传视频方法
	// 1. check the token and get user_id(the id of the author)
	getUserID, _ := c.Get("user_id")
	userID := getUserID.(uint)

	// 2. get the title and data of the uploaded video
	title := c.PostForm("title")
	data, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorVideoUpload.Error(),
		})
		return
	}

	// 3. construct the file name
	fileName := filepath.Base(data.Filename)
	nowTime := time.Now().Unix()
	// format: id_time_filename.mp4
	finalName := fmt.Sprintf("%d_%d_%s", userID, nowTime, fileName)

	fmt.Println("初始文件名：", fileName)
	fmt.Println("最终文件名：", finalName)

	// 4. open and read the local file
	f, err := data.Open()
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorVideoUpload.Error(),
		})
		return
	}
	defer f.Close()

	// 5. upload the local file to COS
	palyURL, err := service.COSUpload(finalName, f)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorCOSUpload.Error(),
		})
		return
	}

	// 6. save the file temporarily to the static path
	savePath := filepath.Join(config.Info.Path.StaticSourcePath, "/", finalName)
	fmt.Println("保存到服务器的临时路径：", savePath)

	if err := c.SaveUploadedFile(data, savePath); err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorVideoUpload.Error(),
		})
		return
	}

	// 7. get the cover from the local video
	coverName := strings.Replace(finalName, ".mp4", ".jpeg", 1)
	img := service.GetCoverFrame(savePath, 2) //获取第2帧封面

	// 8. upload the cover image to the COS
	coverURL, err := service.COSUpload(coverName, img)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorCOSUpload.Error(),
		})
		return
	}

	// 9. remove the local video
	err = os.Remove(savePath) // ignore_security_alert
	if err != nil {
		logging.Info(err)
	}

	// 10. save the video record to "videos" database
	err = service.CreateVideo(userID, palyURL, coverURL, title)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  common.ErrorVideoDBCreateFalse.Error(),
		})
		return
	}

	// 11. no errors found, upload successfully
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " --uploaded successfully",
	})
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
