package igorm

import (
	"fmt"
	"strings"
)

// Element Element
type Element struct {
	Key       string
	Condition string
	Value     interface{}
}

type Filter []Element

func (f Filter) QueryAndArgs() (query interface{}, args []interface{}) {
	fields := make([]string, 0, len(f))
	args = make([]interface{}, 0, len(f))

	for _, element := range f {
		condition := "="
		if element.Condition != "" {
			condition = element.Condition
		}
		fields = append(fields, fmt.Sprintf("%v %v ?", element.Key, condition))
		args = append(args, element.Value)
	}
	return strings.Join(fields, " and "), args
}

func (f Filter) SQL() string {
	fields := make([]string, 0, len(f))

	for _, element := range f {
		condition := "="
		if element.Condition != "" {
			condition = element.Condition
		}
		fields = append(fields, fmt.Sprintf("%v %v %v", element.Key, condition, element.Value))
	}

	return strings.Join(fields, " and ")
}
