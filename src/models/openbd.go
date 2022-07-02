package models

import (
	"time"
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
	Contributor string `json:"contributor"`
}

type Hanmoto struct {
	DateModified *time.Time `json:"datemodified"`
	DateCreated  *time.Time `json:"datecreated"`
	DateKoukai   *time.Time `json:"datekoukai"`
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
