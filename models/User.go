package models

type User struct {
	Id       uint   `json:"id"`
	Email    string `json:"email" gorm:"unique"`
	Username string `json:"username"`
	Password []byte `json:"-"`
	Role     string `json:"role"`
	EmailVerPWhash string  `json:"evpw"`
	Timeout string  `json:"timeout"`
}
