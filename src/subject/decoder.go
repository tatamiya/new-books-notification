package subject

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type SubjectDecoder struct {
	Taishou map[string]string `json:"taishou"`
	Keitai  map[string]string `json:"keitai"`
	Naiyou  map[string]string `json:"naiyou"`
}

type DecodedSubject struct {
	Ccode  string
	Target string
	Format string
	Genre  string
}

func (s *SubjectDecoder) decode(ccode string) (*DecodedSubject, error) {

	if _, err := strconv.Atoi(ccode); err != nil {
		return nil, fmt.Errorf("Invalid Ccode! %s cannot be converted to digits: %s", ccode, err)
	}
	if len(ccode) != 4 {
		return nil, fmt.Errorf("Invalid Ccode! %s is not 4 digits", ccode)
	}
	chars := []rune(ccode)
	target := s.Taishou[string(chars[0])]
	format := s.Keitai[string(chars[1])]
	genre := s.Naiyou[string(chars[2:])]

	return &DecodedSubject{
		ccode, target, format, genre,
	}, nil
}

func NewSubjectDecoder() (*SubjectDecoder, error) {

	var decoder SubjectDecoder
	ccodeData, ioErr := ioutil.ReadFile("./ccode.json")
	if ioErr != nil {
		return nil, fmt.Errorf("Could not read ccode.json!: %s", ioErr)
	}
	jsonErr := json.Unmarshal(ccodeData, &decoder)
	if jsonErr != nil {
		return nil, fmt.Errorf("Could not unmarshal json data!: %s", jsonErr)
	}

	return &decoder, nil
}
