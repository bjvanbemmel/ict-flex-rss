package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	. "github.com/bjvanbemmel/ict-flex-rss/types"
	"github.com/gocolly/colly/v2"
)

const (
	POST_ID_EXPR string = "(post-[0-9]+)"
	DATE_FORMAT     = "2006-01-02T15:04:05+00:00"
)

var (
	post_regex *regexp.Regexp = regexp.MustCompile(POST_ID_EXPR)
)

func main() {
	c := colly.NewCollector()

    c.OnHTML("html", func(e *colly.HTMLElement) {
        nodes := e.DOM.Children().Find("body>main#content.type-post").Nodes

        if len(nodes) < 1 {
            return
        }

        var article Article

        e.ForEach("head>meta[property]", func(_ int, e *colly.HTMLElement) {
            content := e.Attr("content")

            switch e.Attr("property") {
            case "og:url":
                article.Link = content
            case "og:title":
                article.Title = content
            case "og:description":
                article.Description = content
            case "article:published_time":
                date, err := time.Parse(DATE_FORMAT, content)
                if err != nil {
                    log.Fatal(err)
                    return
                }
                article.CreatedAt = date
            }
        })

        var id string
        post := post_regex.FindString(e.ChildAttr("main#content", "class"))
        if post == "" {
            log.Println("Post ID could not be found")
        } else {
            id = strings.Split(post, "-")[1]
        }

        article.Guid.Id = id
        article.Author = e.ChildAttr("head>meta[name='twitter:data1']", "content")

        ArticleFeed.Articles = append(ArticleFeed.Articles, &article)
    })

	c.OnHTML("div.elementor-posts", func(e *colly.HTMLElement) {
		e.ForEach("article", func(_ int, e *colly.HTMLElement) {
			e.Request.Visit(e.ChildAttr(".elementor-post__card>.elementor-post__text>.elementor-post__title>a", "href"))
		})
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
