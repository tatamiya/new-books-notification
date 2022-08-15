package models

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/tatamiya/new-books-notification/src/details"
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
	Content         string
	PubDate         time.Time
	CreatedDate     time.Time
	LastUpdatedDate time.Time
}

func NewBookListFromFeed(feed *gofeed.Feed) *BookList {

	var books []*Book
	for _, item := range feed.Items {
		var trimmedCategories []string
		for _, category := range item.Categories {
			trimmedCategories = append(trimmedCategories, strings.TrimSpace(category))
		}
		book := Book{
			Isbn:       extractISBN(item.Link),
			Title:      strings.TrimSpace(item.Title),
			Url:        item.Link,
			Categories: trimmedCategories,
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

func (bl *BookList) FilterOut(isbns []string) *BookList {

	m := make(map[string]bool)
	for _, isbn := range isbns {
		m[isbn] = true
	}

	filteredBooks := []*Book{}
	for _, book := range bl.Books {
		if _, ok := m[book.Isbn]; !ok {
			filteredBooks = append(filteredBooks, book)
		}
	}
	return &BookList{
		UploadDate: bl.UploadDate,
		Books:      filteredBooks,
	}
}

func (b *Book) UpdateDetails(detailedInfo *details.DetailedInformation) {

	b.Authors = detailedInfo.Author
	b.Publisher = detailedInfo.Publisher

	b.Ccode = detailedInfo.Ccode
	b.Target = detailedInfo.Target
	b.Format = detailedInfo.Format
	b.Content = detailedInfo.Content

	b.CreatedDate = detailedInfo.CreatedDate
	b.LastUpdatedDate = detailedInfo.LastUpdatedDate
}

func (b *Book) AsNotificationMessage() string {
	url := b.Url
	title := strings.TrimSpace(b.Title)
	pubDate := b.PubDate.Format("2006/01/02")
	categories := strings.Join(b.Categories, ",")
	content := b.Content

	message := fmt.Sprintf("<%s|%s>\n発売日: %s\nカテゴリー: %s\n内容: %s", url, title, pubDate, categories, content)
	return message
}
