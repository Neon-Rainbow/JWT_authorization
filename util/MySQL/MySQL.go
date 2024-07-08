package MySQL

import (
	"JWT_authorization/config"
	"JWT_authorization/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var mySQL *gorm.DB

func InitMySQL() error {
	// Initialize MySQL

	appConfig := config.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		appConfig.MySQL.Username,
		appConfig.MySQL.Password,
		appConfig.MySQL.Host,
		appConfig.MySQL.Port,
		appConfig.MySQL.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("连接数据库失败: %v", err)
		return err
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Printf("数据库迁移失败: %v", err)
		return err
	}

	mySQL = db

	return nil
}

func GetMySQL() *gorm.DB {
	return mySQL
}
