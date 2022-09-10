package database

import (
	"github.com/zhaizhonghao/auth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(mysql.Open("root:123456@/blogwebsite"), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database!")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
	connection.AutoMigrate(&models.Blog{})
	connection.AutoMigrate(&models.Comment{})
	connection.AutoMigrate(&models.CommentResponse{})

}
