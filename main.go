package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gosimple/slug"
)

type Article struct {
	Guid        ArticleGuid `xml:"guid"`
	Link        string      `xml:"link"`
	Title       string      `xml:"title"`
	Description string      `xml:"description"`
	Author      Author      `xml:"author"`
	CreatedAt   time.Time   `xml:"pubDate"`
}

type ArticleGuid struct {
	Id          string `xml:",chardata"`
	IsPermaLink bool   `xml:"isPermaLink,attr"`
}

type Author struct {
	Name    string `xml:"name"`
	Profile string `xml:"profile"`
}

type Feed struct {
	XMLName     xml.Name   `xml:"channel"`
	Title       string     `xml:"title"`
	Description string     `xml:"description"`
	Link        string     `xml:"link"`
	Articles    []*Article `xml:"item"`
}

var ArticleFeed Feed = Feed{
	Articles:    []*Article{},
	Title:       "ICT-Flex - Announcements",
	Link:        "https://ict-flex.nl/category/mededelingen/",
	Description: "RSS Feed for ICT-Flex announcements",
}

func main() {
	c := colly.NewCollector()

	c.OnHTML("div.entries", func(e *colly.HTMLElement) {
		e.ForEach("article", func(_ int, e *colly.HTMLElement) {
			var article Article = Article{
				Guid: ArticleGuid{
					IsPermaLink: false,
				},
			}

			article.Guid.Id = e.Attr("id")
			article.Author.Name = e.ChildText(".entry-meta>.meta-author>a>span")
			article.Author.Profile = e.ChildAttr(".entry-meta>.meta-author>a", "href")

			var date string = e.ChildAttr(".entry-meta>.meta-date>time", "datetime")
			article.CreatedAt, _ = time.Parse(time.RFC3339, date)

			ArticleFeed.Articles = append(ArticleFeed.Articles, &article)

			e.Request.Visit(e.ChildAttr(".entry-title>a", "href"))
		})
	})

	c.OnHTML("article.type-post", func(e *colly.HTMLElement) {
		if hero := e.ChildText("div.hero-section"); hero == "" {
			return
		}

		var article *Article

		for _, art := range ArticleFeed.Articles {
			if art.Guid.Id != e.Attr("id") {
				continue
			}

			article = art
		}

		article.Title = e.ChildText(".hero-section>.entry-header>.page-title")
		article.Description = e.ChildText(".entry-content>p")
		article.Link = fmt.Sprintf("https://ict-flex.nl/%s", slug.Make(article.Title))
	})

	if err := c.Visit("https://ict-flex.nl/category/mededelingen/"); err != nil {
		log.Fatal(err)

		return
	}

	buffer := new(bytes.Buffer)
	buffer.WriteString("<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" version=\"2.0\">\n")

	encoder := xml.NewEncoder(buffer)
	encoder.Indent("", "\t")

	err := encoder.Encode(ArticleFeed)
	if err != nil {
		log.Fatal(err)

		return
	}

	buffer.WriteString("</rss>")

	fmt.Println(buffer)
}
