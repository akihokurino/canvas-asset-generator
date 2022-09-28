package model

type Frame struct {
	ID              string `json:"id"`
	WorkID          string `json:"workId"`
	OrgImageUrl     string `json:"orgImageUrl"`
	ResizedImageUrl string `json:"resizedImageUrl"`
	ImageGsPath     string `json:"imageGsPath"`
}
