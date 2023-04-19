package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
)

type Article struct {
	Title     string
	Content   string
	Author    Author
	CreatedAt time.Time
}

type Author struct {
	Name    string
	Profile string
}

var Articles map[string]*Article = make(map[string]*Article)

func main() {
	c := colly.NewCollector()

	c.OnHTML("div.entries", func(e *colly.HTMLElement) {
		e.ForEach("article", func(_ int, e *colly.HTMLElement) {
			var article Article = Article{}

			article.Author.Name = e.ChildText(".entry-meta>.meta-author>a>span")
			article.Author.Profile = e.ChildAttr(".entry-meta>.meta-author>a", "href")

			var date string = e.ChildAttr(".entry-meta>.meta-date>time", "datetime")

			article.CreatedAt, _ = time.Parse(time.RFC3339, date)

			Articles[e.Attr("id")] = &article

			e.Request.Visit(e.ChildAttr(".entry-title>a", "href"))
		})
	})

	c.OnHTML("article.type-post", func(e *colly.HTMLElement) {
		if hero := e.ChildText("div.hero-section"); hero == "" {
			return
		}

		article, found := Articles[e.Attr("id")]
		if !found {
			return
		}

		article.Title = e.ChildText(".hero-section>.entry-header>.page-title")
		article.Content = e.ChildText(".entry-content>p")

	})

	if err := c.Visit("https://ict-flex.nl/category/mededelingen/"); err != nil {
		panic(err)
	}

	raw, err := json.Marshal(Articles)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(raw))
}
