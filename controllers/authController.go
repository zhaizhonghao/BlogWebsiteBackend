package controllers

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/zhaizhonghao/auth/database"
	"github.com/zhaizhonghao/auth/models"
	"golang.org/x/crypto/bcrypt"

	"github.com/zhaizhonghao/auth/staticConfigs"
)

const SecretKey = "zhaizhonghao de secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	//To check whether the email has been registered
	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id != 0 {
		c.Status(fiber.StatusConflict)
		return c.JSON(fiber.Map{
			"message": "邮箱已被注册！",
		})
	}

	//To check whether the username has been registered
	database.DB.Where("username = ?", data["username"]).First(&user)

	if user.Id != 0 {
		c.Status(fiber.StatusConflict)
		return c.JSON(fiber.Map{
			"message": "用户名已被注册！",
		})
	}

	//To check whether the password is correct 
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	newUser := models.User{
		Email:    data["email"],
		Username: data["username"],
		Password: password,
	}

	//assgin the role for the user
	if data["email"] == "390930230@qq.com" {
		newUser.Role = "admin"
	} else {
		newUser.Role = "user"
	}

	//register the user info in the database
	database.DB.Create(&newUser)

	return c.JSON(newUser)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "该邮箱还未注册，请注册！",
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"]))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "密码不对!",
		})
	}

	//To return the JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 dat
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "无法登录!",
		})
	}

	userToken := models.UserToken{
		User: user,
		Token: token,
	}

	return c.JSON(userToken)
}

func ForgotPassword(c *fiber.Ctx) error{

	//解析发送过来的数据
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	//验证邮箱是否注册过
	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "邮箱还未被注册！",
		})
	}

	//生成随机字符串
	// create timeout limit
	now := time.Now()
	// add 45 minutes
	timeout := now.Add(time.Minute * 45)
	fmt.Println(timeout)
	var timeLayoutStr = "2006-01-02 15:04:05"
	timeoutString := timeout.Format(timeLayoutStr)
	
	// create random code for email
	rand.Seed(time.Now().UnixNano())
	// Go rune data type represent Unicode characters
	var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	emailVerRandRune := make([]rune, 64)
	// creat a random slice of runes (characters) to create our emailVerPassword (random string of characters)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	fmt.Println("emailVerRandRune:", emailVerRandRune)
	emailVerPassword := string(emailVerRandRune)
	fmt.Println("emailVerPassword:", emailVerPassword)
	var emailVerPWhash []byte
	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	emailVerPWhash, err = bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		return err
	}
	fmt.Println("emailVerPWhash:", emailVerPWhash)
	emailVerPWhashStr := string(emailVerPWhash)
	fmt.Println("emailVerPWhashStr:", emailVerPWhashStr)

	//更新数据库中用户的邮箱验证哈希值
	newUser := &models.User{
		EmailVerPWhash:emailVerPWhashStr,
		Timeout:timeoutString,
	}
	result := database.DB.Model(&models.User{}).Where("email = ?", data["email"]).Updates(newUser)
	fmt.Println("changed row number",result.RowsAffected)
	if(result.Error != nil){
		fmt.Println(result.Error)
		return result.Error
	}

	//发送邮件
	from := "390930230@qq.com" //ex: "John.Doe@gmail.com"
	password := staticConfigs.SMTP_PASSWORD   // ex: "ieiemcjdkejspqz"
	// receiver address privided through toEmail argument
	to := []string{data["email"]}
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.qq.com"
	port := "587"
	address := host + ":" + port
	// message
	subject := "Subject: 翟中豪个人网站的账号恢复\n"
	// localhost:8080 will be removed by many email service but works with online sites
	// https must be used since we are sending personal data through url parameters
	body := "<body>请点击<a rel=\"nofollow noopener noreferrer\" target=\"_blank\" href=\"https://www.mywebsite.com/forgotpwchange?email=" + data["email"] + "&evpw=" + emailVerPassword + "\">重置密码</a>链接进行密码重置，该链接的有效时间为45分钟。</body>"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(subject + mime + body)
	// athentication data
	// func PlainAuth(identity, username, password, host string) Auth
	auth := smtp.PlainAuth("", from, password, host)
	// func SendMail(addr string, a Auth, from string, to []string, msg []byte) error
	fmt.Println("message:", string(message))
	err = smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		fmt.Println("error sending reset password email, err:", err)

		return err
	}
	return nil
}

func User(c *fiber.Ctx) error {
	//解析发送过来的数据
	fmt.Println("In User")
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(
		data["token"],
		&jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated!",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User
	database.DB.Where("id=?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func ResetPassword(c *fiber.Ctx) error{
	//解析发送过来的数据
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}

	//验证邮箱是否注册过
	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "邮箱还未被注册！",
		})
	}
	fmt.Println(user)


	//检查邮箱验证码有没有过期
	//1.获取email对应的Timeout
	timeoutString := user.Timeout
	//2.看有没有过期
	currentTime := time.Now()
	var timeLayoutStr = "2006-01-02 15:04:05"
	timeout, _ := time.Parse(timeLayoutStr, timeoutString)

	if currentTime.After(timeout){
		fmt.Println("current time",currentTime)
		fmt.Println("timeout",timeout)
		fmt.Println("the email verification password has been outdated")
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "邮箱修改验证码已过期",
		})
	}

	//检查邮箱验证码是不是和之前生成的一致
	// check if db ver_hash is the same as the hash of emailVerPassword from email
	fmt.Println("evpw",data["evpw"])
	err = bcrypt.CompareHashAndPassword( []byte(user.EmailVerPWhash),[]byte(data["evpw"]))
	if err != nil {
		fmt.Println("dbEmailVerHash and hash of emailVerPassword are not the same")
		return err
	}
	fmt.Println("dbEmailVerHash and hash of emailVerPassword are the same")

	//更新密码
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	newUser := &models.User{
		Password: password,
	}
	result := database.DB.Model(&models.User{}).Where("email = ?", data["email"]).Updates(newUser)
	fmt.Println("changed row number",result.RowsAffected)
	if(result.Error != nil){
		fmt.Println(result.Error)
		return result.Error
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Logout(c *fiber.Ctx) error {
	//to remove the cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func GetAllUsers(c *fiber.Ctx) error {
	users := []models.User{}
	database.DB.Find(&users)
	return c.JSON(users)
}

func DeleteUser(c *fiber.Ctx) error {
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}
	user := models.User{}

	database.DB.Where("email=?", data["email"]).Delete(&user)

	return c.JSON(user)
}

func UpdateUserName(c *fiber.Ctx) error {
	//从请求中获取邮箱和用户名
	var data map[string]string

	err := c.BodyParser(&data)

	if err != nil {
		return err
	}
	//检查邮箱是否已经注册
	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "邮箱还未被注册！",
		})
	}

	//更新用户名
	newUser := &models.User{
		Username: data["username"],
	}
	result := database.DB.Model(&models.User{}).Where("email = ?", data["email"]).Updates(newUser)
	fmt.Println("username change successfully! changed row number",result.RowsAffected)
	if(result.Error != nil){
		fmt.Println(result.Error)
		return result.Error
	}

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
