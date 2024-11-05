package ab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAB(t *testing.T) {
	groups := []string{"a", "b", "c"}
	tests := []struct {
		Id    string
		Group string
	}{
		{
			Id:    "1",
			Group: "c",
		},
		{
			Id:    "2",
			Group: "a",
		},
		{
			Id:    "3",
			Group: "c",
		},
		{
			Id:    "4",
			Group: "b",
		},
		{
			Id:    "5",
			Group: "c",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Group, AB(test.Id, groups, WithFixedId("c", "1")), test.Id)
	}
}
