package i18n

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFromZip(t *testing.T) {
	_, err := ReadFromZip("/Users/liyanbing/Downloads/language.zip")
	assert.Nil(t, err)
	str, _ := json.Marshal(languages)
	fmt.Println(string(str))
}
