package details

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenBDResponse struct {
	Onix    Onix    `json:"onix"`
	Hanmoto Hanmoto `json:"hanmoto"`
	Summary Summary `json:"summary"`
}

type Onix struct {
	DescriptiveDetail DescriptiveDetail `json:"DescriptiveDetail"`
}

type DescriptiveDetail struct {
	Subject []Subject `json:"Subject"`
}

type Subject struct {
	MainSubject             string `json:"MainSubject"`
	SubjectSchemeIdentifier string `json:"SubjectSchemeIdentifier"`
	SubjectCode             string `json:"SubjectCode"`
}

type Hanmoto struct {
	DateModified string `json:"datemodified"`
	DateCreated  string `json:"datecreated"`
	DateKoukai   string `json:"datekoukai"`
}

type Summary struct {
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	Series    string `json:"series"`
	Volume    string `json:"volume"`
	Publisher string `json:"publisher"`
	PubDate   string `json:"pubdate"`
	Author    string `json:"author"`
}

type openBDClient struct {
}

func (c *openBDClient) get(isbn string) (*OpenBDResponse, error) {
	openbdUrl := fmt.Sprintf("https://api.openbd.jp/v1/get?isbn=%s&pretty", isbn)
	resp, respErr := http.Get(openbdUrl)
	if respErr != nil {
		err := fmt.Errorf("OpenBD request failed: %s", respErr)
		return nil, err
	}

	var openBDResp []*OpenBDResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&openBDResp)
	if decodeErr != nil {
		err := fmt.Errorf("error in decoding OpenBD JSON response: %s", decodeErr)

		return nil, err
	}

	return openBDResp[0], nil
}
