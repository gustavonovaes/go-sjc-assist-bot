package news

import "encoding/xml"

type News struct {
	Title   string
	Link    string
	Origin  string
	Content string
	Tags    []string
}

type MeonRssChannel struct {
	XMLName xml.Name     `xml:"rss"`
	Channel MeonRssItems `xml:"channel"`
}

type MeonRssItems struct {
	Items []MeonRssItem `xml:"item"`
}

type MeonRssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}
