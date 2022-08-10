package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/tatamiya/new-books-notification/src/config"
	"github.com/tatamiya/new-books-notification/src/models"
	"github.com/tatamiya/new-books-notification/src/notifier"
	"github.com/tatamiya/new-books-notification/src/openbd"
	"github.com/tatamiya/new-books-notification/src/recorder"
	"github.com/tatamiya/new-books-notification/src/subject"
	"github.com/tatamiya/new-books-notification/src/uploader"
)

type Uploader interface {
	Upload(*uploader.UploadObject) error
}

type Recorder interface {
	GetRecordedISBN(context.Context, time.Time) ([]string, error)
	SaveRecords(context.Context, *models.BookList) error
}

type Notifier interface {
	Post(string) error
}

type Filter interface {
	IsFavorite(*models.Book) bool
}

type DetailFetcher interface {
	FetchDetailInfo(string) (*openbd.OpenBDResponse, error)
}

func coreProcess(
	bookList *models.BookList,
	fetcher DetailFetcher,
	subjectDecoder *subject.SubjectDecoder,
	recorder Recorder,
	filter Filter,
	notifier Notifier,
) int {

	ctx := context.Background()

	var newBookList *models.BookList
	if recorder != nil {
		newBookList = bookList
	} else {
		uploadedISBN, err := recorder.GetRecordedISBN(ctx, bookList.UploadDate)
		if err != nil {
			log.Printf("Cannot fetch ISBNs of uploaded books from BigQuery: %s", err)
		}
		newBookList = bookList.FilterOut(uploadedISBN)
	}

	var wg sync.WaitGroup
	for _, book := range newBookList.Books {
		wg.Add(1)
		go func(book *models.Book) {

			defer wg.Done()

			openBDResp, err := fetcher.FetchDetailInfo(book.Isbn)
			if err != nil {
				log.Printf("Cannot fetch data from OpenBD (%s, %s): %s", book.Isbn, book.Title, err)
				return
			} else if openBDResp == nil {
				log.Printf("Response from OpenBD is empty (%s, %s)", book.Isbn, book.Title)
				return
			}

			book.UpdateInfoFrom(openBDResp)
			book.UpdateSubject(subjectDecoder)

			if filter.IsFavorite(book) {
				err = notifier.Post(book.AsNotificationMessage())
				if err != nil {
					log.Printf("Error in notifying %s(%s) to Slack: %s\n", book.Isbn, book.Title, err)
				}
			}

		}(book)

	}
	wg.Wait()

	err := recorder.SaveRecords(ctx, newBookList)
	if err != nil {
		log.Printf("Cannot save newly arrived book records: %s", err)
	}

	return len(newBookList.Books)

}

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(config.FeedURL)
	if err != nil {
		log.Println("Could not get feed!")
		panic(err)
	}

	bookList := models.NewBookListFromFeed(feed)
	log.Println(bookList.UploadDate.String())

	detailFetcher := openbd.NewOpenBDDetailFetcher()

	subjectDecoder, err := subject.NewSubjectDecoder(config.CcodeJsonFilePath)
	if err != nil {
		log.Println("Error in loading SubjectDecoder.")
		panic(err)
	}

	bqSettings := fetchBQSettings()
	ctx := context.Background()
	bqRecorder, err := recorder.NewBQRecorder(ctx, bqSettings)

	favFilter, err := notifier.NewFavoriteFilter(config.FilterSettingFilePath)
	if err != nil {
		log.Println("Error in loading Favorite Filter.")
		panic(err)
	}

	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	slackNotifier, notifierErr := notifier.NewSlackNotifier(webhookURL)
	if notifierErr != nil {
		log.Println("Error in loading SlackNotifier.")
	}

	numUploaded := coreProcess(bookList, detailFetcher, subjectDecoder, bqRecorder, favFilter, slackNotifier)

	log.Printf("Reported %d new book(s)", numUploaded)

	uploadFeed, err := generateJsonUploadObject(feed)
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	objectUploader, uploaderErr := uploader.NewGCSUploader(ctx, bucketName, "")
	if uploaderErr != nil {
		log.Printf("Cannot create feed uploader: %s", uploaderErr)
		return
	}
	uploadErr := objectUploader.Upload(uploadFeed)
	if uploadErr != nil {
		log.Printf("Feed upload failed: %s", err)
	}

	return
}

func generateJsonUploadObject(feed *gofeed.Feed) (*uploader.UploadObject, error) {
	b, err := json.Marshal(feed)
	if err != nil {

		return nil, fmt.Errorf("Failed in converting feed into JSON: %s", err)
	}
	feedJsonFilename := fmt.Sprintf("feed%s.json", feed.PublishedParsed.Format("20060102"))

	uploadObject := uploader.UploadObject{
		ObjectName:  feedJsonFilename,
		ContentType: "application/json",
		Binary:      b,
	}

	return &uploadObject, nil
}

func fetchBQSettings() *recorder.BQSettings {

	projectID, err := getProjectID()
	if err != nil || projectID == "" {
		log.Printf("Cannot get projectID from metadata API: %s", err)
		projectID = os.Getenv("GCP_PROJECT_ID")
	}

	return &recorder.BQSettings{
		ProjectID:   projectID,
		DatasetName: os.Getenv("GCP_BIGQUERY_DATASET"),
		TableName:   os.Getenv("GCP_BIGQUERY_TABLE"),
	}
}

func getProjectID() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://metadata.google.internal/computeMetadata/v1/project/project-id", nil)
	if err != nil {
		return "", fmt.Errorf("Could not create request: %s", err)
	}
	req.Header.Add("Metadata-Flavor", "Google")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Request Failed: %s", err)
	}

	bProjectID, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read metadata response: %s", err)
	}

	return string(bProjectID), err
}
