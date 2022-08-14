package details

import (
	"log"
	"time"
)

type DetailedInformation struct {
	Author          string
	Publisher       string
	CreatedDate     time.Time
	LastUpdatedDate time.Time
	Ccode           string
	Format          string
	Target          string
	Content         string
}

type openBDClientInterface interface {
	get(string) (*OpenBDResponse, error)
}

type OpenBDDetailsFetcher struct {
	client  openBDClientInterface
	decoder *SubjectDecoder
}

func (f *OpenBDDetailsFetcher) FetchDetailInfo(isbn string) (*DetailedInformation, error) {

	res, err := f.client.get(isbn)
	if err != nil {
		return nil, err
	} else if res == nil {
		return nil, nil
	}
	summary := res.Summary
	author := summary.Author
	publisher := summary.Publisher

	hanmoto := res.Hanmoto
	loc, _ := time.LoadLocation("Asia/Tokyo")

	createdDate, err := time.ParseInLocation("2006-01-02 15:04:05", hanmoto.DateCreated, loc)
	if err != nil {
		log.Printf("Error in parsing timestamp: %s", hanmoto.DateCreated)
	}

	lastUpdatedDate, err := time.ParseInLocation("2006-01-02 15:04:05", hanmoto.DateModified, loc)
	if err != nil {
		log.Printf("Error in parsing timestamp: %s", hanmoto.DateModified)
	}

	var ccode string
	subjects := res.Onix.DescriptiveDetail.Subject
	if len(subjects) > 0 {
		ccode = subjects[0].SubjectCode
	}
	decoded, err := f.decoder.Decode(ccode)
	if err != nil {
		log.Println("Failed in decoding Ccode: ", err)
	}

	return &DetailedInformation{
		Author:          author,
		Publisher:       publisher,
		CreatedDate:     createdDate,
		LastUpdatedDate: lastUpdatedDate,
		Ccode:           ccode,
		Format:          decoded.Format,
		Target:          decoded.Target,
		Content:         decoded.Content,
	}, nil

}
