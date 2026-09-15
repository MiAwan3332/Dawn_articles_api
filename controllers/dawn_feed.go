package controllers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/abdulalikhan/Dawn-News-API/models"
	"golang.org/x/net/html"
)

const dawnLatestFeedURL = "https://www.dawn.com/feeds/latest-news"

var feedHTTPClient = &http.Client{Timeout: 30 * time.Second}

type dawnFeed struct {
	Channel struct {
		Items []dawnFeedItem `xml:"item"`
	} `xml:"channel"`
}

type dawnFeedItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Content     string `xml:"encoded"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Media       struct {
		URL string `xml:"url,attr"`
	} `xml:"content"`
}

func fetchLatestFeed(ctx context.Context) (dawnFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dawnLatestFeedURL, nil)
	if err != nil {
		return dawnFeed{}, err
	}
	req.Header.Set("User-Agent", "LCA-Dawn-Articles-API/1.0")

	response, err := feedHTTPClient.Do(req)
	if err != nil {
		return dawnFeed{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return dawnFeed{}, fmt.Errorf("Dawn feed returned HTTP %d", response.StatusCode)
	}

	var feed dawnFeed
	if err := xml.NewDecoder(io.LimitReader(response.Body, 5<<20)).Decode(&feed); err != nil {
		return dawnFeed{}, err
	}
	return feed, nil
}

func (feed dawnFeed) articleDetails() []models.ArticleDetails {
	articles := make([]models.ArticleDetails, 0, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		articles = append(articles, models.ArticleDetails{
			Headline:    strings.TrimSpace(item.Title),
			URL:         strings.TrimSpace(item.Link),
			ImageUrl:    strings.TrimSpace(item.Media.URL),
			PublishTime: formatPublishTime(item.PubDate),
			Excerpt:     truncate(cleanArticleText(item.Description), 300),
		})
	}
	return articles
}

func (feed dawnFeed) article(id string) (models.Article, bool) {
	needle := "/news/" + strings.TrimSpace(id)
	for _, item := range feed.Channel.Items {
		if !strings.Contains(item.Link, needle) && !strings.HasSuffix(item.GUID, needle) {
			continue
		}

		story := item.Content
		if strings.TrimSpace(story) == "" {
			story = item.Description
		}
		return models.Article{
			Title:       strings.TrimSpace(item.Title),
			URL:         strings.TrimSpace(item.Link),
			PublishTime: formatPublishTime(item.PubDate),
			Story:       cleanArticleText(story),
		}, true
	}
	return models.Article{}, false
}

func cleanArticleText(value string) string {
	document, err := html.Parse(strings.NewReader(value))
	if err != nil {
		return strings.Join(strings.Fields(value), " ")
	}

	parts := make([]string, 0)
	var collectText func(*html.Node)
	collectText = func(node *html.Node) {
		if node.Type == html.ElementNode {
			switch node.Data {
			case "script", "style", "noscript":
				return
			}
		}
		if node.Type == html.TextNode {
			if text := strings.TrimSpace(node.Data); text != "" {
				parts = append(parts, text)
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			collectText(child)
		}
	}
	collectText(document)

	return strings.Join(strings.Fields(strings.Join(parts, " ")), " ")
}

func truncate(value string, maxRunes int) string {
	if utf8.RuneCountInString(value) <= maxRunes {
		return value
	}
	runes := []rune(value)
	return strings.TrimSpace(string(runes[:maxRunes])) + "..."
}

func formatPublishTime(value string) string {
	publishedAt, err := time.Parse(time.RFC1123Z, strings.TrimSpace(value))
	if err != nil {
		return strings.TrimSpace(value)
	}
	return publishedAt.Format("2006-01-02 03:04 PM")
}
