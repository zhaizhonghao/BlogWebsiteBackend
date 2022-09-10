package models

type Comment struct {
	Id          uint              `json:"id"`
	BlogId      uint              `json:"blogId"`
	Username    string            `json:"userName"`
	Timestamp   string            `json:"timestamp"`
	MainContent string            `json:"mainContent"`
	Responses   []CommentResponse `json:"responses"`
}