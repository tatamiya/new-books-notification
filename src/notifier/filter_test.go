package notifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/new-books-notification/src/models"
)

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
		Categories: "自然科学",
		Content:    "物理学",
	}
	assert.Equal(t, true, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndContent))

	bookWithFavoriteCategoryAndUnfavoriteContent := models.Book{
		Categories: "自然科学",
		Content:    "その他の工業",
	}
	assert.Equal(t, true, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndUnfavoriteContent))

	bookWithFavoriteContentAndUnfavoriteCategory := models.Book{
		Categories: "児童書",
		Content:    "物理学",
	}
	assert.Equal(t, true, sampleFilter.IsFavorite(&bookWithFavoriteContentAndUnfavoriteCategory))

	bookWithUnfavoriteCategoryAndContent := models.Book{
		Categories: "児童書",
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
		Categories: "自然科学",
		Content:    "物理学",
	}
	assert.Equal(t, true, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndContent))

	bookWithFavoriteCategoryAndUnfavoriteContent := models.Book{
		Categories: "自然科学",
		Content:    "その他の工業",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndUnfavoriteContent))

	bookWithFavoriteContentAndUnfavoriteCategory := models.Book{
		Categories: "児童書",
		Content:    "物理学",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithFavoriteContentAndUnfavoriteCategory))

	bookWithUnfavoriteCategoryAndContent := models.Book{
		Categories: "児童書",
		Content:    "その他の工業",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithUnfavoriteCategoryAndContent))
}

func TestEmptyFilterDoesNotFail(t *testing.T) {

	sampleFilter := ComplexFilter{}

	bookWithFavoriteCategoryAndContent := models.Book{
		Categories: "自然科学",
		Content:    "物理学",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndContent))
}

func TestFilterReturnsFalseForEmptyCategoryAndContent(t *testing.T) {
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

	bookWithEmptyCategoryAndContent := models.Book{}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithEmptyCategoryAndContent))

	bookWithEmptyContent := models.Book{
		Categories: "hoge",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithEmptyContent))

	bookWithEmptyCategory := models.Book{
		Content: "hoge",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithEmptyCategory))
}
