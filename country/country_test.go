package country

import (
	"fmt"
	"testing"
)

func TestGetAreaCodeByCode(t *testing.T) {
	fmt.Println("二字码数量：", len(twoCharCodeToCountryCodeMapping))
	fmt.Println("三字码数量：", len(threeCharCodeToCountryCodeMapping))
	fmt.Println("国家总数量：", len(countryCodeMapping))
	fmt.Println("国家区域总数量：", len(countryAreaMapping))
	for key, _ := range countryCodeMapping {
		areaCode, ok := countryAreaMapping[key]
		if !ok || areaCode == "" {
			fmt.Printf("国家：%v 不存在\n", key)
			continue
		}
	}
}

func TestGetCountryByCode(t *testing.T) {
	fmt.Println(GetAreaCodeByCode("156"))
}

func TestGetAreaNameByAreaCode(t *testing.T) {
	m := make(map[string]struct{})
	for _, areaCode := range countryAreaMapping {
		m[areaCode] = struct{}{}
	}

	fmt.Println("总区域熟练：", len(m))
}
