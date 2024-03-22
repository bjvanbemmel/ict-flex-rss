package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bjvanbemmel/ict-flex-rss/types"
	"github.com/gocolly/colly/v2"
)

const (
	POST_ID_EXPR string = "(post-[0-9]+)"
	DATE_FORMAT         = "2 January 2006"
)

var (
	post_regex *regexp.Regexp = regexp.MustCompile(POST_ID_EXPR)
)

func main() {
	c := colly.NewCollector()
	c.UserAgent = "This is the ICT-Flex-RSS scraper. I come in peace. The generated RSS feed can be found at https://rss.bjvanbemmel.nl/ict-flex."

	c.OnHTML("html", func(e *colly.HTMLElement) {
		nodes := e.DOM.Children().Find("div.elementor-posts-container > article")

		if len(nodes.Nodes) < 1 {
			return
		}

		nodes.Children().Each(func(index int, x *goquery.Selection) {
			var article types.Article

			title := x.Find(".elementor-post__title").
				First().
				Text()

			description := x.Find(".elementor-post__excerpt").
				First().
				Text()

			link, _ := x.Find(".elementor-post__title>a[href]").
				First().
				Attr("href")

			createdAt := x.Find(".elementor-post__meta-data>.elementor-post-date").
				First().
				Text()

			title = strings.ReplaceAll(title, "\n", "")
			description = strings.Trim(description, "\n")
			createdAt = strings.ReplaceAll(createdAt, "\n", "")
			createdAt = strings.ReplaceAll(createdAt, "\x09", "")

			identifier := fmt.Sprintf("%s%s%s%s", title, description, link, createdAt)
			hashedId := sha256.Sum256([]byte(identifier))

			splitCreatedAt := strings.Split(createdAt, " ")
			if len(splitCreatedAt) < 3 {
				return
			}

			splitCreatedAt[1] = translateMonth(splitCreatedAt[1])
			createdAt = strings.Join(splitCreatedAt, " ")

			date, err := time.Parse(DATE_FORMAT, createdAt)
			if err != nil {
				log.Println(err)
			}

			article.Title = title
			article.Description = description
			article.Link = link
			article.CreatedAt = date
			article.Guid = types.ArticleGuid{
				Id:          fmt.Sprintf("%x", hashedId),
				IsPermaLink: false,
			}

			types.ArticleFeed.Articles = append(types.ArticleFeed.Articles, &article)
		})
	})

	// c.OnHTML("html", func(e *colly.HTMLElement) {
	// 	nodes := e.DOM.Children().Find("body>main#content.type-post").Nodes

	// 	if len(nodes) < 1 {
	// 		return
	// 	}

	// 	var article Article

	// 	e.ForEach("head>meta[property]", func(_ int, e *colly.HTMLElement) {
	// 		content := e.Attr("content")

	// 		switch e.Attr("property") {
	// 		case "og:url":
	// 			article.Link = content
	// 		case "og:title":
	// 			article.Title = content
	// 		case "og:description":
	// 			article.Description = content
	// 		case "article:published_time":
	// 			date, err := time.Parse(DATE_FORMAT, content)
	// 			if err != nil {
	// 				log.Fatal(err)
	// 				return
	// 			}
	// 			article.CreatedAt = date
	// 		}
	// 	})

	// 	var id string
	// 	post := post_regex.FindString(e.ChildAttr("main#content", "class"))
	// 	if post == "" {
	// 		log.Println("Post ID could not be found")
	// 	} else {
	// 		id = strings.Split(post, "-")[1]
	// 	}

	// 	article.Guid.Id = id
	// 	article.Author = e.ChildAttr("head>meta[name='twitter:data1']", "content")

	// 	ArticleFeed.Articles = append(ArticleFeed.Articles, &article)
	// })

	// c.OnHTML("div.elementor-posts", func(e *colly.HTMLElement) {
	// 	e.ForEach("article", func(_ int, e *colly.HTMLElement) {
	// 		e.Request.Visit(e.ChildAttr(".elementor-post__card>.elementor-post__text>.elementor-post__title>a", "href"))
	// 	})
	// })

	if err := c.Visit("https://ict-flex.nl/mededelingen/"); err != nil {
		log.Fatal(err)
		return
	}

	buffer := new(bytes.Buffer)
	buffer.WriteString("<rss xmlns:atom=\"http://www.w3.org/2005/Atom\" version=\"2.0\">\n")

	encoder := xml.NewEncoder(buffer)
	encoder.Indent("", "\t")

	err := encoder.Encode(types.ArticleFeed)
	if err != nil {
		log.Fatal(err)

		return
	}

	buffer.WriteString("</rss>")

	fmt.Println(buffer)
}

func translateMonth(month string) string {
	switch month {
	case "januari":
		return "January"
	case "februari":
		return "February"
	case "maart":
		return "March"
	case "april":
		return "April"
	case "mei":
		return "May"
	case "juni":
		return "June"
	case "juli":
		return "July"
	case "augustus":
		return "August"
	case "september":
		return "September"
	case "oktober":
		return "October"
	case "november":
		return "November"
	case "december":
		return "December"
	default:
		return ""
	}
}
