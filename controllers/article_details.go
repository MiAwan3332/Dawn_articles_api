package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abdulalikhan/Dawn-News-API/models"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

// List : Return a list of latest news articles
// @Summary Return a list of latest news articles
// @Description List all latest news articles published on Dawn.com
// @Tags news
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ArticleDetails
// @Router /latest-news/ [get]
func FetchLatestNews(c *gin.Context) {
	articles := make([]models.ArticleDetails, 0)

	cly := colly.NewCollector(
		colly.AllowedDomains("www.dawn.com", "dawn.com"),
	)

	cly.OnHTML("article", func(e *colly.HTMLElement) {
		headline := e.ChildText(".story__link  ")
		link := e.ChildAttr(".story__link  ", "href")
		publishTime := e.ChildAttr(".timestamp--time", "title")
		excerpt := e.ChildText(".story__excerpt")
		layout := "2006-01-02T15:04:05+05:00"
		publishTimeParsed, _ := time.Parse(layout, publishTime)
		if len(publishTime) > 0 {
			article := models.ArticleDetails{
				Headline:    headline,
				URL:         link,
				PublishTime: publishTimeParsed.Format("2006-01-02 03:04 PM"),
				Excerpt:     excerpt,
			}
			articles = append(articles, article)
		}
	})
	cly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	cly.Visit("https://www.dawn.com/latest-news")
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	//enc.Encode(articles)
	c.JSON(http.StatusOK, articles)
}
