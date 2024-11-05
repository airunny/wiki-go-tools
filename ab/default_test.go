package ab

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
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
			Group: "b",
		},
		{
			Id:    "2",
			Group: "a",
		},
		{
			Id:    "3",
			Group: "b",
		},
		{
			Id:    "4",
			Group: "c",
		},
		{
			Id:    "5",
			Group: "a",
		},
		{
			Id:    "6",
			Group: "a",
		},
		{
			Id:    "7",
			Group: "a",
		},
		{
			Id:    "8",
			Group: "c",
		},
		{
			Id:    "9",
			Group: "c",
		},
		{
			Id:    "10",
			Group: "a",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.Group, ABTest(test.Id, groups), test.Id)
	}
}

func TestStringToMD5Int(t *testing.T) {
	hash := md5.New()
	hash.Write([]byte("1"))
	hashBytes := hash.Sum(nil)
	hexDigest := hex.EncodeToString(hashBytes)
	decimalValue := new(big.Int)
	decimalValue, success := decimalValue.SetString(hexDigest, 16)
	if !success {
		return
	}

	num2 := new(big.Int)
	num2.SetInt64(int64(3))

	result := new(big.Int)
	result.Mod(decimalValue, num2)
	fmt.Println("内容：", result.Uint64())
}
