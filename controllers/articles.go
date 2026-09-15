package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FetchLatestNews returns the latest articles from Dawn's official RSS feed.
// @Summary Return a list of latest news articles
// @Description List all latest news articles published on Dawn.com
// @Tags news
// @Accept json
// @Produce json
// @Success 200 {object} models.ArticleDetails
// @Router /latest-news/ [get]
func FetchLatestNews(c *gin.Context) {
	feed, err := fetchLatestFeed(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Unable to fetch the Dawn news feed"})
		return
	}

	c.JSON(http.StatusOK, feed.articleDetails())
}
