package models

type ArticleDetails struct {
	Headline    string `json:"headline"`
	URL         string `json:"url"`
	PublishTime string `json:"timestamp"`
	Excerpt     string `json:"excerpt"`
}
