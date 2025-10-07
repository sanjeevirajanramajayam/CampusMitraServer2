package models

type Project struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Github      string   `json:"github"`
	Stack       []string `json:"stack"`
}