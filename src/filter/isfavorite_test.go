package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/new-books-notification/src/models"
)

var sampleFavoriteFilter = FavoriteFilter{
	FavoriteCategories: []string{"自然科学", "コンピュータ"},
	FavoriteGenres:     []string{"物理学", "地理"},
}

func TestReturnTrueForFavoriteCategory(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{"自然科学"},
		Ccode:      "1040",
		Target:     "教養",
		Format:     "単行本",
		Genre:      "自然科学総記",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, true, isfav)

}

func TestReturnFalseForNotFavoriteCategory(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{"趣味・実用"},
		Ccode:      "1040",
		Target:     "教養",
		Format:     "単行本",
		Genre:      "自然科学総記",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, false, isfav)

}

func TestReturnFalseForEmptyCategory(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{},
		Ccode:      "1040",
		Target:     "教養",
		Format:     "単行本",
		Genre:      "自然科学総記",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, false, isfav)

}

func TestReturnTrueForFavoriteGenre(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{},
		Ccode:      "1042",
		Target:     "教養",
		Format:     "単行本",
		Genre:      "物理学",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, true, isfav)

}

func TestReturnFalseForNotFavoriteGenre(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{"工業・工学"},
		Ccode:      "1058",
		Target:     "教養",
		Format:     "単行本",
		Genre:      "その他の工業",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, false, isfav)

}

func TestReturnFalseForEmptyGenre(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{"工業・工学"},
		Ccode:      "",
		Target:     "",
		Format:     "",
		Genre:      "",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, false, isfav)

}
