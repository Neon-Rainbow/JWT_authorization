package MySQL

import (
	"JWT_authorization/config"
	"JWT_authorization/model"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var mySQL *gorm.DB

func InitMySQL() error {
	appConfig := config.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		appConfig.MySQL.Username,
		appConfig.MySQL.Password,
		appConfig.MySQL.Host,
		appConfig.MySQL.Port,
		appConfig.MySQL.Database,
	)
	//fmt.Printf("dsn: %v\n", dsn)

	var err error
	mySQL, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("连接数据库失败: %v", err)
		return err
	}

	// 设置一个 20 秒的上下文超时
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 使用 AutoMigrate 进行数据库迁移操作
	err = mySQL.WithContext(ctx).AutoMigrate(&model.User{})
	if err != nil {
		log.Printf("数据库迁移失败: %v", err)
		return err
	}

	return nil
}

func GetMySQL() *gorm.DB {
	if mySQL == nil {
		log.Fatal("Database not initialized")
	}
	return mySQL
}
