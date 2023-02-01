package main

import (
	"github.com/charfole/simple-tiktok/config"
	"github.com/charfole/simple-tiktok/dao/mysql"
	"github.com/charfole/simple-tiktok/router"
	"github.com/gin-gonic/gin"
)

var InitError error

func main() {
	// go service.RunMessageServer()

	r := gin.Default()

	config.InitEnv()
	mysql.InitMySQL()
	// 用于测试mysql和gorm功能是否正常
	// mysql.DB.AutoMigrate(model.Todo{})
	// todo := model.Todo{ID: 1, Title: "title", Date: "date", Status: true}
	// err := mysql.CreateATodo(&todo)
	// if err != nil {
	// 	panic(err)
	// }
	router.InitRouter(r)

	port := ":" + config.Info.Server.Port
	InitError = r.Run(port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if InitError != nil {
		panic(InitError)
	}

}
