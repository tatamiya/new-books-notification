package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/new-books-notification/src/details"
	"github.com/tatamiya/new-books-notification/src/models"
	"github.com/tatamiya/new-books-notification/src/notifier"
)

func TestGenerateUploadObjectOfFeed(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	date := time.Date(2022, time.July, 1, 12, 13, 24, 0, loc)
	inputFeed := gofeed.Feed{
		PublishedParsed: &date,
		Title:           "This is a Sample Feed!",
	}

	uploadObject, err := generateJsonUploadObject(&inputFeed)

	assert.Nil(t, err)
	assert.Equal(t, "feed20220701.json", uploadObject.ObjectName)
	assert.Equal(t, "application/json", uploadObject.ContentType)
}

type RecorderStub struct {
	RecordedISBN []string
	IsError      bool
}

func (r *RecorderStub) GetRecordedISBN(ctx context.Context, today time.Time) ([]string, error) {
	if r.IsError {
		return []string{}, fmt.Errorf("Could not get uploaded ISBNs!")
	}
	return r.RecordedISBN, nil
}

func (r *RecorderStub) SaveRecords(ctx context.Context, bookList *models.BookList) error {
	if r.IsError {
		return fmt.Errorf("Could not save records!")
	}
	for _, book := range bookList.Books {
		r.RecordedISBN = append(r.RecordedISBN, book.Isbn)
	}
	return nil
}

type NotifierStub struct {
	Messages []string
	IsError  bool
}

func (n *NotifierStub) Post(message string) error {
	if n.IsError {
		return fmt.Errorf("Could not send message!")
	}
	n.Messages = append(n.Messages, message)

	return nil
}

type DetailFetcherStub struct {
	details map[string]*details.DetailedInformation
	IsError bool
}

func (d *DetailFetcherStub) FetchDetailInfo(isbn string) (*details.DetailedInformation, error) {
	if d.IsError {
		return nil, fmt.Errorf("Could not get detailed information!")
	}

	detail := d.details[isbn]

	return detail, nil
}

func TestCoreProcessSkipsAlreadyUploadedBook(t *testing.T) {

	loc, _ := time.LoadLocation("Asia/Tokyo")
	dateUploaded := time.Date(2024, time.August, 1, 22, 42, 0, 0, loc)
	datePublished := time.Date(2024, time.September, 1, 22, 42, 0, 0, loc)
	inputBookList := models.BookList{
		UploadDate: dateUploaded,
		Books: []*models.Book{
			{
				Isbn:       "1111111111111",
				Title:      "Newly arrived book",
				Url:        "http://example.com/bd/isbn/1111111111111",
				PubDate:    datePublished,
				Categories: "",
			},
			{
				Isbn:       "2222222222222",
				Title:      "Newly arrived book",
				Url:        "http://example.com/bd/isbn/2222222222222",
				PubDate:    datePublished,
				Categories: "",
			},
			{
				Isbn:       "3333333333333",
				Title:      "Already uploaded book",
				Url:        "http://example.com/bd/isbn/3333333333333",
				PubDate:    datePublished,
				Categories: "",
			},
			{
				Isbn:       "4444444444444",
				Title:      "Newly arrived book",
				Url:        "http://example.com/bd/isbn/4444444444444",
				PubDate:    datePublished,
				Categories: "",
			},
		},
	}

	testRecorder := RecorderStub{
		RecordedISBN: []string{"3333333333333"},
		IsError:      false,
	}

	testDetailFetcher := DetailFetcherStub{
		details: map[string]*details.DetailedInformation{},
		IsError: false,
	}
	testNotifier := NotifierStub{
		IsError: false,
	}

	testFavoriteFilter := notifier.FavoriteFilter{}

	numUploaded := coreProcess(
		&inputBookList,
		&testDetailFetcher,
		&testRecorder,
		&testFavoriteFilter,
		&testNotifier,
	)

	assert.Equal(t, 3, numUploaded)
	assert.ElementsMatch(t, []string{"1111111111111", "2222222222222", "3333333333333", "4444444444444"}, testRecorder.RecordedISBN)

}

func TestCoreProcessNotifyingFavoriteBooks(t *testing.T) {

	loc, _ := time.LoadLocation("Asia/Tokyo")
	dateUploaded := time.Date(2024, time.August, 1, 22, 42, 0, 0, loc)
	datePublished := time.Date(2024, time.September, 1, 22, 42, 0, 0, loc)
	inputBookList := models.BookList{
		UploadDate: dateUploaded,
		Books: []*models.Book{
			{
				Isbn:       "1111111111111",
				Title:      "Book with favorite category",
				Url:        "http://example.com/bd/isbn/1111111111111",
				PubDate:    datePublished,
				Categories: "自然科学",
			},
			{
				Isbn:       "2222222222222", // Content: "物理学"
				Title:      "Book with unfavorite category and favorite content",
				Url:        "http://example.com/bd/isbn/2222222222222",
				PubDate:    datePublished,
				Categories: "趣味・実用",
			},
			{
				Isbn:       "3333333333333", // Content: "その他の工業"
				Title:      "Book with unfavorite category and content",
				Url:        "http://example.com/bd/isbn/3333333333333",
				PubDate:    datePublished,
				Categories: "趣味・実用",
			},
		},
	}

	testDetails := map[string]*details.DetailedInformation{
		"1111111111111": {Content: "物理学"},
		"2222222222222": {Content: "物理学"},
		"3333333333333": {Content: "その他の工業"},
	}
	testDetailFetcher := DetailFetcherStub{
		details: testDetails,
		IsError: false,
	}
	testRecorder := RecorderStub{
		IsError: false,
	}
	testNotifier := NotifierStub{
		IsError: false,
	}

	testFavoriteFilter := notifier.FavoriteFilter{
		FavoriteCategories: []string{"自然科学"},
		FavoriteContents:   []string{"物理学"},
	}

	_ = coreProcess(
		&inputBookList,
		&testDetailFetcher,
		&testRecorder,
		&testFavoriteFilter,
		&testNotifier,
	)

	assert.Equal(t, 2, len(testNotifier.Messages))

}

func TestCoreProcessSkipsNotifyingAlreadyUploadedFavoriteBooks(t *testing.T) {

	loc, _ := time.LoadLocation("Asia/Tokyo")
	dateUploaded := time.Date(2024, time.August, 1, 22, 42, 0, 0, loc)
	datePublished := time.Date(2024, time.September, 1, 22, 42, 0, 0, loc)
	inputBookList := models.BookList{
		UploadDate: dateUploaded,
		Books: []*models.Book{
			{
				Isbn:       "1111111111111",
				Title:      "Newly arrived favorite book",
				Url:        "http://example.com/bd/isbn/1111111111111",
				PubDate:    datePublished,
				Categories: "自然科学",
			},
			{
				Isbn:       "2222222222222",
				Title:      "Already uploaded favorite book",
				Url:        "http://example.com/bd/isbn/2222222222222",
				PubDate:    datePublished,
				Categories: "自然科学",
			},
		},
	}

	testRecorder := RecorderStub{
		RecordedISBN: []string{"2222222222222"},
		IsError:      false,
	}

	testDetailFetcher := DetailFetcherStub{
		details: map[string]*details.DetailedInformation{},
		IsError: false,
	}
	testNotifier := NotifierStub{
		IsError: false,
	}
	testFavoriteFilter := notifier.FavoriteFilter{
		FavoriteCategories: []string{"自然科学"},
	}

	_ = coreProcess(
		&inputBookList,
		&testDetailFetcher,
		&testRecorder,
		&testFavoriteFilter,
		&testNotifier,
	)

	assert.Equal(t, 1, len(testNotifier.Messages))

}

func TestCoreProcessMakesNotificationWhenFetchedDetailIsEmpty(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	dateUploaded := time.Date(2024, time.August, 1, 22, 42, 0, 0, loc)
	datePublished := time.Date(2024, time.September, 1, 22, 42, 0, 0, loc)
	inputBookList := models.BookList{
		UploadDate: dateUploaded,
		Books: []*models.Book{
			{
				Isbn:       "1111111111111",
				Title:      "Newly arrived book with favorite category but no details",
				Url:        "http://example.com/bd/isbn/1111111111111",
				PubDate:    datePublished,
				Categories: "自然科学",
			},
		},
	}

	testRecorder := RecorderStub{
		IsError: false,
	}

	testDetailFetcher := DetailFetcherStub{
		details: map[string]*details.DetailedInformation{},
		IsError: false,
	}
	testNotifier := NotifierStub{
		IsError: false,
	}
	testFavoriteFilter := notifier.FavoriteFilter{
		FavoriteCategories: []string{"自然科学"},
	}

	_ = coreProcess(
		&inputBookList,
		&testDetailFetcher,
		&testRecorder,
		&testFavoriteFilter,
		&testNotifier,
	)

	assert.Equal(t, 1, len(testNotifier.Messages))

}

func TestCoreProcessMakesNotificationWhenFetchingDetailsFails(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	dateUploaded := time.Date(2024, time.August, 1, 22, 42, 0, 0, loc)
	datePublished := time.Date(2024, time.September, 1, 22, 42, 0, 0, loc)
	inputBookList := models.BookList{
		UploadDate: dateUploaded,
		Books: []*models.Book{
			{
				Isbn:       "1111111111111",
				Title:      "Newly arrived book with favorite category",
				Url:        "http://example.com/bd/isbn/1111111111111",
				PubDate:    datePublished,
				Categories: "自然科学",
			},
		},
	}

	testRecorder := RecorderStub{
		IsError: false,
	}

	testDetailFetcher := DetailFetcherStub{
		IsError: true,
	}
	testNotifier := NotifierStub{
		IsError: false,
	}
	testFavoriteFilter := notifier.FavoriteFilter{
		FavoriteCategories: []string{"自然科学"},
	}

	_ = coreProcess(
		&inputBookList,
		&testDetailFetcher,
		&testRecorder,
		&testFavoriteFilter,
		&testNotifier,
	)

	assert.Equal(t, 1, len(testNotifier.Messages))

}
