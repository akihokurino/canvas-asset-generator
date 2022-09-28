package model

type Frame struct {
	ID          string `json:"id"`
	WorkID      string `json:"workId"`
	ImageUrl    string `json:"imagePath"`
	ImageGsPath string `json:"imageGsPath"`
}
