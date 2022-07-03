package models

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
