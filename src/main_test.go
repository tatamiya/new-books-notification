package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/new-books-notification/src/models"
	"github.com/tatamiya/new-books-notification/src/notifier"
	"github.com/tatamiya/new-books-notification/src/openbd"
	"github.com/tatamiya/new-books-notification/src/subject"
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
	details []*openbd.OpenBDResponse
	IsError bool
}

func (d *DetailFetcherStub) FetchDetailInfo(isbn string) (*openbd.OpenBDResponse, error) {
	if d.IsError {
		return nil, fmt.Errorf("Could not get detailed information!")
	}

	for _, detail := range d.details {
		if detail.Summary.ISBN == isbn {
			return detail, nil
		}
	}

	return &openbd.OpenBDResponse{}, nil
}

var testDecoder = subject.SubjectDecoder{
	Taishou: map[string]string{
		"0": "一般",
	},
	Keitai: map[string]string{
		"0": "単行本",
	},
	Naiyou: map[string]string{
		"42": "物理学",
		"58": "その他の工業",
	},
}

func generageMockOpenBDResponse(isbn string, content string) *openbd.OpenBDResponse {

	ccode := "0099"
	for k, v := range testDecoder.Naiyou {
		if v == content {
			ccode = k
		}
	}

	return &openbd.OpenBDResponse{
		Onix: openbd.Onix{
			DescriptiveDetail: openbd.DescriptiveDetail{
				Subject: []openbd.Subject{
					{SubjectCode: ccode}},
			},
		},
		Summary: openbd.Summary{
			ISBN: isbn,
		},
	}
}

func TestCoreProcess(t *testing.T) {

	loc, _ := time.LoadLocation("Asia/Tokyo")
	dateUploaded := time.Date(2024, time.August, 1, 22, 42, 0, 0, loc)
	datePublished := time.Date(2024, time.September, 1, 22, 42, 0, 0, loc)
	inputBookList := models.BookList{
		UploadDate: dateUploaded,
		Books: []*models.Book{
			{
				Isbn:       "1111111111111",
				Title:      "Normal test book",
				Url:        "http://example.com/bd/isbn/1111111111111",
				PubDate:    datePublished,
				Categories: []string{"自然科学"},
			},
			{
				Isbn:       "2222222222222",
				Title:      "Test book with unfavorite category and content",
				Url:        "http://example.com/bd/isbn/2222222222222",
				PubDate:    datePublished,
				Categories: []string{"学参"},
			},
			{
				Isbn:       "3333333333333",
				Title:      "Test book already uploaded",
				Url:        "http://example.com/bd/isbn/3333333333333",
				PubDate:    datePublished,
				Categories: []string{"自然科学"},
			},
			{
				Isbn:       "4444444444444",
				Title:      "Test book not registered in OpenBD",
				Url:        "http://example.com/bd/isbn/4444444444444",
				PubDate:    datePublished,
				Categories: []string{""},
			},
		},
	}

	testOpenBDResponses := []*openbd.OpenBDResponse{
		generageMockOpenBDResponse("1111111111111", "物理学"),
		generageMockOpenBDResponse("2222222222222", "その他の工業"),
		generageMockOpenBDResponse("3333333333333", "物理学"),
	}
	testDetailFetcher := DetailFetcherStub{
		details: testOpenBDResponses,
		IsError: false,
	}
	testRecorder := RecorderStub{
		RecordedISBN: []string{"3333333333333"},
		IsError:      false,
	}
	testNotifier := NotifierStub{
		IsError: false,
	}

	testFavoriteFilter := notifier.FavoriteFilter{
		FavoriteCategories: []string{"自然科学"},
		FavoriteContents:   []string{"物理学"},
	}

	numUploaded := coreProcess(
		&inputBookList,
		&testDetailFetcher,
		&testDecoder,
		&testRecorder,
		&testFavoriteFilter,
		&testNotifier,
	)

	assert.Equal(t, 3, numUploaded)
	assert.Equal(t, 1, len(testNotifier.Messages))
	assert.ElementsMatch(t, []string{"1111111111111", "2222222222222", "3333333333333", "4444444444444"}, testRecorder.RecordedISBN)

}
