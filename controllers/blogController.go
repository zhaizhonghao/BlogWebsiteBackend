package controllers

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zhaizhonghao/auth/database"
	"github.com/zhaizhonghao/auth/models"
)

func UploadBlog(c *fiber.Ctx) error {

	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}
	
	//创建blog
	newBlog := models.Blog{
		BlogTitle:    data["blogTitle"],
		BlogHTML: data["blogHTML"],
		BlogCoverPhotoPath: data["blogCoverPhotoPath"],
		BlogCoverPhotoName:data["blogCoverPhotoName"],
		Creator:data["creator"],
		CreatorId: data["creatorId"],
		CreateTime:data["createTime"],
	}

	//将博客插入数据库
	database.DB.Create(&newBlog)

	return c.JSON(newBlog)
}

func UploadImage(c *fiber.Ctx) error{

	file, err := c.FormFile("image")

    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message":   err.Error(),
        })
    }

	// Get Buffer from file
	buffer, err := file.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":   err.Error(),
		})
	}
	defer buffer.Close()


	//给图片随机生成一个名字
	rand.Seed(time.Now().UnixNano())
	// Go rune data type represent Unicode characters
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	fileName := make([]rune, 64)
	// creat a random slice of runes (characters) to create our emailVerPassword (random string of characters)
	for i := 0; i < 64; i++ {
		fileName[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	imageSuffix := strings.Split(filepath.Base(file.Filename), ".")[1]
	fileNameStr := filepath.Base(string(fileName))+"."+imageSuffix
	fmt.Println("random file name:", fileNameStr)

	//创建文件
	filePath := filepath.Join("C:/uploads", fileNameStr)
	dst, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, buffer); err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(fiber.Map{
		"message":"image upload successfully",
		"photoFilePath": "http://localhost:8081/"+fileNameStr,
	})
}

func GetAllBlogs(c *fiber.Ctx) error{
	blogs := []models.Blog{}
	database.DB.Find(&blogs)
	return c.JSON(blogs)
}

func DeleteBlogById(c *fiber.Ctx) error{
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	blogId := data["blogid"]

	//查看该blog在不在
	var blog models.Blog
	database.DB.Where("id = ?", blogId).First(&blog)
	if blog.Id == 0 {
		c.Status(fiber.StatusConflict)
		return c.JSON(fiber.Map{
			"message": "没有该blog，或许之前已被删除了！",
		})
	}

	//删除该博客
	database.DB.Delete(&models.Blog{}, blogId)

	return c.JSON(fiber.Map{
		"message":"delete the blog successfully",
	})
}

func UpdateBlog(c *fiber.Ctx) error{
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return err
	}

	blogId := data["id"]
	fmt.Println("updated blog id",blogId)

	//查看该blog在不在
	var blog models.Blog
	database.DB.Where("id = ?", blogId).First(&blog)
	if blog.Id == 0 {
		c.Status(fiber.StatusConflict)
		return c.JSON(fiber.Map{
			"message": "没有该blog，或许之前已被删除了！",
		})
	}

	//更新blog
	newBlog := &models.Blog{
		BlogTitle: data["blogTitle"],
		BlogHTML : data["blogHTML"],
		BlogCoverPhotoPath: data["blogCoverPhotoPath"],
		BlogCoverPhotoName:data["blogCoverPhotoName"],
		Creator:data["creator"],
		CreatorId:data["creatorId"],
		CreateTime:data["createTime"],
	}
	result := database.DB.Model(&models.Blog{}).Where("id = ?", blogId).Updates(newBlog)
	fmt.Println("changed row number",result.RowsAffected)
	if(result.Error != nil){
		fmt.Println(result.Error)
		return result.Error
	}

	return c.JSON(fiber.Map{
		"message":"update the blog successfully",
	})
}

