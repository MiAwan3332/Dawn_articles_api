package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FetchArticle returns an article from Dawn's official RSS feed by news ID.
// @Summary Return a news article by ID
// @Description Return a news article by ID
// @Tags news
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} models.Article
// @Router /article/{id} [get]
func FetchArticle(c *gin.Context) {
	feed, err := fetchLatestFeed(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Unable to fetch the Dawn news feed"})
		return
	}

	article, found := feed.article(c.Param("id"))
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article is not present in the latest-news feed"})
		return
	}

	c.JSON(http.StatusOK, article)
}
