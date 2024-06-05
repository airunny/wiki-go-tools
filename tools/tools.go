package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/airunny/wiki-go-tools/i18n"
	"github.com/shopspring/decimal"
)

const (
	rectCountryImage  = "https://img.souhei.com.cn/flag/%v/%v.png_wiki-template-global"
	roundCountryImage = "https://img.souhei.com.cn/flag-c/%v.png_wiki-template-global"
)

type Color string

const (
	Yellow        Color = "yellow"
	Red           Color = "red"           // 红色
	Green         Color = "green"         // 绿色
	ClosedBlue    Color = "close_blue"    // 休市蓝色
	DarkGreen     Color = "dark_green"    // 深绿色
	Gray          Color = "gray"          // 灰色
	Blue          Color = "blue"          // 蓝色
	LightGreen    Color = "light_green"   // 浅绿色
	Orange_Yellow Color = "orange_yellow" // 橘黄
	Orange_Red    Color = "orange_red"    // 橘红
	Red_360       Color = "red_360"       // 券商360红色
	Blue_360      Color = "blue_360"      // 券商360蓝色
	Black         Color = "black"         // 黑色
)

var (
	year                = time.Hour * 24 * 365
	defaultRate         = decimal.NewFromInt(1)
	toDollarRateMapping = map[string]decimal.Decimal{
		"EUR": decimal.NewFromFloat(1.10000),
		"HKD": decimal.NewFromFloat(0.13000),
		"JPY": decimal.NewFromFloat(0.00769),
		"GBP": decimal.NewFromFloat(1.20000),
		"AUD": decimal.NewFromFloat(0.64000),
		"CAD": decimal.NewFromFloat(0.72000),
		"SGD": decimal.NewFromFloat(0.73000),
		"CNY": decimal.NewFromFloat(0.1393),
	}

	colorMapping = map[Color]string{
		Red:           "#B2243C",
		DarkGreen:     "#338066",
		Green:         "#2BB351",
		Gray:          "#9EA1A2",
		Blue:          "#468084",
		LightGreen:    "#E6EFEC",
		Yellow:        "#FFFF00",
		Orange_Yellow: "#D7A638",
		Orange_Red:    "#B2243C",
		Red_360:       "#FF621F",
		Blue_360:      "#188DE3",
		Black:         "#1D2129",
		ClosedBlue:    "#004788",
	}
)

// CountryCodeToRectImage 根据城市码获取国家的矩形国旗图片地址
func CountryCodeToRectImage(in string) string {
	if in == "" {
		return ""
	}

	h := md5.New()
	h.Write([]byte(in))

	out := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf(rectCountryImage, out[8:8+16], in)
}

// TwoCharCodeToRoundImage 根据国家二字码获取国家圆形国旗图片地址
func TwoCharCodeToRoundImage(in string) string {
	h := md5.New()
	h.Write([]byte(in))

	out := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf(roundCountryImage, out[8:8+16])
}

// BrokerFoundDate 券商成立时间
func BrokerFoundDate(in time.Time) string {
	duration := time.Since(in)

	switch {
	case duration <= year:
		return "<1"
	case duration <= 2*year:
		return "1-2"
	case duration <= 5*year:
		return "2-5"
	case duration <= 10*year:
		return "5-10"
	case duration <= 15*year:
		return "10-15"
	case duration <= 20*year:
		return "15-20"
	case duration > 20*year:
		return ">20"
	default:
		return ""
	}
}

func BrokerFoundDateWithCustomTrans(in time.Time) *i18n.CustomTrans {
	duration := time.Since(in)

	switch {
	case duration <= year:
		return &i18n.CustomTrans{
			KeyType:    i18n.CustomTransTypeTrans,
			Key:        "34519",
			ChineseKey: "{0}年以内",
			KeyValues: []*i18n.CustomTransValue{
				{
					Value: "1",
				},
			},
			LabelType: "101",
			Number:    "L10101",
		}
	case duration <= 2*year:
		return &i18n.CustomTrans{
			KeyType:    i18n.CustomTransTypeTrans,
			Key:        "36046",
			ChineseKey: "{0}-{1}年",
			KeyValues: []*i18n.CustomTransValue{
				{
					Value: "1",
				},
				{
					Value: "2",
				},
			},
			LabelType: "101",
			Number:    "L10102",
		}
	case duration <= 5*year:
		return &i18n.CustomTrans{
			KeyType:    i18n.CustomTransTypeTrans,
			Key:        "36046",
			ChineseKey: "{0}-{1}年",
			KeyValues: []*i18n.CustomTransValue{
				{
					Value: "2",
				},
				{
					Value: "5",
				},
			},
			LabelType: "101",
			Number:    "L10103",
		}
	case duration <= 10*year:
		return &i18n.CustomTrans{
			KeyType:    i18n.CustomTransTypeTrans,
			Key:        "36046",
			ChineseKey: "{0}-{1}年",
			KeyValues: []*i18n.CustomTransValue{
				{
					Value: "5",
				},
				{
					Value: "10",
				},
			},
			LabelType: "101",
			Number:    "L10104",
		}
	case duration <= 15*year:
		return &i18n.CustomTrans{
			KeyType:    i18n.CustomTransTypeTrans,
			Key:        "36046",
			ChineseKey: "{0}-{1}年",
			KeyValues: []*i18n.CustomTransValue{
				{
					Value: "10",
				},
				{
					Value: "15",
				},
			},
			LabelType: "101",
			Number:    "L10105",
		}
	case duration <= 20*year:
		return &i18n.CustomTrans{
			KeyType:    i18n.CustomTransTypeTrans,
			Key:        "36046",
			ChineseKey: "{0}-{1}年",
			KeyValues: []*i18n.CustomTransValue{
				{
					Value: "15",
				},
				{
					Value: "20",
				},
			},
			LabelType: "101",
			Number:    "L10106",
		}
	case duration > 20*year:
		return &i18n.CustomTrans{
			KeyType:    i18n.CustomTransTypeTrans,
			Key:        "34517",
			ChineseKey: "{0}年以上",
			KeyValues: []*i18n.CustomTransValue{
				{
					Value: "20",
				},
			},
			LabelType: "101",
			Number:    "L10107",
		}
	}
	return &i18n.CustomTrans{}
}

// PercentageToScope 百分比转换成指定范围内的数字
func PercentageToScope(per, number int) int {
	// 先将100分成number分
	return int(math.Round(decimal.NewFromInt(int64(per)).Div(decimal.NewFromInt(100)).Mul(decimal.NewFromInt(int64(number))).InexactFloat64()))
}

// Percentage 百分比计算
func Percentage(v1, v2 decimal.Decimal) int {
	if v1.IsZero() {
		return 0
	}

	if v1.GreaterThanOrEqual(v2) {
		v1 = v2
	}

	out := int(math.Round(v1.Div(v2).Mul(decimal.NewFromInt(100)).InexactFloat64()))
	if out > 100 { // nolint:gomnd
		out = 100 // nolint:gomnd
	}
	return out
}

// PercentageToString 百分比计算，包含%
func PercentageToString(v1, v2 decimal.Decimal) string {
	out := Percentage(v1, v2)
	return fmt.Sprintf("%v%%", out)
}

// ConvertNumberToShortScale 将数字转化成短些格式
/**
 - 999,999（<百万）直接显示具体数值
- 1,000,000（1百万）可以简化为1M
- 10,000,000（1千万）可以简化为10M
- 100,000,000（1亿）可以简化为100M
- 1,000,000,000（10亿）可以简化为1B
- 10,000,000,000（1百亿）可以简化为10B
- 100,000,000,000（1千亿）可以简化为100B
- 1000,000,000,000（1万亿）可以简化为1T
*/

const (
	oneHundredMillion = 100000000 // 一亿
	tenThousand       = 10000     // 一万
)

func ConvertDecimalToShortScale(in decimal.Decimal) string {
	var (
		out    = in.String()
		suffix string
	)

	switch {
	case in.GreaterThanOrEqual(decimal.NewFromInt(10000 * oneHundredMillion)): // 1万亿
		out = in.DivRound(decimal.NewFromInt(10000*oneHundredMillion), 2).String()
		suffix = "T"
	case in.GreaterThanOrEqual(decimal.NewFromInt(1000 * oneHundredMillion)): // 1千亿
		out = in.DivRound(decimal.NewFromInt(10*oneHundredMillion), 2).String()
		suffix = "B"
	case in.GreaterThanOrEqual(decimal.NewFromInt(100 * oneHundredMillion)): // 1百亿
		out = in.DivRound(decimal.NewFromInt(10*oneHundredMillion), 2).String()
		suffix = "B"
	case in.GreaterThanOrEqual(decimal.NewFromInt(10 * oneHundredMillion)): // 10亿
		out = in.DivRound(decimal.NewFromInt(10*oneHundredMillion), 2).String()
		suffix = "B"
	case in.GreaterThanOrEqual(decimal.NewFromInt(oneHundredMillion)): // 1亿
		out = in.DivRound(decimal.NewFromInt(100*tenThousand), 2).String()
		suffix = "M"
	case in.GreaterThanOrEqual(decimal.NewFromInt(1000 * tenThousand)): // 1千万
		out = in.DivRound(decimal.NewFromInt(100*tenThousand), 2).String()
		suffix = "M"
	case in.GreaterThanOrEqual(decimal.NewFromInt(100 * tenThousand)): // 1百万
		out = in.DivRound(decimal.NewFromInt(100*tenThousand), 2).String()
		suffix = "M"
	}

	var (
		pointIndex  = strings.Index(out, ".")
		pointSuffix string
	)

	if pointIndex > 0 {
		pointSuffix = out[pointIndex:]
		out = out[:pointIndex]
	}

	if suffix == "" && len(out) > 3 {
		out = fmt.Sprintf("%v,%v", out[0:len(out)-3], out[len(out)-3:])
	}
	return fmt.Sprintf("%v%v%v", out, pointSuffix, suffix)
}

func ConvertIntToShortScale(in int64) string {
	return ConvertDecimalToShortScale(decimal.NewFromInt(in))
}

func ConvertFloatToShortScale(in float64) string {
	return ConvertDecimalToShortScale(decimal.NewFromFloat(in))
}

// ConvertToUSD 转换成美元
func ConvertToUSD(currency string, in decimal.Decimal) decimal.Decimal {
	rate, ok := toDollarRateMapping[strings.ToUpper(currency)]
	if !ok {
		rate = defaultRate
	}

	return in.Mul(rate)
}

func ConvertToUSDAndShortScale(currency string, in decimal.Decimal) string {
	return ConvertDecimalToShortScale(ConvertToUSD(currency, in))
}

func GetColor(color Color) string {
	return colorMapping[color]
}
