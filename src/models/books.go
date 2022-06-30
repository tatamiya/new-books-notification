package models

import (
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
