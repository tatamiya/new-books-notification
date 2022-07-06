package models

import (
	"log"
	"net/url"
	"regexp"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/tatamiya/new-books-notification/src/openbd"
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

func (b *Book) UpdateInfoFrom(openbd *openbd.OpenBDResponse) {
	summary := openbd.Summary
	b.Authors = summary.Author
	b.Publisher = summary.Publisher

	subjects := openbd.Onix.DescriptiveDetail.Subject
	if len(subjects) > 0 {
		b.Ccode = subjects[0].SubjectCode
	}

	hanmoto := openbd.Hanmoto
	loc, _ := time.LoadLocation("Asia/Tokyo")

	dateCreated, err := time.ParseInLocation("2006-01-02 15:04:05", hanmoto.DateCreated, loc)
	if err != nil {
		log.Printf("Error in parsing timestamp: %s", hanmoto.DateCreated)
	}
	b.CreatedDate = dateCreated

	dateModified, err := time.ParseInLocation("2006-01-02 15:04:05", hanmoto.DateModified, loc)
	if err != nil {
		log.Printf("Error in parsing timestamp: %s", hanmoto.DateModified)
	}
	b.LastUpdatedDate = dateModified

}
