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
	groups := []string{"A", "B"}
	tests := []struct {
		Id    string
		Group string
	}{
		{
			Id:    "1",
			Group: "B",
		},
		{
			Id:    "2",
			Group: "A",
		},
		{
			Id:    "3",
			Group: "B",
		},
		{
			Id:    "4",
			Group: "A",
		},
		{
			Id:    "5",
			Group: "B",
		},
		{
			Id:    "6",
			Group: "A",
		},
		{
			Id:    "7",
			Group: "B",
		},
		{
			Id:    "8",
			Group: "B",
		},
		{
			Id:    "9",
			Group: "A",
		},
		{
			Id:    "10",
			Group: "A",
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
