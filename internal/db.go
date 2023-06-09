package internal

import (
	"dapp-backend/config"
	"dapp-backend/model"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() {
	var err error
	log.Println("connect to mysql:", config.Setting.MysqlConfig.Host)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Setting.MysqlConfig.UserName,
		config.Setting.MysqlConfig.PassWord,
		config.Setting.MysqlConfig.Host,
		config.Setting.MysqlConfig.Port,
		config.Setting.MysqlConfig.DataBase,
	)
	if DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Fatalln(err)
	}

	if err = DB.AutoMigrate(
		&model.UserVsion{},
	); err != nil {
		log.Fatalln("auto migrate db model with error:", err)
	}
}
