package i18n

import (
	"context"
	"strings"
)

type Convert interface {
	Value(ctx context.Context, valueType, value, languageCode string) string
	Key(ctx context.Context, values []string, languageCode string) string
}

var (
	_              Convert = (*convert)(nil)
	defaultConvert Convert = &convert{}
)

func RegisterConvert(in Convert) {
	defaultConvert = in
}

type convert struct{}

func (c *convert) Value(ctx context.Context, valueType, value, languageCode string) string {
	return value
}

func (c *convert) Key(ctx context.Context, in []string, languageCode string) string {
	return strings.Join(in, ",")
}

type CustomTransValue struct {
	Type  string `json:"type,omitempty"`  // 值类型
	Value string `json:"value,omitempty"` // 需要翻译的标识
}

type CustomTransType string

const (
	CustomTransTypeCustom CustomTransType = "custom" // 根据KeyValues 自动生成key
	CustomTransTypeTrans  CustomTransType = "trans"  // 根据key中的字段进行翻译，并且使用KeyValues 字段进行填充
)

/**
key_type: 有两种类型：custom，trans
custom: 这种类型说明不需要翻译，key值可以直接通过数据查询获取到；比如：业务值的名称 {业务值名称}
trans: 这种类型说明key是需要进行翻译的，比如{0}监管,{0}需要替换成国家,监管需要翻译成不同的语言；

key_values：中的type
""：空字符串说明value就是对应的值，不需要再进行转换
country: 说明value中存储的是 国家码
license_plate: 说明value中存储的是 牌照业务code

示例1：{0}-{1}年
1、这里对key_values 不进行翻译，只需要将 对应的值填充进翻译语言中的占位符即可
2、将 10 替换 {0}，将 20 替换 {1} => 10-20年
{
	"key_type":"trans",
	"key": "34518",
    "key_values": [
		{
			"type":"",
			"value":"10"
		},
		{
			"type":"",
			"value":"20"
		}
	]
}

实例2：{0}监管
1、将156 转换成国家名称
2、将转换之后的国家名称 在翻译中进行占位符填充 中国监管
{
	"key_type":"trans",
	"key": "34518",
    "key_values": [
		{
			"type":"country",
			"value":"156"
		},
	]
}

示例3：{业务值名称}
1、将 1697768978 转换成 对应的业务值名称，多语言
2、将转换之后的对应语言的 业务值名称 直接赋值给key即可
{
	"key_type":"custom",
	"key": "",
    "key_values": [
		{
			"type":"license_plate",
			"value":"1697768978"
		},
	]
}

*/

type CustomTrans struct {
	KeyType    CustomTransType     `json:"key_type,omitempty"`    // key 翻译类型
	Key        string              `json:"key,omitempty"`         // key
	ChineseKey string              `json:"chinese_key,omitempty"` // 中文
	KeyValues  []*CustomTransValue `json:"key_values,omitempty"`  // key翻译中占位值列表
	DefaultKey string              `json:"default_key,omitempty"` // 如果没有找到翻译的key，则使用当前值
	Value      string              `json:"value,omitempty"`       // value ,如果是纯标签的，value一定为空
	LabelType  string              `json:"label_type,omitempty"`  // 标签类型
	Number     string              `json:"number,omitempty"`      // 标签编号
	OutKey     string              `json:"out_key,omitempty"`
	OutValue   string              `json:"out_value,omitempty"`
}

func GetLanguageWithCustomAndDefaultEnglish(ctx context.Context, values []*CustomTrans, languageCode string) {
	for _, value := range values {
		keysValues := convertValues(ctx, value.KeyValues, languageCode)
		switch value.KeyType {
		case CustomTransTypeCustom:
			value.OutKey = defaultConvert.Key(ctx, keysValues, languageCode)
		default:
			value.OutKey = GetWithTemplateDataDefault(value.Key, languageCode, value.DefaultKey, keysValues)
		}

		if value.Value != "" {
			value.OutValue = value.Value
		}
	}
}

func convertValues(ctx context.Context, values []*CustomTransValue, languageCode string) []string {
	if defaultConvert == nil {
		return nil
	}

	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, defaultConvert.Value(ctx, value.Type, value.Value, languageCode))
	}
	return out
}
