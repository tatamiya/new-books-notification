package notifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/new-books-notification/src/models"
)

var sampleFavoriteFilter = FavoriteFilter{
	FavoriteCategories: []string{"自然科学", "コンピュータ"},
	FavoriteContents:   []string{"物理学", "地理"},
}

func TestReturnTrueForFavoriteCategory(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{"自然科学"},
		Ccode:      "1040",
		Target:     "教養",
		Format:     "単行本",
		Content:    "自然科学総記",
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
		Content:    "自然科学総記",
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
		Content:    "自然科学総記",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, false, isfav)

}

func TestReturnTrueForFavoriteContent(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{},
		Ccode:      "1042",
		Target:     "教養",
		Format:     "単行本",
		Content:    "物理学",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, true, isfav)

}

func TestReturnFalseForNotFavoriteContent(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{"工業・工学"},
		Ccode:      "1058",
		Target:     "教養",
		Format:     "単行本",
		Content:    "その他の工業",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, false, isfav)

}

func TestReturnFalseForEmptyContent(t *testing.T) {

	sampleBook := models.Book{
		Categories: []string{"工業・工学"},
		Ccode:      "",
		Target:     "",
		Format:     "",
		Content:    "",
	}

	isfav := sampleFavoriteFilter.IsFavorite(&sampleBook)

	assert.Equal(t, false, isfav)

}

func TestNewFavoriteFilter(t *testing.T) {

	favFilter, err := NewFavoriteFilter("./sample_favorites.json")
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"コンピュータ"}, favFilter.FavoriteCategories)
	assert.EqualValues(t, []string{"情報科学", "電子通信"}, favFilter.FavoriteContents)
}

func TestLoadNewComplexFilterSuccessfuly(t *testing.T) {

	actualComplexFilter, err := NewComplexFilter("./test_complex_filter.json")
	expectedComplexFilter := ComplexFilter{
		conditionBlocks: []*conditionBlock{
			{
				conditions: []condition{
					&categoryContainsCondition{words: []string{"自然科学"}},
					&contentContainsCondition{words: []string{"数学", "物理学"}},
				},
			},
		},
	}

	assert.EqualValues(t, expectedComplexFilter, *actualComplexFilter)
	assert.Nil(t, err)
}

func TestBuildComplexFilter(t *testing.T) {
	inputFilterSettings := filterSettings{
		Blocks: []filterBlocks{
			{
				Conditions: []filterCondition{
					{
						FilterBy:   "category",
						FilterType: "contain",
						Words:      []string{"自然科学"},
					},
					{
						FilterBy:   "content",
						FilterType: "contain",
						Words:      []string{"数学", "物理学"},
					},
				},
			},
		},
	}

	expectedComplexFilter := ComplexFilter{
		conditionBlocks: []*conditionBlock{
			{
				conditions: []condition{
					&categoryContainsCondition{words: []string{"自然科学"}},
					&contentContainsCondition{words: []string{"数学", "物理学"}},
				},
			},
		},
	}

	actualComplexFilter := buildComplexFilter(&inputFilterSettings)

	assert.EqualValues(t, expectedComplexFilter, *actualComplexFilter)

}

func TestComplexFilterForOrConditions(t *testing.T) {

	sampleFilter := ComplexFilter{
		conditionBlocks: []*conditionBlock{
			{
				conditions: []condition{
					&categoryContainsCondition{words: []string{"自然科学"}},
					&contentContainsCondition{words: []string{"数学", "物理学"}},
				},
			},
		},
	}

	bookWithFavoriteCategoryAndContent := models.Book{
		Categories: []string{"自然科学"},
		Content:    "物理学",
	}
	assert.Equal(t, true, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndContent))

	bookWithFavoriteCategoryAndUnfavoriteContent := models.Book{
		Categories: []string{"自然科学"},
		Content:    "その他の工業",
	}
	assert.Equal(t, true, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndUnfavoriteContent))

	bookWithFavoriteContentAndUnfavoriteCategory := models.Book{
		Categories: []string{"児童書"},
		Content:    "物理学",
	}
	assert.Equal(t, true, sampleFilter.IsFavorite(&bookWithFavoriteContentAndUnfavoriteCategory))

	bookWithUnfavoriteCategoryAndContent := models.Book{
		Categories: []string{"児童書"},
		Content:    "その他の工業",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithUnfavoriteCategoryAndContent))
}

func TestComplexFilterForAndConditions(t *testing.T) {

	sampleFilter := ComplexFilter{
		conditionBlocks: []*conditionBlock{
			{
				conditions: []condition{
					&categoryContainsCondition{words: []string{"自然科学"}},
				},
			},
			{
				conditions: []condition{
					&contentContainsCondition{words: []string{"数学", "物理学"}},
				},
			},
		},
	}

	bookWithFavoriteCategoryAndContent := models.Book{
		Categories: []string{"自然科学"},
		Content:    "物理学",
	}
	assert.Equal(t, true, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndContent))

	bookWithFavoriteCategoryAndUnfavoriteContent := models.Book{
		Categories: []string{"自然科学"},
		Content:    "その他の工業",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndUnfavoriteContent))

	bookWithFavoriteContentAndUnfavoriteCategory := models.Book{
		Categories: []string{"児童書"},
		Content:    "物理学",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithFavoriteContentAndUnfavoriteCategory))

	bookWithUnfavoriteCategoryAndContent := models.Book{
		Categories: []string{"児童書"},
		Content:    "その他の工業",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithUnfavoriteCategoryAndContent))
}
