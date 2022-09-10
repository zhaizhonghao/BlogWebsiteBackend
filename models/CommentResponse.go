package models

type CommentResponse struct {
	Id          uint   `json:"id"`
	BlogId      uint   `json:"blogId"`
	CommentId   uint   `json:"commentId"`
	From        string `json:"from"`
	To          string `json:"to"`
	Timestamp   string `json:"timestamp"`
	MainContent string `json:"mainContent"`
}
