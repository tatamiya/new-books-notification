package models

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestGenerateNewBookListFromFeedCorrectly(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	date1 := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	item1 := gofeed.Item{
		Title:           "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Link:            "http://example.com/bd/isbn/1111111111111",
		PublishedParsed: &date1,
		Categories:      []string{"自然科学"},
	}

	date2 := time.Date(2124, time.February, 29, 0, 0, 0, 0, loc)
	item2 := gofeed.Item{
		Title:           "流体力学（後編） - 今井功(著 / 文) | 裳華房",
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
				Title:      "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
				Url:        "http://example.com/bd/isbn/1111111111111",
				PubDate:    date1,
				Categories: []string{"自然科学"},
			},
			{
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
