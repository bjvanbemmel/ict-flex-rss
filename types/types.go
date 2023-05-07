package types

import (
	"encoding/xml"
	"time"
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
