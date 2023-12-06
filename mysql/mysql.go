package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitMysql() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/ginredis?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("connect mysql error", err)
		return
	}
	fmt.Println("mysql connect success")
	return nil
}
