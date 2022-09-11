package database

import (
	"github.com/zhaizhonghao/auth/models"
	"github.com/zhaizhonghao/auth/staticConfigs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := staticConfigs.DATABASE_USERNAME+":"+staticConfigs.DATABASE_PASWORD+"@/"+staticConfigs.DATABASE_NAME
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database!")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
	connection.AutoMigrate(&models.Blog{})
	connection.AutoMigrate(&models.Comment{})
	connection.AutoMigrate(&models.CommentResponse{})

}
