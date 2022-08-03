package models

type Entry struct {
	Id    uint   `json:"id"`
	Email string `json:"email" gorm:"unique"`
}
