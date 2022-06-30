package models

import (
	"github.com/mmcdole/gofeed"
	"google.golang.org/genproto/googleapis/type/datetime"
)

type BookList struct {
	UploadDate datetime.DateTime
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
	PubDate         datetime.DateTime
	CreatedDate     datetime.DateTime
	LastUpdatedDate datetime.DateTime
}

func NewBookListFromFeed(feed *gofeed.Feed) *BookList {

	var books []*Book
	for _, item := range feed.Items {
		book := Book{
			Title:      item.Title,
			Url:        item.Link,
			Categories: item.Categories,
			PubDate:    item.PublishedParsed,
		}
		books = append(books, &book)
	}

	return &BookList{
		UploadDate: feed.PublishedParsed,
		Books:      books,
	}
}
