package details

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
