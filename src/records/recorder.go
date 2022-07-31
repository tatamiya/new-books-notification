package records

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/tatamiya/new-books-notification/src/models"
	"google.golang.org/api/iterator"
)

type Record struct {
	ISBN          string
	Title         string
	Url           string
	Authors       string
	Publisher     string
	Categories    string
	Ccode         string
	Target        string
	Format        string
	Content       string
	PubDate       civil.Date
	CreatedAt     time.Time
	LastUpdatedAt time.Time
	UploadedAt    time.Time
	UploadedDate  civil.Date
}

func (s *BQRecorder) CreateTable(ctx context.Context) error {
	schema := bigquery.Schema{
		{Name: "ISBN", Required: true, Type: bigquery.StringFieldType},
		{Name: "PubDate", Required: true, Type: bigquery.DateFieldType},
		{Name: "Title", Required: true, Type: bigquery.StringFieldType},
		{Name: "Url", Required: true, Type: bigquery.StringFieldType},
		{Name: "Authors", Required: false, Type: bigquery.StringFieldType},
		{Name: "Publisher", Required: false, Type: bigquery.StringFieldType},
		{Name: "Categories", Required: false, Type: bigquery.StringFieldType},
		{Name: "Ccode", Required: false, Type: bigquery.StringFieldType},
		{Name: "Target", Required: false, Type: bigquery.StringFieldType},
		{Name: "Format", Required: false, Type: bigquery.StringFieldType},
		{Name: "Content", Required: false, Type: bigquery.StringFieldType},
		{Name: "CreatedAt", Required: false, Type: bigquery.TimestampFieldType},
		{Name: "LastUpdatedAt", Required: false, Type: bigquery.TimestampFieldType},
		{Name: "UploadedAt", Required: true, Type: bigquery.TimestampFieldType},
		{Name: "UploadedDate", Required: true, Type: bigquery.DateFieldType},
	}
	metadata := bigquery.TableMetadata{
		Schema: schema,
		TimePartitioning: &bigquery.TimePartitioning{
			Type:  bigquery.DayPartitioningType,
			Field: "UploadedDate",
		},
	}

	return s.table.Create(ctx, &metadata)
}

func convertIntoRecordsFrom(bookList *models.BookList) []*Record {
	uploadedAt := bookList.UploadDate

	var records []*Record
	for _, book := range bookList.Books {
		bookRecord := convertIntoRecord(book, uploadedAt)
		records = append(records, bookRecord)
	}
	return records
}

func convertIntoRecord(book *models.Book, uploadedAt time.Time) *Record {

	var category string
	if len(book.Categories) > 0 {
		category = book.Categories[0]
	} else {
		category = ""
	}
	return &Record{
		ISBN:          book.Isbn,
		Title:         book.Title,
		Url:           book.Url,
		Authors:       book.Authors,
		Publisher:     book.Publisher,
		Categories:    category,
		Ccode:         book.Ccode,
		Target:        book.Target,
		Format:        book.Format,
		Content:       book.Content,
		PubDate:       civil.DateOf(book.PubDate),
		CreatedAt:     book.CreatedDate,
		LastUpdatedAt: book.LastUpdatedDate,
		UploadedAt:    uploadedAt,
		UploadedDate:  civil.DateOf(uploadedAt),
	}
}

type BQRecorder struct {
	client *bigquery.Client
	table  *bigquery.Table
}

type BQSettings struct {
	ProjectID   string
	DatasetName string
	TableName   string
}

func NewBQRecorder(ctx context.Context, settings *BQSettings) (*BQRecorder, error) {

	client, err := bigquery.NewClient(ctx, settings.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to BigQuery: %s", err)
	}
	dataset := client.Dataset(settings.DatasetName)
	table := dataset.Table(settings.TableName)

	recorder := BQRecorder{client: client, table: table}

	_, err = table.Metadata(ctx)
	if err != nil {
		log.Printf("Cannot find the table %s: %s", settings.TableName, err)
		if err = recorder.CreateTable(ctx); err != nil {
			return nil, fmt.Errorf("Cannot create a table: %s", err)
		}
		log.Printf("Successfully created the table %s", settings.TableName)
	}

	return &recorder, nil
}

func (s *BQRecorder) SaveRecords(ctx context.Context, bookList *models.BookList) error {
	records := convertIntoRecordsFrom(bookList)

	inserter := s.table.Inserter()
	err := inserter.Put(ctx, records)
	if err != nil {
		return fmt.Errorf("Upload book records failed: %s", err)
	}

	return nil
}

func (s *BQRecorder) GetRecordedISBN(ctx context.Context, targetDate time.Time) ([]string, error) {

	table := s.table
	fullTableID := fmt.Sprintf("`%s.%s.%s`", table.ProjectID, table.DatasetID, table.TableID)
	uploadedDate := fmt.Sprintf("\"%s\"", targetDate.Format("2006-01-02"))
	q := s.client.Query(`SELECT DISTINCT ISBN FROM ` + fullTableID + ` WHERE UploadedDate=` + uploadedDate)

	it, err := q.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("Query execution failed: %s", err)
	}
	type QueryResult struct {
		ISBN string
	}
	uploadedISBN := []string{}
	for {
		var r QueryResult
		err := it.Next(&r)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Unexpected query results: %s", err)
			continue
		}
		uploadedISBN = append(uploadedISBN, r.ISBN)
	}

	return uploadedISBN, nil
}
