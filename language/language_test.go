package language

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	list := [][]string{
		{"aar", "aa", "aar", "aar", "Afar", "阿法尔语"},
		{"abk", "ab", "abk", "abk", "Abkhazian", "阿布哈兹语"},
		{"afr", "af", "afr", "afr", "Afrikaans", "南非语"},
		{"aka", "ak", "aka", "aka", "Akan", "阿肯语"},
		{"amh", "am", "amh", "amh", "Amharic", "阿姆哈拉语"},
		{"ara", "ar", "ara", "ara", "Arabic", "阿拉伯语"},
		{"arg", "an", "arg", "arg", "Aragonese", "阿拉贡语"},
		{"asm", "as", "asm", "asm", "Assamese", "阿萨姆语"},
		{"ava", "av", "ava", "ava", "Avaric", "阿瓦尔语"},
		{"ave", "ae", "ave", "ave", "Avestan", "阿维斯陀语"},
		{"aym", "ay", "aym", "aym", "Aymara", "艾马拉语"},
		{"aze", "az", "aze", "aze", "Azerbaijani", "阿塞拜疆语"},
		{"bak", "ba", "bak", "bak", "Bashkir", "巴什基尔语"},
		{"bam", "bm", "bam", "bam", "Bambara", "班巴拉语"},
		{"bel", "be", "bel", "bel", "Belarusian", "白俄罗斯语"},
		{"ben", "bn", "ben", "ben", "Bengali", "孟加拉语"},
		{"bih", "bh", "bih", "bih", "Bihari languages", "比哈尔语"},
		{"bis", "bi", "bis", "bis", "Bislama", "比斯拉马语"},
		{"bod", "bo", "bod", "tib", "Tibetan", "藏语"},
		{"bos", "bs", "bos", "bos", "Bosnian", "波斯尼亚语"},
		{"bre", "br", "bre", "bre", "Breton", "布列塔尼语"},
		{"bul", "bg", "bul", "bul", "Bulgarian", "保加利亚语"},
		{"cat", "ca", "cat", "cat", "Catalan", "加泰罗尼亚语"},
		{"ces", "cs", "ces", "cze", "Czech", "捷克语"},
		{"cha", "ch", "cha", "cha", "Chamorro", "查莫罗语"},
		{"che", "ce", "che", "che", "Chechen", "车臣语"},
		{"chv", "cv", "chv", "chv", "Chuvash", "楚瓦什语"},
		{"cor", "kw", "cor", "cor", "Cornish", "康沃尔语"},
		{"cos", "co", "cos", "cos", "Corsican", "科西嘉语"},
		{"cre", "cr", "cre", "cre", "Cree", "克里语"},
		{"cym", "cy", "cym", "wel", "Welsh", "威尔士语"},
		{"dan", "da", "dan", "dan", "Danish", "丹麦语"},
		{"deu", "de", "deu", "ger", "German", "德语"},
		{"div", "dv", "div", "div", "Dhivehi", "迪维希语"},
		{"dzo", "dz", "dzo", "dzo", "Dzongkha", "宗卡语"},
		{"ell", "el", "ell", "gre", "Greek", "希腊语"},
		{"eng", "en", "eng", "eng", "English", "英语"},
		{"est", "et", "est", "est", "Estonian", "爱沙尼亚语"},
		{"eus", "eu", "eus", "baq", "Basque", "巴士克语"},
		{"ewe", "ee", "ewe", "ewe", "Ewe", "埃维语"},
		{"fao", "fo", "fao", "fao", "Faroese", "法罗语"},
		{"fas", "fa", "fas", "per", "Persian", "法斯语"},
		{"fij", "fj", "fij", "fij", "Fijian", "斐济语"},
		{"fil", "fil", "fil", "fil", "Filipino", "菲律宾语"},
		{"fin", "fi", "fin", "fin", "Finnish", "芬兰语"},
		{"fra", "fr", "fra", "fre", "French", "法语"},
		{"fry", "fy", "fry", "fry", "Western Frisian", "西弗里西亚语"},
		{"ful", "ff", "ful", "ful", "Fulah", "弗拉尼语"},
		{"gla", "gd", "gla", "gla", "Gaelic", "盖尔语"},
		{"gle", "ga", "gle", "gle", "Irish", "爱尔兰语"},
		{"glg", "gl", "glg", "glg", "Galician", "加里西亚语"},
		{"glv", "gv", "glv", "glv", "Manx", "马恩岛语"},
		{"grn", "gn", "grn", "grn", "Guaraní", "瓜拉尼语"},
		{"guj", "gu", "guj", "guj", "Gujarati", "古吉拉特语"},
		{"hat", "ht", "hat", "hat", "Haitian", "海地语"},
		{"hau", "ha", "hau", "hau", "Hausa", "豪萨语"},
		{"heb", "he", "heb", "heb", "Hebrew", "希伯来语"},
		{"her", "hz", "her", "her", "Herero", "赫雷罗语"},
		{"hin", "hi", "hin", "hin", "Hindi", "印地语"},
		{"hmo", "ho", "hmo", "hmo", "Hiri Motu", "希利摩陀语"},
		{"hrv", "hr", "hrv", "hrv", "Croatian", "克罗地亚语"},
		{"hun", "hu", "hun", "hun", "Hungarian", "匈牙利语"},
		{"hye", "hy", "hye", "arm", "Armenian", "亚美尼亚语"},
		{"ibo", "ig", "ibo", "ibo", "Igbo", "伊博语"},
		{"iii", "ii", "iii", "iii", "Sichuan Yi", "四川彝族语"},
		{"iku", "iu", "iku", "iku", "Inuktitut", "因纽特语"},
		{"ind", "id", "ind", "ind", "Indonesian", "印度尼西亚语"},
		{"ipk", "ik", "ipk", "ipk", "Inupiaq", "伊努皮克语"},
		{"isl", "is", "isl", "ice", "Icelandic", "冰岛语"},
		{"ita", "it", "ita", "ita", "Italian", "意大利语"},
		{"jav", "jv", "jav", "jav", "Javanese", "爪哇语"},
		{"jpn", "ja", "jpn", "jpn", "Japanese", "日语"},
		{"kal", "kl", "kal", "kal", "Greenlandic", "格陵兰语"},
		{"kan", "kn", "kan", "kan", "Kannada", "卡纳拉语"},
		{"kas", "ks", "kas", "kas", "Kashmiri", "克什米尔语"},
		{"kat", "ka", "kat", "geo", "Georgian", "格鲁吉亚语"},
		{"kau", "kr", "kau", "kau", "Kanuri", "卡努里语"},
		{"kaz", "kk", "kaz", "kaz", "Kazakh", "哈萨克语"},
		{"khm", "km", "khm", "khm", "Central Khmer", "高棉语"},
		{"kik", "ki", "kik", "kik", "Kikuyu", "基库尤语"},
		{"kin", "rw", "kin", "kin", "Kinyarwanda", "卢旺达语"},
		{"kir", "ky", "kir", "kir", "Kirghiz", "吉尔吉斯语"},
		{"kom", "kv", "kom", "kom", "Komi", "科米语"},
		{"kon", "kg", "kon", "kon", "Kongo", "刚果语"},
		{"kor", "ko", "kor", "kor", "Korean", "朝鲜语"},
		{"kua", "kj", "kua", "kua", "Kwanyama", "库瓦亚马语"},
		{"kur", "ku", "kur", "kur", "Kurdish", "库尔德语"},
		{"lao", "lo", "lao", "lao", "Lao", "老挝语"},
		{"lat", "la", "lat", "lat", "Latin", "拉丁语"},
		{"lav", "lv", "lav", "lav", "Latvian", "拉脱维亚语"},
		{"lim", "li", "lim", "lim", "Limburgish", "林堡语"},
		{"lin", "ln", "lin", "lin", "Lingala", "林加拉语"},
		{"lit", "lt", "lit", "lit", "Lithuanian", "立陶宛语"},
		{"ltz", "lb", "ltz", "ltz", "Luxembourgish", "卢森堡语"},
		{"lub", "lu", "lub", "lub", "Luba-Katanga", "鲁巴加丹加语"},
		{"lug", "lg", "lug", "lug", "Ganda", "干达语"},
		{"mah", "mh", "mah", "mah", "Marshallese", "马绍尔语"},
		{"mal", "ml", "mal", "mal", "Malayalam", "马拉雅拉姆语"},
		{"mar", "mr", "mar", "mar", "Marathi", "马拉地语"},
		{"mkd", "mk", "mkd", "mac", "Macedonian", "马其顿语"},
		{"mlg", "mg", "mlg", "mlg", "Malagasy", "马达加斯加语"},
		{"mlt", "mt", "mlt", "mlt", "Maltese", "马耳他语"},
		{"mon", "mn", "mon", "mon", "Mongolian", "蒙古语"},
		{"mri", "mi", "mri", "mao", "Maori", "毛利语"},
		{"msa", "ms", "msa", "may", "Malay", "马来语"},
		{"mya", "my", "mya", "bur", "Burmese", "缅甸语"},
		{"nbl", "nr", "nbl", "nbl", "South Ndebele", "南恩德贝勒语"},
		{"nde", "nd", "nde", "nde", "North Ndebele", "北恩德贝勒语"},
		{"ndo", "ng", "ndo", "ndo", "Ndonga", "尼东阁语"},
		{"nep", "ne", "nep", "nep", "Nepali", "尼泊尔语"},
		{"nld", "nl", "nld", "dut", "Dutch, Flemish", "荷兰语"},
		{"nno", "nn", "nno", "nno", "Norwegian Nynorsk", "挪威语"},
		{"nob", "nb", "nob", "nob", "Norwegian Bokm?l", "挪威语"},
		{"nor", "no", "nor", "nor", "Norwegian", "Norwegian"},
		{"nya", "ny", "nya", "nya", "Chichewa", "奇契瓦语"},
		{"oci", "oc", "oci", "oci", "Occitan", "欧西坦语"},
		{"oji", "oj", "oji", "oji", "Ojibwa", "奥吉布瓦语"},
		{"ori", "or", "ori", "ori", "Oriya", "奥里亚语"},
		{"orm", "om", "orm", "orm", "Oromo", "奥罗莫语"},
		{"oss", "os", "oss", "oss", "Ossetian", "奥塞特语"},
		{"pli", "pi", "pli", "pli", "Pali", "巴利语"},
		{"pol", "pl", "pol", "pol", "Polish", "波兰语"},
		{"por", "pt", "por", "por", "Portuguese", "葡萄牙语"},
		{"ron", "ro", "ron", "rum", "Romanian", "罗马尼亚语"},
		{"rus", "ru", "rus", "rus", "Russian", "俄语"},
		{"sag", "sg", "sag", "sag", "Sango", "桑戈语"},
		{"san", "sa", "san", "san", "Sanskrit", "梵语"},
		{"slk", "sk", "slk", "slo", "Slovak", "斯洛伐克语"},
		{"slv", "sl", "slv", "slv", "Slovenian", "斯洛文尼亚语"},
		{"sme", "se", "sme", "sme", "Northern Sami", "北萨摩斯语"},
		{"smo", "sm", "smo", "smo", "Samoan", "萨摩亚语"},
		{"sna", "sn", "sna", "sna", "Shona", "修纳语"},
		{"som", "so", "som", "som", "Somali", "索马里语"},
		{"sot", "st", "sot", "sot", "Southern Sotho", "南索托语"},
		{"spa", "es", "spa", "spa", "Spanish", "西班牙语"},
		{"sqi", "sq", "sqi", "alb", "Albanian", "阿尔巴尼亚语"},
		{"srd", "sc", "srd", "srd", "Sardinian", "撒丁语"},
		{"srp", "sr", "srp", "srp", "Serbian", "塞尔维亚语"},
		{"swa", "sw", "swa", "swa", "Swahili", "斯瓦希里语"},
		{"swe", "sv", "swe", "swe", "Swedish", "瑞典语"},
		{"tah", "ty", "tah", "tah", "Tahitian", "塔希提语"},
		{"tam", "ta", "tam", "tam", "Tamil", "泰米尔语"},
		{"tel", "te", "tel", "tel", "Telugu", "泰卢固语"},
		{"tgk", "tg", "tgk", "tgk", "Tajik", "塔吉克语"},
		{"tgl", "tl", "tgl", "tgl", "Tagalog", "塔加拉族语"},
		{"tha", "th", "tha", "tha", "Thai", "泰语"},
		{"tir", "ti", "tir", "tir", "Tigrinya", "提格里尼亚语"},
		{"tuk", "tk", "tuk", "tuk", "Turkmen", "土库曼语"},
		{"tur", "tr", "tur", "tur", "Turkish", "土耳其语"},
		{"uig", "ug", "uig", "uig", "Uighur", "维吾尔语"},
		{"ukr", "uk", "ukr", "ukr", "Ukrainian", "乌克兰语"},
		{"urd", "ur", "urd", "urd", "Urdu", "乌尔都语"},
		{"uzb", "uz", "uzb", "uzb", "Uzbek", "乌兹别克语"},
		{"ven", "ve", "ven", "ven", "Venda", "文达语"},
		{"vie", "vi", "vie", "vie", "Vietnamese", "越南语"},
		{"yor", "yo", "yor", "yor", "Yoruba", "约鲁巴语"},
		{"zh-CN", "zh-CN", "zh-CN", "chi", "Chinese", "中文(简体)"},
		{"zh-HK", "zh-HK", "zh-HK", "chi", "HK", "香港(繁体)"},
		{"zh-TW", "zh-TW", "zh-TW", "chi", "TaiWan", "台湾(繁体)"},
		{"zha", "za", "zha", "zha", "Zhuang", "壮语"},
		{"zul", "zu", "zul", "zul", "Zulu", "祖鲁语"},
	}

	var (
		a = make(map[string]language)
		b = make(map[string]language)
		c = make(map[string]language)
		d = make(map[string]language)
	)

	for _, v := range list {
		a[strings.ToLower(v[0])] = language{
			English: v[4],
			Chinese: v[5],
		}

		b[strings.ToLower(v[1])] = language{
			English: v[4],
			Chinese: v[5],
		}

		c[strings.ToLower(v[2])] = language{
			English: v[4],
			Chinese: v[5],
		}

		d[strings.ToLower(v[3])] = language{
			English: v[4],
			Chinese: v[5],
		}
	}

	aStr, _ := json.Marshal(a)
	fmt.Println(string(aStr))

	bStr, _ := json.Marshal(b)
	fmt.Println(string(bStr))

	cStr, _ := json.Marshal(c)
	fmt.Println(string(cStr))

	dStr, _ := json.Marshal(d)
	fmt.Println(string(dStr))
}

func Test0021(t *testing.T) {
	str, _ := json.Marshal(codeMapping)
	fmt.Println(string(str))
}
