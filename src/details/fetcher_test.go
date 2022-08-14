package details

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var sampleDecoder = SubjectDecoder{
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

var sampleOpenBDResp = OpenBDResponse{
	Onix: Onix{
		DescriptiveDetail: DescriptiveDetail{
			Subject: []Subject{
				{
					MainSubject:             "",
					SubjectSchemeIdentifier: "00",
					SubjectCode:             "1040",
				},
			},
		},
	},
	Hanmoto: Hanmoto{
		DateModified: "2022-08-01 18:18:39",
		DateCreated:  "2022-06-30 18:22:39",
		DateKoukai:   "20220701",
	},
	Summary: Summary{
		ISBN:      "1111111111111",
		Title:     "ご冗談でしょう、tatamiyaさん",
		Series:    "シリーズ畳の不思議",
		Volume:    "1",
		Publisher: "畳屋書店",
		PubDate:   "20240831",
		Author:    "tatamiya tamiya／著 畳の科学／編集",
	},
}

type openBDClientStub struct {
	Responses map[string]*OpenBDResponse
	IsError   bool
}

func (c *openBDClientStub) get(isbn string) (*OpenBDResponse, error) {
	if c.IsError {
		return nil, fmt.Errorf("OpenBD request failed!")
	}

	res, ok := c.Responses[isbn]

	if !ok {
		return nil, nil
	}

	return res, nil
}

func TestFetcherParsesOpenBDResopnseCorrectly(t *testing.T) {

	testFetcher := OpenBDDetailsFetcher{
		client: &openBDClientStub{
			Responses: map[string]*OpenBDResponse{
				"1111111111111": &sampleOpenBDResp,
			},
			IsError: false,
		},
		decoder: &sampleDecoder,
	}

	actualDetailedInfo, err := testFetcher.FetchDetailInfo("1111111111111")

	loc, _ := time.LoadLocation("Asia/Tokyo")
	createdDate, _ := time.ParseInLocation("2006-01-02 15:04:05", "2022-06-30 18:22:39", loc)
	lastUpdatedDate, _ := time.ParseInLocation("2006-01-02 15:04:05", "2022-08-01 18:18:39", loc)
	expectedDetailedInfo := DetailedInformation{
		Author:          "tatamiya tamiya／著 畳の科学／編集",
		Publisher:       "畳屋書店",
		CreatedDate:     createdDate,
		LastUpdatedDate: lastUpdatedDate,
		Ccode:           "1040",
		Format:          "単行本",
		Target:          "教養",
		Content:         "自然科学総記",
	}

	assert.Nil(t, err)
	assert.EqualValues(t, expectedDetailedInfo, *actualDetailedInfo)

}
