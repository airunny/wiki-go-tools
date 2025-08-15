package country

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"
)

var (
	countryFlagData      []countryFlag
	flagCodeMapping      map[string]*countryFlag
	flagTwoCharMapping   map[string]*countryFlag
	flagThreeCharMapping map[string]*countryFlag
	flagAreaCodeMapping  map[string]*countryFlag

	ErrCountryNotFound = errors.New("country not found")
	ErrEmptyFlagData   = errors.New("country flag data not loaded")
)

type countryFlag struct {
	TwoCharCode string            `json:"two_char_code"`
	CountryCode string            `json:"country_code"`
	AreaCode    string            `json:"area_code"`
	FlagURL     string            `json:"flag_url"`
	CountryName map[string]string `json:"country_name"`
}

func init() {
	flagPath := os.Getenv("COUNTRY_FLAG_PATH")
	if flagPath == "" {
		return
	}

	var filePath string

	// 检查 flagPath 是否是一个文件
	if info, err := os.Stat(flagPath); err == nil && !info.IsDir() {
		// 如果是文件，直接使用
		filePath = flagPath
	} else {
		// 如果是目录，拼接文件名
		flagFileName := os.Getenv("COUNTRY_FLAG_FILE_NAME")
		if flagFileName == "" {
			flagFileName = "countryFlag.json"
		}
		filePath = path.Join(flagPath, flagFileName)
	}

	err := loadCountryFlagData(filePath)
	if err != nil {
		panic(err)
	}
}

func loadCountryFlagData(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &countryFlagData)
	if err != nil {
		return err
	}

	initFlagMappings()
	return nil
}

func initFlagMappings() {
	flagCodeMapping = make(map[string]*countryFlag, len(countryFlagData))
	flagTwoCharMapping = make(map[string]*countryFlag, len(countryFlagData))
	flagThreeCharMapping = make(map[string]*countryFlag, len(countryFlagData))
	flagAreaCodeMapping = make(map[string]*countryFlag, len(countryFlagData))

	// 先建立国家码到旗帜数据的映射
	for i := range countryFlagData {
		flag := &countryFlagData[i]

		if flag.AreaCode != "" {
			areaCode := strings.TrimSpace(flag.AreaCode)
			flagAreaCodeMapping[areaCode] = flag

			if strings.HasPrefix(areaCode, "00") {
				flagAreaCodeMapping[areaCode[2:]] = flag
			}

			if strings.HasPrefix(areaCode, "+") {
				flagAreaCodeMapping[areaCode[1:]] = flag
			}
		}

		if flag.CountryCode != "" {
			flagCodeMapping[flag.CountryCode] = flag
		}

		if flag.TwoCharCode != "" {
			flagTwoCharMapping[strings.ToUpper(flag.TwoCharCode)] = flag
		}

	}

	// 利用现有的三字码映射建立三字码到旗帜数据的映射
	for threeChar, countryCode := range threeCharCodeToCountryCodeMapping {
		if flag, exists := flagCodeMapping[countryCode]; exists {
			flagThreeCharMapping[threeChar] = flag
		}
	}
}

func GetCountryInfoWithUserCountry(languageCode, countryCode, userCountryCode string) (countryName string, flagURL string, err error) {
	country, ok := GetCountryByCode(countryCode)
	if !ok {
		err = ErrCountryNotFound
		return
	}
	countryCode = country.CountryCode

	userCountry, ok := GetCountryByCode(userCountryCode)
	if !ok {
		err = ErrCountryNotFound
		return
	}
	userCountryCode = userCountry.CountryCode

	if countryCode == "156" {
		countryCode = "344"
	}

	flag := findCountryFlag(countryCode)
	if flag == nil {
		err = ErrCountryNotFound
		return
	}

	flagURL = flag.FlagURL
	countryName = getCountryNameByLang(flag, languageCode)
	// 如果用户在中国，且业务实在香港
	if userCountryCode == "156" && countryCode == "344" {
		countryName, ok = hongkongNameMapping[languageCode]
		if !ok {
			countryName = hongkongNameMapping["en"]
		}
	}
	return
}

// GetCountryInfo 根据语言代码和国家代码获取国家名称和旗帜URL
// 支持二字码、三字码、数字国家码、区号，语言不存在时回退到英语
func GetCountryInfo(langCode, code string) (countryName string, flagURL string, err error) {
	if len(countryFlagData) == 0 {
		err = ErrEmptyFlagData
		return
	}

	flag := findCountryFlag(code)
	if flag == nil {
		err = ErrCountryNotFound
		return
	}

	flagURL = flag.FlagURL
	countryName = getCountryNameByLang(flag, langCode)
	return
}

func findCountryFlag(countryCode string) *countryFlag {
	code := strings.ToUpper(strings.TrimSpace(countryCode))
	if code == "" {
		return nil
	}

	// 按优先级查找：数字国家码 > 二字码 > 三字码 > 区号
	if flag, ok := flagCodeMapping[code]; ok {
		return flag
	}

	if flag, ok := flagTwoCharMapping[code]; ok {
		return flag
	}

	if flag, ok := flagThreeCharMapping[code]; ok {
		return flag
	}

	// 查找区号（保持原始大小写，因为区号通常包含数字和符号）
	originalCode := strings.TrimSpace(countryCode)
	if flag, ok := flagAreaCodeMapping[originalCode]; ok {
		return flag
	}

	// 尝试添加常见的区号前缀进行查找
	if !strings.HasPrefix(originalCode, "+") && !strings.HasPrefix(originalCode, "00") {
		// 尝试添加 "+" 前缀
		if flag, ok := flagAreaCodeMapping["+"+originalCode]; ok {
			return flag
		}
		// 尝试添加 "00" 前缀
		if flag, ok := flagAreaCodeMapping["00"+originalCode]; ok {
			return flag
		}
	}

	return nil
}

func getCountryNameByLang(flag *countryFlag, langCode string) string {
	if flag == nil || flag.CountryName == nil {
		return ""
	}

	if name, ok := flag.CountryName[langCode]; ok && name != "" {
		return name
	}

	if name, ok := flag.CountryName["en"]; ok && name != "" {
		return name
	}

	return ""
}

func GetCountryName(langCode string, countryCode string) (string, error) {
	countryName, _, err := GetCountryInfo(langCode, countryCode)
	return countryName, err
}

func GetFlagURL(countryCode string) (string, error) {
	_, flagURL, err := GetCountryInfo("en", countryCode)
	return flagURL, err
}

func IsCountryCodeValid(countryCode string) bool {
	return findCountryFlag(countryCode) != nil
}

func GetAllSupportedLanguages() []string {
	if len(countryFlagData) == 0 {
		return nil
	}

	langSet := make(map[string]bool)
	for _, flag := range countryFlagData {
		for langCode := range flag.CountryName {
			langSet[langCode] = true
		}
	}

	languages := make([]string, 0, len(langSet))
	for langCode := range langSet {
		languages = append(languages, langCode)
	}

	return languages
}

func GetCountryInfoBatch(langCode string, countryCodes []string) map[string]struct {
	CountryName string
	FlagURL     string
	Error       error
} {
	result := make(map[string]struct {
		CountryName string
		FlagURL     string
		Error       error
	}, len(countryCodes))

	for _, code := range countryCodes {
		countryName, flagURL, err := GetCountryInfo(langCode, code)
		result[code] = struct {
			CountryName string
			FlagURL     string
			Error       error
		}{
			CountryName: countryName,
			FlagURL:     flagURL,
			Error:       err,
		}
	}

	return result
}
