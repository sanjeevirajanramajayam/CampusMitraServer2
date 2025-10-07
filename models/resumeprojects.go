package models

type Projects struct {
	Title        string   `json:"title"`
	Descrription string   `json:"description"`
	Stack        []string `json:"stack"`
	Github       string   `json:"github"`
}