package models

type Article struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	PublishTime string `json:"timestamp"`
	Story       string `json:"story"`
}
