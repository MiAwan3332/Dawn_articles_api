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

// Get : Return a news article by ID
// @Summary Return a news article by ID
// @Description Return a news article by ID
// @Tags news
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Article
// @Router /article/{id} [get]
func FetchArticle(c *gin.Context) {

	id := c.Param("id")

	var article models.Article

	cly := colly.NewCollector(
		colly.AllowedDomains("www.dawn.com", "dawn.com"),
	)

	cly.OnHTML("article", func(e *colly.HTMLElement) {
		title := e.ChildText(".story__title:first-child")
		link := e.ChildAttr(".story__link  ", "href")
		publishTime := e.ChildAttr(".timestamp--time", "title")
		story := e.ChildText(".story__content")
		layout := "2006-01-02T15:04:05+05:00"
		publishTimeParsed, _ := time.Parse(layout, publishTime)
		if len(publishTime) > 0 {
			article = models.Article{
				Title:       title,
				URL:         link,
				PublishTime: publishTimeParsed.Format("2006-01-02 03:04 PM"),
				Story:       story,
			}
		}
	})
	cly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	cly.Visit("https://www.dawn.com/news/" + id)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	c.JSON(http.StatusOK, article)
}
