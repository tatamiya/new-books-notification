package openbd

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenBDDetailFetcher struct {
}

func NewOpenBDDetailFetcher() *OpenBDDetailFetcher {
	return &OpenBDDetailFetcher{}
}

func (f *OpenBDDetailFetcher) GetDetailedBookInfo(isbn string) (*OpenBDResponse, error) {
	openbdUrl := fmt.Sprintf("https://api.openbd.jp/v1/get?isbn=%s&pretty", isbn)
	resp, respErr := http.Get(openbdUrl)
	if respErr != nil {
		err := fmt.Errorf("OpenBD request failed: %s", respErr)
		return nil, err
	}

	var openBDResp []*OpenBDResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&openBDResp)
	if decodeErr != nil {
		err := fmt.Errorf("Error in decoding OpenBD JSON response: %s", decodeErr)

		return nil, err
	}

	return openBDResp[0], nil
}
