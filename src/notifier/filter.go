package notifier

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

	category := book.Categories
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

type ComplexFilter struct {
	conditionBlocks []*conditionBlock
}

func (cf *ComplexFilter) IsFavorite(book *models.Book) bool {
	matchAll := true
	for _, conditionBlock := range cf.conditionBlocks {
		if !conditionBlock.matchAny(book) {
			matchAll = false
			break
		}
	}
	return matchAll
}

type conditionBlock struct {
	conditions []condition
}

type condition interface {
	match(*models.Book) bool
}

func (cb *conditionBlock) matchAny(book *models.Book) bool {
	match := false
	for _, condition := range cb.conditions {
		if condition.match(book) {
			match = true
			break
		}
	}
	return match
}

type categoryContainsCondition struct {
	words []string
}

func (c *categoryContainsCondition) match(book *models.Book) bool {

	for _, favCategory := range c.words {
		if book.Categories == favCategory {
			return true
		}
	}

	return false
}

type contentContainsCondition struct {
	words []string
}

func (c *contentContainsCondition) match(book *models.Book) bool {
	content := book.Content
	for _, favContent := range c.words {
		if content == favContent {
			return true
		}
	}

	return false
}

type filterSettings struct {
	Blocks []filterBlocks `json:"blocks"`
}

type filterBlocks struct {
	Conditions []filterCondition `json:"conditions"`
}
type filterCondition struct {
	FilterBy   string   `json:"filter_by"`
	FilterType string   `json:"type"`
	Words      []string `json:"words"`
}

func NewComplexFilter(filterPath string) (*ComplexFilter, error) {

	var settings filterSettings
	filterData, ioErr := ioutil.ReadFile(filterPath)
	if ioErr != nil {
		return nil, fmt.Errorf("could not read %s!: %s", filterPath, ioErr)
	}
	jsonErr := json.Unmarshal(filterData, &settings)
	if jsonErr != nil {
		return nil, fmt.Errorf("could not unmarshal json data!: %s", jsonErr)
	}

	return buildComplexFilter(&settings), nil
}

func buildComplexFilter(settings *filterSettings) *ComplexFilter {
	var blocks []*conditionBlock
	for _, filterBlock := range settings.Blocks {
		var conditions []condition
		for _, filterCondition := range filterBlock.Conditions {
			if filterCondition.FilterBy == "content" && filterCondition.FilterType == "contain" {
				conditions = append(conditions, &contentContainsCondition{
					words: filterCondition.Words,
				})
			} else if filterCondition.FilterBy == "category" && filterCondition.FilterType == "contain" {
				conditions = append(conditions, &categoryContainsCondition{
					words: filterCondition.Words,
				})
			}
		}
		if len(conditions) > 0 {
			blocks = append(blocks, &conditionBlock{conditions: conditions})
		}
	}

	return &ComplexFilter{
		conditionBlocks: blocks,
	}
}
