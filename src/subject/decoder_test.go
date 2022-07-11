package subject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDecoder = SubjectDecoder{
	Taishou: map[string]string{
		"0": "一般",
		"1": "教養",
		"3": "専門書",
	},
	Keitai: map[string]string{
		"0": "単行本",
		"1": "文庫",
	},
	Naiyou: map[string]string{
		"40": "自然科学総記",
		"42": "物理学",
	},
}

func TestDecodeSubjectCorrectly(t *testing.T) {

	expectedDecoded := DecodedSubject{
		Ccode:  "0040",
		Target: "一般",
		Format: "単行本",
		Genre:  "自然科学総記",
	}

	actualDecoded, err := testDecoder.decode("0040")
	assert.Nil(t, err)
	assert.EqualValues(t, expectedDecoded, *actualDecoded)
}

func TestDecodingFailsWhenCcodeIsNotFourDigits(t *testing.T) {

	actualDecoded, err := testDecoder.decode("040")
	assert.NotNil(t, err)
	assert.Nil(t, actualDecoded)
}

func TestDecodingFailsWhenCcodeContainsNonDigit(t *testing.T) {

	actualDecoded, err := testDecoder.decode("1a40")
	assert.NotNil(t, err)
	assert.Nil(t, actualDecoded)
}

func TestDecodedResultIsEmptyWhenNotFound(t *testing.T) {

	expectedDecoded := DecodedSubject{
		Ccode:  "0049",
		Target: "一般",
		Format: "単行本",
		Genre:  "",
	}

	actualDecoded, err := testDecoder.decode("0049")
	assert.Nil(t, err)
	assert.EqualValues(t, expectedDecoded, *actualDecoded)
}

func TestNewSubjectDecoder(t *testing.T) {
	decoder, err := NewSubjectDecoder()
	assert.Nil(t, err)
	assert.EqualValues(t, "一般", decoder.Taishou["0"])
	assert.EqualValues(t, "単行本", decoder.Keitai["0"])
	assert.EqualValues(t, "総記", decoder.Naiyou["00"])
}
