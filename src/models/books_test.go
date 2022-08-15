package models

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/new-books-notification/src/details"
)

func TestGenerateNewBookListFromFeedCorrectly(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	date1 := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	item1 := gofeed.Item{
		Title:           "\tご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店 ",
		Link:            "http://example.com/bd/isbn/1111111111111",
		PublishedParsed: &date1,
		Categories:      []string{" 自然科学 "},
	}

	date2 := time.Date(2124, time.February, 29, 0, 0, 0, 0, loc)
	item2 := gofeed.Item{
		Title:           "\t流体力学（後編） - 今井功(著 / 文) | 裳華房 ",
		Link:            "http://example.com/bd/isbn/9999999999999",
		PublishedParsed: &date2,
		Categories:      []string{""},
	}

	datePublished := time.Date(2024, time.September, 1, 22, 42, 0, 0, loc)
	inputFeed := gofeed.Feed{
		PublishedParsed: &datePublished,
		Items: []*gofeed.Item{
			&item1,
			&item2,
		},
	}

	expectedBookList := BookList{
		UploadDate: datePublished,
		Books: []*Book{
			{
				Isbn:       "1111111111111",
				Title:      "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
				Url:        "http://example.com/bd/isbn/1111111111111",
				PubDate:    date1,
				Categories: []string{"自然科学"},
			},
			{
				Isbn:       "9999999999999",
				Title:      "流体力学（後編） - 今井功(著 / 文) | 裳華房",
				Url:        "http://example.com/bd/isbn/9999999999999",
				PubDate:    date2,
				Categories: []string{""},
			},
		},
	}

	actualBookList := NewBookListFromFeed(&inputFeed)

	assert.EqualValues(t, expectedBookList, *actualBookList)

}

func TestFilterOutBookListByISBNCorrectly(t *testing.T) {
	sampleBookList := BookList{
		Books: []*Book{
			{
				Isbn:  "1111111111111",
				Title: "Book1",
				Url:   "http://example.com/bd/isbn/1111111111111",
			},
			{
				Isbn:  "2222222222222",
				Title: "Book2",
				Url:   "http://example.com/bd/isbn/2222222222222",
			},
			{
				Isbn:  "3333333333333",
				Title: "Book3",
				Url:   "http://example.com/bd/isbn/3333333333333",
			},
		},
	}

	inputISBNsToFilter := []string{"1111111111111", "3333333333333", "4444444444444"}
	expectedFileteredBookList := BookList{
		Books: []*Book{
			{
				Isbn:  "2222222222222",
				Title: "Book2",
				Url:   "http://example.com/bd/isbn/2222222222222",
			},
		},
	}

	actualFilteredBookList := sampleBookList.FilterOut(inputISBNsToFilter)

	assert.EqualValues(t, expectedFileteredBookList, *actualFilteredBookList)

}

func TestExtractISBN(t *testing.T) {
	inputURL := "http://example.com/bd/isbn/9999999999999"
	expectedISBN := "9999999999999"

	actualISBN := extractISBN(inputURL)
	assert.Equal(t, expectedISBN, actualISBN)
}

func TestUpdateDetails(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	pubDate := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	sampleBook := Book{
		Title:      "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Isbn:       "1111111111111",
		Url:        "http://example.com/bd/isbn/1111111111111",
		PubDate:    pubDate,
		Categories: []string{"自然科学"},
	}

	createdDate, _ := time.ParseInLocation("2006-01-02 15:04:05", "2022-06-30 18:22:39", loc)
	lastUpdatedDate, _ := time.ParseInLocation("2006-01-02 15:04:05", "2022-08-01 18:18:39", loc)
	inputDetailedInfo := details.DetailedInformation{
		Author:          "tatamiya tamiya／著 畳の科学／編集",
		Publisher:       "畳屋書店",
		CreatedDate:     createdDate,
		LastUpdatedDate: lastUpdatedDate,
		Ccode:           "1040",
		Format:          "単行本",
		Target:          "教養",
		Content:         "自然科学総記",
	}

	expectedUpdatedBook := Book{
		Isbn:            "1111111111111",
		Title:           "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Url:             "http://example.com/bd/isbn/1111111111111",
		Authors:         "tatamiya tamiya／著 畳の科学／編集",
		Publisher:       "畳屋書店",
		Categories:      []string{"自然科学"},
		Ccode:           "1040",
		Target:          "教養",
		Format:          "単行本",
		Content:         "自然科学総記",
		PubDate:         pubDate,
		CreatedDate:     createdDate,
		LastUpdatedDate: lastUpdatedDate,
	}

	sampleBook.UpdateDetails(&inputDetailedInfo)

	assert.EqualValues(t, expectedUpdatedBook, sampleBook)

}

func TestCreateNotificationMessageCorrectly(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	date1 := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	sampleBook := Book{
		Isbn:       "1111111111111",
		Title:      "\nご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Url:        "http://example.com/bd/isbn/1111111111111",
		PubDate:    date1,
		Categories: []string{"自然科学"},
		Ccode:      "1042",
		Target:     "教養",
		Format:     "単行本",
		Content:    "物理学",
	}

	expectedMessage := `<http://example.com/bd/isbn/1111111111111|ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店>
発売日: 2024/08/31
カテゴリー: 自然科学
内容: 物理学`

	actualMessage := sampleBook.AsNotificationMessage()

	assert.EqualValues(t, expectedMessage, actualMessage)
}
