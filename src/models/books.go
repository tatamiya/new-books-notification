package models

import (
	"net/url"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
)

type BookList struct {
	UploadDate time.Time
	Books      []*Book
}

type Book struct {
	Isbn            string
	Title           string
	Url             string
	Authors         string
	Publisher       string
	Categories      []string
	Ccode           string
	Target          string
	Format          string
	Genre           string
	PubDate         time.Time
	CreatedDate     time.Time
	LastUpdatedDate time.Time
}

func NewBookListFromFeed(feed *gofeed.Feed) *BookList {

	var books []*Book
	for _, item := range feed.Items {
		book := Book{
			Isbn:       extractISBN(item.Link),
			Title:      item.Title,
			Url:        item.Link,
			Categories: item.Categories,
			PubDate:    *item.PublishedParsed,
		}
		books = append(books, &book)
	}

	return &BookList{
		UploadDate: *feed.PublishedParsed,
		Books:      books,
	}
}

func extractISBN(link string) string {
	u, _ := url.Parse(link)
	re := regexp.MustCompile(`[0-9]{13}`)
	return re.FindString(u.Path)
}
