package models

type Hackathon struct{
	Title string `json:"title"`
	Place string `json:"place"`
	ImgUrl string `json:"img_url,omitempty"`
}