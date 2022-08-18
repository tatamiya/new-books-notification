package recorder

import (
	"context"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/new-books-notification/src/models"
)

func TestConvertBookIntoRecordCorrectly(t *testing.T) {

	loc, _ := time.LoadLocation("Asia/Tokyo")
	date := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	inputBook := models.Book{
		Isbn:            "1111111111111",
		Title:           "\nご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Authors:         "tatamiya tamiya",
		Publisher:       "畳屋書店",
		Url:             "http://example.com/bd/isbn/1111111111111",
		PubDate:         date,
		Categories:      []string{"自然科学"},
		Ccode:           "1042",
		Target:          "教養",
		Format:          "単行本",
		Content:         "物理学",
		LastUpdatedDate: date,
		CreatedDate:     date,
	}
	uploadedDate := time.Date(2022, time.August, 1, 12, 30, 0, 0, loc)
	actualRecord := convertIntoRecord(&inputBook, uploadedDate)
	expectedRecord := Record{
		ISBN:          "1111111111111",
		Title:         "\nご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Authors:       "tatamiya tamiya",
		Publisher:     "畳屋書店",
		Url:           "http://example.com/bd/isbn/1111111111111",
		PubDate:       civil.DateOf(date),
		Categories:    "自然科学",
		Ccode:         "1042",
		Target:        "教養",
		Format:        "単行本",
		Content:       "物理学",
		LastUpdatedAt: date,
		CreatedAt:     date,
		UploadedAt:    uploadedDate,
		UploadedDate:  civil.DateOf(uploadedDate),
	}

	assert.EqualValues(t, &expectedRecord, actualRecord)

}

func TestSaveRecordsCorrectly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping")
	}
	bqSettings := BQSettings{
		ProjectID:   os.Getenv("GCP_PROJECT_ID"),
		DatasetName: os.Getenv("GCP_BIGQUERY_DATASET"),
		TableName:   os.Getenv("GCP_BIGQUERY_TABLE"),
	}
	ctx := context.Background()
	recorder, err := NewBQRecorder(ctx, &bqSettings)

	assert.Nil(t, err)

	loc, _ := time.LoadLocation("Asia/Tokyo")
	date1 := time.Date(2024, time.August, 31, 12, 13, 24, 0, loc)
	sampleBook1 := models.Book{
		Isbn:            "1111111111111",
		Title:           "ご冗談でしょう、tatamiyaさん - tatamiya tamiya(著 / 文) | 畳屋書店",
		Authors:         "tatamiya tamiya",
		Publisher:       "畳屋書店",
		Url:             "http://example.com/bd/isbn/1111111111111",
		PubDate:         date1,
		Categories:      []string{"自然科学"},
		Ccode:           "1042",
		Target:          "教養",
		Format:          "単行本",
		Content:         "物理学",
		LastUpdatedDate: date1,
		CreatedDate:     date1,
	}
	date2 := time.Date(2124, time.February, 29, 0, 0, 0, 0, loc)
	sampleBook2 := models.Book{
		Isbn:            "9999999999999",
		Title:           "流体力学（後編） - 今井功(著 / 文) | 裳華房 ",
		Url:             "http://example.com/bd/isbn/9999999999999",
		Authors:         "今井功",
		Publisher:       "裳華房",
		PubDate:         date2,
		Categories:      []string{""},
		Ccode:           "",
		Target:          "",
		Format:          "",
		Content:         "",
		LastUpdatedDate: date2,
		CreatedDate:     date2,
	}
	uploadedDate := time.Date(2022, time.August, 1, 12, 30, 0, 0, loc)
	inputBookList := models.BookList{UploadDate: uploadedDate, Books: []*models.Book{
		&sampleBook1, &sampleBook2,
	}}

	err = recorder.SaveRecords(ctx, &inputBookList)
	assert.Nil(t, err)

}

func TestGetISBNCorrectly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping")
	}
	bqSettings := BQSettings{
		ProjectID:   os.Getenv("GCP_PROJECT_ID"),
		DatasetName: os.Getenv("GCP_BIGQUERY_DATASET"),
		TableName:   os.Getenv("GCP_BIGQUERY_TABLE"),
	}
	ctx := context.Background()
	recorder, err := NewBQRecorder(ctx, &bqSettings)

	assert.Nil(t, err)

	loc, _ := time.LoadLocation("Asia/Tokyo")
	date := time.Date(2022, time.August, 1, 0, 0, 0, 0, loc)

	actualUploadedISBN, err := recorder.GetRecordedISBN(ctx, date)
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"1111111111111", "9999999999999"}, actualUploadedISBN)

}

func TestGetEmptyWhenNoRecordIsUploaded(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping")
	}
	bqSettings := BQSettings{
		ProjectID:   os.Getenv("GCP_PROJECT_ID"),
		DatasetName: os.Getenv("GCP_BIGQUERY_DATASET"),
		TableName:   os.Getenv("GCP_BIGQUERY_TABLE"),
	}
	ctx := context.Background()
	recorder, err := NewBQRecorder(ctx, &bqSettings)

	assert.Nil(t, err)

	loc, _ := time.LoadLocation("Asia/Tokyo")
	date := time.Date(2122, time.August, 1, 0, 0, 0, 0, loc)

	actualUploadedISBN, err := recorder.GetRecordedISBN(ctx, date)
	assert.Nil(t, err)
	assert.EqualValues(t, []string{}, actualUploadedISBN)

}
