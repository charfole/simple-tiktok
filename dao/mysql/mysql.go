package mysql

import (
	"fmt"

	"github.com/charfole/simple-tiktok/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func getDSN() string {
	arg := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%v&loc=%s",
		config.Info.DB.Username, config.Info.DB.Password, config.Info.DB.Host, config.Info.DB.Port, config.Info.DB.Database,
		config.Info.DB.Charset, config.Info.DB.ParseTime, config.Info.DB.Loc)
	fmt.Println(arg)
	return arg
}

func InitMySQL() {
	var DBError error
	// 主函数中加载了配置文件 dao层则不需要再加载
	// config.InitEnv()
	dsn := getDSN()
	DB, DBError = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if DBError != nil {
		panic(DBError)
	}
	fmt.Println(DB.Name())
}
