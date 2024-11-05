package ab

import (
	"crypto/md5"
	"math/big"
)

// AB 测试分组，目前没有存储，只是对入参id分摊到groups中
func AB(id string, groups []string, opts ...Option) string {
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
	groupIndex := int(hashInt % uint64(len(groups)))
	return groups[groupIndex]
}

func stringToMD5Int(str string) uint64 {
	hash := md5.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	hashBigInt := new(big.Int).SetBytes(hashBytes)
	return hashBigInt.Uint64()
}
