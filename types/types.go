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
	Author      string      `xml:"author"`
	CreatedAt   time.Time   `xml:"pubDate"`
}

type ArticleGuid struct {
	Id          string `xml:",chardata"`
	IsPermaLink bool   `xml:"isPermaLink,attr"`
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
	Link:        "https://ict-flex.nl/mededelingen/",
	Description: "RSS Feed for ICT-Flex announcements",
}
