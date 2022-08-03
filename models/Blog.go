package models

type Blog struct {
	Id                 uint   `json:"id"`
	BlogTitle          string `json:"blogTitle"`
	BlogHTML           string `json:"blogHTML"`
	BlogCoverPhotoPath string `json:"blogCoverPhotoPath"`
	BlogCoverPhotoName string `json:"blogCoverPhotoName"`
	Creator            string `json:"creator"`
	CreatorId          string `json:"creatorId"`
	CreateTime         string `json:"createTime"`
}
