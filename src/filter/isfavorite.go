package filter

import "github.com/tatamiya/new-books-notification/src/models"

type FavoriteFilter struct {
	FavoriteCategories []string
	FavoriteGenres     []string
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

	return false

}
