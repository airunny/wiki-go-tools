package i18n

import (
	"archive/zip"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Language map[string]string

var (
	mux       = sync.Mutex{}
	languages = map[string]Language{
		"en": {
			"110120": "All",
			"120110": "Custody assets {0}",
			"110119": "{0}-{1} year",
			"00001":  "Regulated in {0}",
			"00002":  "监管状态{3} 监管机构{1} 国家{0} （监管号：{2}）！",
		},
		"zh-cn": {
			"110120": "全部",
			"120110": "超过 {0} 年",
			"110119": "{0}-{1} 年",
			"00001":  "{0}监管",
			"00002":  "国家{0} 监管机构{1}（监管号：{2}）监管状态{3}！",
		},
	}
	languageNameReplacer = strings.NewReplacer("TXT_", "", ".json", "")
)

func SetLanguage(in map[string]Language) {
	mux.Lock()
	defer mux.Unlock()

	languages = make(map[string]Language, len(in))
	for k, v := range in {
		languages[strings.ToLower(k)] = v
	}
}

func init() {
	languagePath := os.Getenv("LANGUAGE_PATH")
	if languagePath != "" {
		_, err := ReadFromZip(languagePath)
		if err != nil {
			panic(err)
		}
	}
}

type Entry struct {
	Language string
	Values   map[string]string
}

func ReadFromZip(zipFilePath string) ([]*Entry, error) {
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	out := make([]*Entry, 0, len(r.File))
	for _, f := range r.File {
		if !strings.HasSuffix(f.Name, ".json") {
			continue
		}

		var values map[string]string
		values, err = readFile(f)
		if err != nil {
			return nil, err
		}

		languageName := languageNameReplacer.Replace(f.Name)
		mergeToLanguage(languageName, values)

		out = append(out, &Entry{
			Language: languageName,
			Values:   values,
		})
	}
	return out, nil
}

func mergeToLanguage(languageCode string, values map[string]string) {
	languageCode = strings.ToLower(languageCode)
	langValues, ok := languages[languageCode]
	if !ok {
		languages[languageCode] = values
		return
	}

	for name, value := range langValues {
		values[name] = value
	}

	languages[languageCode] = values
}

func readFile(f *zip.File) (map[string]string, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var data map[string]string
	err = json.NewDecoder(rc).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetWithDefaultLanguage(key, languageCode, defaultLanguageCode string) string {
	out := GetLanguage(key, languageCode)
	if out != "" {
		return out
	}

	return GetLanguage(key, defaultLanguageCode)
}

func GetWithDefaultEnglish(key, languageCode string) string {
	return GetWithDefaultLanguage(key, languageCode, "en")
}

func GetLanguage(key, languageCode string) string {
	languageCode = strings.ToLower(languageCode)
	lan, ok := languages[languageCode]
	if !ok {
		return ""
	}

	return lan[key]
}

func GetWithTemplateDataDefault(key, languageCode, defaultValue string, data []string) string {
	lang := GetLanguage(key, languageCode)
	if lang == "" {
		lang = defaultValue
	}
	return genWithTemplate(lang, data)
}

func GetWithTemplateDataDefaultEnglish(key, languageCode string, data []string) string {
	lang := GetLanguage(key, languageCode)
	if lang == "" {
		lang = GetLanguage(key, "en")
	}
	return genWithTemplate(lang, data)
}

func GetWithChineseValueDefaultEnglish(value, languageCode string) string {
	values, ok := languages["zh-cn"]
	if !ok {
		return ""
	}

	var key string
	for k, v := range values {
		if strings.ToLower(v) == strings.ToLower(value) {
			key = k
			break
		}
	}

	return GetWithDefaultEnglish(key, languageCode)
}

func GetWithTemplateData(key, languageCode string, data []string) string {
	lang := GetLanguage(key, languageCode)
	return genWithTemplate(lang, data)
}

func genWithTemplate(content string, data []string) string {
	oldNews := make([]string, 0, len(data))
	for i, v := range data {
		oldNews = append(oldNews, fmt.Sprintf("{%v}", i), v)
	}

	return strings.NewReplacer(oldNews...).Replace(content)
}
