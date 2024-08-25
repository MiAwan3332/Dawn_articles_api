package models

type ArticleDetails struct {
	Headline    string `json:"headline"`
	URL         string `json:"url"`
	ImageUrl    string `json:"imageurl"`
	PublishTime string `json:"timestamp"`
	Excerpt     string `json:"excerpt"`
}
