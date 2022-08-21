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
					&containCondition{filterBy: "Categories", words: []string{"自然科学"}},
					&containCondition{filterBy: "Content", words: []string{"数学", "物理学"}},
				},
			},
			{
				conditions: []condition{
					&notContainCondition{filterBy: "Categories", words: []string{"学参"}},
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
					{
						FilterBy:   "INVALID",
						FilterType: "contain",
						Words:      []string{"hogehoge"},
					},
				},
			},
		},
	}

	expectedNotificationFilter := NotificationFilter{
		conditionBlocks: []*conditionBlock{
			{
				conditions: []condition{
					&containCondition{filterBy: "Categories", words: []string{"自然科学"}},
					&containCondition{filterBy: "Content", words: []string{"数学", "物理学"}},
				},
			},
		},
	}

	actualNotificationFilter := buildNotificationFilter(&inputFilterSettings)

	assert.EqualValues(t, expectedNotificationFilter, *actualNotificationFilter)

}

func TestContainConditionFiltersCorrectly(t *testing.T) {
	// Filter by Categories
	categoriesCondition := containCondition{
		filterBy: "Categories",
		words:    []string{"自然科学"},
	}
	bookWithFavoriteCategory := models.Book{
		Categories: "自然科学",
	}
	assert.Equal(t, true, categoriesCondition.match(&bookWithFavoriteCategory))
	bookWithUnfavoriteCategory := models.Book{
		Categories: "児童書",
	}
	assert.Equal(t, false, categoriesCondition.match(&bookWithUnfavoriteCategory))

	// Filter by Content
	contentCondition := containCondition{
		filterBy: "Content",
		words:    []string{"数学", "物理学"},
	}
	bookWithFavoriteContent := models.Book{
		Content: "物理学",
	}
	assert.Equal(t, true, contentCondition.match(&bookWithFavoriteContent))

	// Filter by Target
	targetCondition := containCondition{
		filterBy: "Target",
		words:    []string{"専門", "一般"},
	}
	bookWithFavoriteTarget := models.Book{
		Target: "専門",
	}
	assert.Equal(t, true, targetCondition.match(&bookWithFavoriteTarget))
}

func TestContainConditionWithInvalidTargetFieldAlwaysReturnsFalse(t *testing.T) {
	categoriesCondition := containCondition{
		filterBy: "INVALID_FIELD_NAME",
		words:    []string{"自然科学"},
	}

	bookWithFavoriteCategory := models.Book{
		Categories: "自然科学",
	}
	assert.Equal(t, false, categoriesCondition.match(&bookWithFavoriteCategory))

	bookWithUnfavoriteCategory := models.Book{
		Categories: "児童書",
	}
	assert.Equal(t, false, categoriesCondition.match(&bookWithUnfavoriteCategory))
}

func TestNotContainConditionFiltersCorrectly(t *testing.T) {
	// Filter by Categories
	testCondition := notContainCondition{
		filterBy: "Categories",
		words:    []string{"児童書", "学参"},
	}

	bookWithFavoriteCategory := models.Book{
		Categories: "自然科学",
	}
	assert.Equal(t, true, testCondition.match(&bookWithFavoriteCategory))

	bookWithUnfavoriteCategory := models.Book{
		Categories: "児童書",
	}
	assert.Equal(t, false, testCondition.match(&bookWithUnfavoriteCategory))
}

func TestNotStartWithConditionFiltersCorrectly(t *testing.T) {
	// Filter by Target
	testCondition := notStartWithCondition{
		filterBy: "Target",
		words:    []string{"学参"},
	}

	bookWithFavoriteTarget := models.Book{
		Target: "専門",
	}
	assert.Equal(t, true, testCondition.match(&bookWithFavoriteTarget))

	bookWithUnfavoriteTarget := models.Book{
		Target: "学参II（高校）",
	}
	assert.Equal(t, false, testCondition.match(&bookWithUnfavoriteTarget))
}

func TestNotificationFilterForOrConditions(t *testing.T) {

	sampleFilter := NotificationFilter{
		conditionBlocks: []*conditionBlock{
			{
				conditions: []condition{
					&containCondition{filterBy: "Categories", words: []string{"自然科学"}},
					&containCondition{filterBy: "Content", words: []string{"数学", "物理学"}},
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
					&containCondition{filterBy: "Categories", words: []string{"自然科学"}},
				},
			},
			{
				conditions: []condition{
					&containCondition{filterBy: "Content", words: []string{"数学", "物理学"}},
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
					&containCondition{filterBy: "Categories", words: []string{"自然科学"}},
					&containCondition{filterBy: "Content", words: []string{"数学", "物理学"}},
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
