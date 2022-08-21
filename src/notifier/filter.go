package notifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/tatamiya/new-books-notification/src/models"
)

type NotificationFilter struct {
	conditionBlocks []*conditionBlock
}

func (cf *NotificationFilter) IsFavorite(book *models.Book) bool {
	if len(cf.conditionBlocks) == 0 {
		return false
	}
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

type containCondition struct {
	filterBy string
	words    []string
}

func (c *containCondition) match(book *models.Book) bool {
	bookValue := reflect.ValueOf(*book)
	targetFieldValue := bookValue.FieldByName(c.filterBy)

	if !targetFieldValue.IsValid() {
		return false
	}

	for _, favWord := range c.words {
		if targetFieldValue.Interface() == favWord {
			return true
		}
	}

	return false
}

type notContainCondition struct {
	filterBy string
	words    []string
}

func (c *notContainCondition) match(book *models.Book) bool {
	bookValue := reflect.ValueOf(*book)
	targetFieldValue := bookValue.FieldByName(c.filterBy)

	if !targetFieldValue.IsValid() {
		return false
	}

	notContain := true
	for _, unfavWord := range c.words {
		if targetFieldValue.Interface() == unfavWord {
			notContain = false
		}
	}

	return notContain
}

type notStartWithCondition struct {
	filterBy string
	words    []string
}

func (c *notStartWithCondition) match(book *models.Book) bool {
	bookValue := reflect.ValueOf(*book)
	targetFieldValue := bookValue.FieldByName(c.filterBy)

	if !targetFieldValue.IsValid() {
		return false
	}

	notStartWith := true
	for _, unfavWord := range c.words {
		if strings.HasPrefix(targetFieldValue.Interface().(string), unfavWord) {
			notStartWith = false
		}
	}

	return notStartWith
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

func NewNotificationFilter(filterPath string) (*NotificationFilter, error) {

	var settings filterSettings
	filterData, ioErr := ioutil.ReadFile(filterPath)
	if ioErr != nil {
		return nil, fmt.Errorf("could not read %s!: %s", filterPath, ioErr)
	}
	jsonErr := json.Unmarshal(filterData, &settings)
	if jsonErr != nil {
		return nil, fmt.Errorf("could not unmarshal json data!: %s", jsonErr)
	}

	return buildNotificationFilter(&settings), nil
}

func buildNotificationFilter(settings *filterSettings) *NotificationFilter {
	var blocks []*conditionBlock
	for _, filterBlock := range settings.Blocks {
		var conditions []condition
		for _, filterCondition := range filterBlock.Conditions {
			filterBy := strings.Title(filterCondition.FilterBy)
			if isValidFieldName(filterBy) && filterCondition.FilterType == "contain" {
				conditions = append(conditions, &containCondition{
					filterBy: filterBy,
					words:    filterCondition.Words,
				})
			}
			if isValidFieldName(filterBy) && filterCondition.FilterType == "not_contain" {
				conditions = append(conditions, &notContainCondition{
					filterBy: filterBy,
					words:    filterCondition.Words,
				})
			}
			if isValidFieldName(filterBy) && filterCondition.FilterType == "not_start_with" {
				conditions = append(conditions, &notStartWithCondition{
					filterBy: filterBy,
					words:    filterCondition.Words,
				})
			}
		}
		if len(conditions) > 0 {
			blocks = append(blocks, &conditionBlock{conditions: conditions})
		}
	}

	return &NotificationFilter{
		conditionBlocks: blocks,
	}
}

func isValidFieldName(fieldName string) bool {
	bookValue := reflect.ValueOf(models.Book{})
	return bookValue.FieldByName(fieldName).IsValid()
}
