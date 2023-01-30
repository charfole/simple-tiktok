package main

import (
	"github.com/charfole/simple-tiktok/config"
	"github.com/charfole/simple-tiktok/dao"
	"github.com/charfole/simple-tiktok/router"
	"github.com/gin-gonic/gin"
)

var InitError error

func main() {
	// go service.RunMessageServer()

	r := gin.Default()

	config.InitEnv()
	dao.InitMySQL()
	// dao.DB.AutoMigrate(models.Todo{})
	// todo := models.Todo{ID: 1, Title: "title", Date: "date", Status: true}
	// err := models.CreateATodo(&todo)
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
