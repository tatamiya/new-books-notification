package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/mmcdole/gofeed"
	"github.com/tatamiya/new-books-notification/src/filter"
	"github.com/tatamiya/new-books-notification/src/models"
	"github.com/tatamiya/new-books-notification/src/notifier"
	"github.com/tatamiya/new-books-notification/src/openbd"
	"github.com/tatamiya/new-books-notification/src/subject"
)

func main() {
	url := "https://www.hanmoto.com/ci/bd/search/hdt/%E6%96%B0%E3%81%97%E3%81%8F%E7%99%BB%E9%8C%B2%E3%81%95%E3%82%8C%E3%81%9F%E6%9C%AC/sdate/today/created/today/order/desc/vw/rss20"
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)

	// gopath := os.Getenv("GOPATH")
	// rootPath := gopath + "/src/github.com/tatamiya/new-books-notification/"
	rootPath := "./"

	jsonFilePath := rootPath + "src/subject/ccode.json"
	subjectDecoder, err := subject.NewSubjectDecoder(jsonFilePath)
	if err != nil {
		log.Println("Error in loading SubjectDecoder.")
		panic(err)
	}

	bookList := models.NewBookListFromFeed(feed)
	log.Println(bookList.UploadDate.String())

	favoritesPath := rootPath + "favorites.json"
	favFilter, err := filter.NewFavoriteFilter(favoritesPath)
	if err != nil {
		log.Println("Error in loading Favorite Filter.")
		panic(err)
	}

	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	slackNotifier, notifierErr := notifier.NewSlackNotifier(webhookURL)
	fmt.Println(notifierErr)
	if notifierErr != nil {
		log.Println("Error in loading SlackNotifier.")
	}

	var wg sync.WaitGroup
	for _, book := range bookList.Books {
		wg.Add(1)
		go func(book *models.Book) {

			defer wg.Done()

			openBDResp, err := openbd.GetDetailedBookInfo(book.Isbn)

			book.UpdateInfoFrom(openBDResp)
			book.UpdateSubject(subjectDecoder)

			log.Println(book)
			if favFilter.IsFavorite(book) {
				err = slackNotifier.Post(book.AsNotificationMessage())
				if err != nil {
					log.Printf("Error in notifying %s(%s) to Slack: %s\n", book.Isbn, strings.TrimSpace(book.Title), err)
				}
			}

		}(book)

	}
	wg.Wait()

}
