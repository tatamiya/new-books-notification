package models

import (
	"google.golang.org/genproto/googleapis/type/datetime"
)

type BookList struct {
	UploadDate datetime.DateTime
	Books      []Book
}

type Book struct {
	Isbn            string
	Title           string
	Url             string
	Authors         string
	Publisher       string
	Categories      []string
	Ccode           string
	Target          string
	Format          string
	Genre           string
	PubDate         datetime.DateTime
	CreatedDate     datetime.DateTime
	LastUpdatedDate datetime.DateTime
}
