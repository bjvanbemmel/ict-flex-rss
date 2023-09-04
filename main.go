package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/bjvanbemmel/ict-flex-rss/types"
	"github.com/gocolly/colly/v2"
	"github.com/gosimple/slug"
)

const (
	POST_ID_EXPR string = "(post-[0-9]+)"
)

var (
	post_regex *regexp.Regexp = regexp.MustCompile(POST_ID_EXPR)
)

func main() {
	c := colly.NewCollector()

	c.OnHTML("div.elementor-posts", func(e *colly.HTMLElement) {
		e.ForEach("article", func(_ int, e *colly.HTMLElement) {
			var article Article = Article{
				Guid: ArticleGuid{
					IsPermaLink: false,
				},
			}

			classes := e.Attr("class")
			post := post_regex.FindString(classes)
			if post == "" {
				return
			}
			article.Guid.Id = strings.Split(post, "-")[1]

			var rawDate string = e.ChildText(".elementor-post__card>.elementor-post__meta-data>span") //.entry-meta>.meta-date>time", "datetime

			var err error
			article.CreatedAt, err = getDateFromNaturalLanguage(rawDate)
			if err != nil {
				return
			}

			ArticleFeed.Articles = append(ArticleFeed.Articles, &article)

			e.Request.Visit(e.ChildAttr(".elementor-post__card>.elementor-post__text>.elementor-post__title>a", "href"))
		})
	})

	c.OnHTML("main.post", func(e *colly.HTMLElement) {
		// if hero := e.ChildText("div.hero-section"); hero == "" {
		// 	return
		// }

		classes := e.Attr("class")
		post := post_regex.FindString(classes)
		if post == "" {
			return
		}
		id := strings.Split(post, "-")[1]

		var article *Article

		for _, art := range ArticleFeed.Articles {
			if art.Guid.Id != id {
				continue
			}

			article = art
		}

		article.Title = e.ChildText(".page-header>.entry-title")
		article.Description = e.ChildText(".page-content>p")
		article.Link = fmt.Sprintf("https://ict-flex.nl/%s", slug.Make(article.Title))
	})

	if err := c.Visit("https://ict-flex.nl/mededelingen/"); err != nil {
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

// This entire function is incredibly cursed but Go simply lacks the language features
// to do this somewhat relatively cleanly.
func getDateFromNaturalLanguage(date string) (time.Time, error) {
	dateFields := strings.Split(date, " ")

	day, err := strconv.Atoi(dateFields[0])
	if err != nil {
		return time.Time{}, err
	}

	naturalMonth := dateFields[1]

	year, err := strconv.Atoi(dateFields[2])
	if err != nil {
		return time.Time{}, err
	}

	var month int
	switch naturalMonth {
	case "Januari":
		month = 1
	case "Februari":
		month = 2
	case "Maart":
		month = 3
	case "April":
		month = 4
	case "Mei":
		month = 5
	case "Juni":
		month = 6
	case "Juli":
		month = 7
	case "Augustus":
		month = 8
	case "September":
		month = 9
	case "Oktober":
		month = 10
	case "November":
		month = 11
	case "December":
		month = 12
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), err
}
