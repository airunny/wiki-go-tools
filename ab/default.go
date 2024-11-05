package ab

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
)

// ABTest 测试分组，目前没有存储，只是对入参id分摊到groups中
func ABTest(id string, groups []string, opts ...Option) string {
	o := newDefaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if group, ok := o.fixedId[id]; ok {
		return group
	}

	// 获取用户ID的 MD5 数值
	hashInt := stringToMD5Int(id)
	// 计算分组索引
	num2 := new(big.Int)
	num2.SetInt64(int64(len(groups)))

	result := new(big.Int)
	result.Mod(hashInt, num2)
	groupIndex := int(result.Uint64())
	return groups[groupIndex]
}

func stringToMD5Int(str string) *big.Int {
	hash := md5.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	hexDigest := hex.EncodeToString(hashBytes)
	decimalValue := new(big.Int)
	decimalValue, success := decimalValue.SetString(hexDigest, 16)
	if !success {
		return new(big.Int)
	}
	return decimalValue
}
