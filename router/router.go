package router

import (
	"github.com/charfole/simple-tiktok/controller"
	"github.com/charfole/simple-tiktok/middleware"
	"github.com/charfole/simple-tiktok/mycontroller"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	// apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/feed/", mycontroller.Feed)
	// apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.GET("/user/", middleware.JWTMiddleware(), mycontroller.UserInfo)
	// apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/register/", mycontroller.UserRegister)
	// apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/user/login/", mycontroller.UserLogin)
	// apiRouter.POST("/publish/action/", controller.Publish)
	apiRouter.POST("/publish/action/", middleware.JWTMiddleware(), mycontroller.Publish)
	// apiRouter.GET("/publish/list/", controller.PublishList)
	apiRouter.GET("/publish/list/", middleware.JWTMiddleware(), mycontroller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction)
	apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.JWTMiddleware(), mycontroller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.JWTMiddleware(), mycontroller.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.JWTMiddleware(), mycontroller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", controller.MessageChat)
	apiRouter.POST("/message/action/", controller.MessageAction)
}
