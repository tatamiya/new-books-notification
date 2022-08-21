package notifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tatamiya/new-books-notification/src/models"
)

func TestLoadNewNotificationFilterSuccessfuly(t *testing.T) {

	actualNotificationFilter, err := NewNotificationFilter("./test_notification_filter.json")
	expectedNotificationFilter := NotificationFilter{
		conditionBlocks: []*conditionBlock{
			{
				conditions: []condition{
					&categoryContainsCondition{words: []string{"自然科学"}},
					&contentContainsCondition{words: []string{"数学", "物理学"}},
				},
			},
		},
	}

	assert.EqualValues(t, expectedNotificationFilter, *actualNotificationFilter)
	assert.Nil(t, err)
}

func TestBuildNotificationFilter(t *testing.T) {
	inputFilterSettings := filterSettings{
		Blocks: []filterBlocks{
			{
				Conditions: []filterCondition{
					{
						FilterBy:   "categories",
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

	expectedNotificationFilter := NotificationFilter{
		conditionBlocks: []*conditionBlock{
			{
				conditions: []condition{
					&categoryContainsCondition{words: []string{"自然科学"}},
					&contentContainsCondition{words: []string{"数学", "物理学"}},
				},
			},
		},
	}

	actualNotificationFilter := buildNotificationFilter(&inputFilterSettings)

	assert.EqualValues(t, expectedNotificationFilter, *actualNotificationFilter)

}

func TestNotificationFilterForOrConditions(t *testing.T) {

	sampleFilter := NotificationFilter{
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

func TestNotificationFilterForAndConditions(t *testing.T) {

	sampleFilter := NotificationFilter{
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

	sampleFilter := NotificationFilter{}

	bookWithFavoriteCategoryAndContent := models.Book{
		Categories: "自然科学",
		Content:    "物理学",
	}
	assert.Equal(t, false, sampleFilter.IsFavorite(&bookWithFavoriteCategoryAndContent))
}

func TestFilterReturnsFalseForEmptyCategoryAndContent(t *testing.T) {
	sampleFilter := NotificationFilter{
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
