package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zhaizhonghao/auth/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
	app.Get("/api/user/all", controllers.GetAllUsers)
	app.Post("/api/user/delete", controllers.DeleteUser)
	app.Post("/api/forgotPassword", controllers.ForgotPassword)
	app.Post("/api/resetPassword", controllers.ResetPassword)
	app.Post("/api/updateUserName", controllers.UpdateUserName)

	app.Post("/api/uploadBlog", controllers.UploadBlog)
	app.Post("/api/uploadImage", controllers.UploadImage)
	app.Get("/api/getAllBlogs",controllers.GetAllBlogs)
	app.Post("/api/deleteBlogById",controllers.DeleteBlogById)
	app.Post("/api/updateBlog",controllers.UpdateBlog)

}
