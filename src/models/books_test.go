package models

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/new-books-notification/src/openbd"
	"github.com/tatamiya/new-books-notification/src/subject"
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

func TestExtractISBN(t *testing.T) {
	inputURL := "http://example.com/bd/isbn/9999999999999"
	expectedISBN := "9999999999999"

	actualISBN := extractISBN(inputURL)
	assert.Equal(t, expectedISBN, actualISBN)
}

var inputOpenBDResp = openbd.OpenBDResponse{
	Onix: openbd.Onix{
		DescriptiveDetail: openbd.DescriptiveDetail{
			Subject: []openbd.Subject{
				{
					MainSubject:             "",
					SubjectSchemeIdentifier: "00",
					SubjectCode:             "1040",
				},
			},
		},
	},
	Hanmoto: openbd.Hanmoto{
		DateModified: "2022-08-01 18:18:39",
		DateCreated:  "2022-06-30 18:22:39",
		DateKoukai:   "20220701",
	},
	Summary: openbd.Summary{
		ISBN:      "1111111111111",
		Title:     "ご冗談でしょう、tatamiyaさん",
		Series:    "シリーズ畳の不思議",
		Volume:    "1",
		Publisher: "畳屋書店",
		PubDate:   "20240831",
		Author:    "tatamiya tamiya／著 畳の科学／編集",
	},
}

func TestUpdateBookInfo(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	date1 := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	sampleBook := Book{
		Isbn:       "1111111111111",
		Title:      "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Url:        "http://example.com/bd/isbn/1111111111111",
		PubDate:    date1,
		Categories: []string{"自然科学"},
	}

	sampleBook.UpdateInfoFrom(&inputOpenBDResp)
	assert.Equal(t, "畳屋書店", sampleBook.Publisher)
	assert.Equal(t, "tatamiya tamiya／著 畳の科学／編集", sampleBook.Authors)
	assert.Equal(t, "1040", sampleBook.Ccode)

	expectedLastUpdatedDate := time.Date(2022, time.August, 1, 18, 18, 39, 0, loc)
	assert.Equal(t, expectedLastUpdatedDate, sampleBook.LastUpdatedDate)

	expectedCreatedDate := time.Date(2022, time.June, 30, 18, 22, 39, 0, loc)
	assert.Equal(t, expectedCreatedDate, sampleBook.CreatedDate)

	assert.Equal(t, "", sampleBook.Target)
	assert.Equal(t, "", sampleBook.Format)
	assert.Equal(t, "", sampleBook.Content)
}

func TestUpdateWithEmptyCcodeWhenSubjectIsEmpty(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	date1 := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	sampleBook := Book{
		Isbn:       "1111111111111",
		Title:      "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Url:        "http://example.com/bd/isbn/1111111111111",
		PubDate:    date1,
		Categories: []string{"自然科学"},
	}

	inputOpenBDResp.Onix = openbd.Onix{
		DescriptiveDetail: openbd.DescriptiveDetail{
			Subject: []openbd.Subject{},
		}}

	sampleBook.UpdateInfoFrom(&inputOpenBDResp)
	assert.Equal(t, "畳屋書店", sampleBook.Publisher)
	assert.Equal(t, "tatamiya tamiya／著 畳の科学／編集", sampleBook.Authors)
	assert.Equal(t, "", sampleBook.Ccode)

	expectedLastUpdatedDate := time.Date(2022, time.August, 1, 18, 18, 39, 0, loc)
	assert.Equal(t, expectedLastUpdatedDate, sampleBook.LastUpdatedDate)

	expectedCreatedDate := time.Date(2022, time.June, 30, 18, 22, 39, 0, loc)
	assert.Equal(t, expectedCreatedDate, sampleBook.CreatedDate)

	assert.Equal(t, "", sampleBook.Target)
	assert.Equal(t, "", sampleBook.Format)
	assert.Equal(t, "", sampleBook.Content)
}

var sampleDecoder = subject.SubjectDecoder{
	Taishou: map[string]string{
		"0": "一般",
		"1": "教養",
		"3": "専門書",
	},
	Keitai: map[string]string{
		"0": "単行本",
		"1": "文庫",
	},
	Naiyou: map[string]string{
		"40": "自然科学総記",
		"42": "物理学",
	},
}

func TestUpdateSubject(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	date1 := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	sampleBook := Book{
		Isbn:       "1111111111111",
		Title:      "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Url:        "http://example.com/bd/isbn/1111111111111",
		PubDate:    date1,
		Categories: []string{"自然科学"},
		Ccode:      "1042",
	}
	err := sampleBook.UpdateSubject(&sampleDecoder)

	assert.Nil(t, err)

	assert.EqualValues(t, "教養", sampleBook.Target)
	assert.EqualValues(t, "単行本", sampleBook.Format)
	assert.EqualValues(t, "物理学", sampleBook.Content)

}

func TestNotUpdateSubjectWhenCcodeIsInvalid(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	date1 := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	sampleBook := Book{
		Isbn:       "1111111111111",
		Title:      "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Url:        "http://example.com/bd/isbn/1111111111111",
		PubDate:    date1,
		Categories: []string{"自然科学"},
		Ccode:      "",
	}

	err := sampleBook.UpdateSubject(&sampleDecoder)

	assert.NotNil(t, err)

	assert.EqualValues(t, "", sampleBook.Target)
	assert.EqualValues(t, "", sampleBook.Format)
	assert.EqualValues(t, "", sampleBook.Content)

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
