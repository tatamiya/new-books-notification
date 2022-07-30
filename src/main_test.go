package main

import (
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
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
