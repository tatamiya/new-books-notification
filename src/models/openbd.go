package models

import (
	"github.com/tatamiya/gcp-cost-notification/src/datetime"
	"google.golang.org/genproto/googleapis/type/datetime"
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
	DateModified datetime.DateTime `json:"datemodified"`
	DateCreated  datetime.DateTime `json:"datecreated"`
	DateKoukai   datetime.DateTime `json:"datekoukai"`
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
