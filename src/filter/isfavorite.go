package filter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tatamiya/new-books-notification/src/models"
)

type FavoriteFilter struct {
	FavoriteCategories []string `json:"categories"`
	FavoriteContents   []string `json:"contents"`
}

func NewFavoriteFilter(filterPath string) (*FavoriteFilter, error) {

	var favFilter FavoriteFilter
	filterData, ioErr := ioutil.ReadFile(filterPath)
	if ioErr != nil {
		return nil, fmt.Errorf("Could not read favorites.json!: %s", ioErr)
	}
	jsonErr := json.Unmarshal(filterData, &favFilter)
	if jsonErr != nil {
		return nil, fmt.Errorf("Could not unmarshal json data!: %s", jsonErr)
	}

	return &favFilter, nil
}

func (f *FavoriteFilter) IsFavorite(book *models.Book) bool {
	var category string
	if len(book.Categories) > 0 {
		category = book.Categories[0]
	}

	for _, favCategory := range f.FavoriteCategories {
		if category == favCategory {
			return true
		}

	}

	content := book.Content
	for _, favContent := range f.FavoriteContents {
		if content == favContent {
			return true
		}
	}

	return false

}
