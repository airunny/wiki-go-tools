package iheader

import (
	"strings"

	"github.com/go-kratos/kratos/v2/transport"
)

const (
	ResponseContentJsonType       = "application/json"          // json 数据
	ResponseContentTextType       = "text/plain"                // 文本数据
	TraceIdHeaderKey              = "X-Trace-Id"                // 链路追踪ID
	TokenHeaderKey                = "X-Token"                   // 用户token
	ForwardForHeaderKey           = "X-Forwarded-For"           // 客户端ip
	XRealIpHeaderKey              = "X-Real-Ip"                 // 客户端IP
	UserIdHeaderKey               = "X-User-Id"                 // 用户ID
	RequestIdKey                  = "X-Request-Id"              // request_id
	RequestIdKeyOld               = "WikidataCenter-Request-Id" // kong传递下来的request_id
	RequestIdKeyOld2              = "Wikidatacenter-Request-Id" // kong传递下来的request_idv2
	CountryCodeHeaderKeyOld       = "Countrycode"               // 国家code，后续使用X-Country-Code" 代替
	CountryCodeHeaderKey          = "X-Country-Code"            // 国家code
	LanguageCodeHeaderKeyOld      = "Languagecode"              // 语言code，使用X-Language-Code 代替
	LanguageCodeHeaderKey         = "X-Language-Code"           // 语言code
	BasicDataHeaderKey            = "Basicdata"                 // 其他信息
	PreferredLanguageHeaderKeyOld = "Preferredlanguagecode"     // 偏好语言，后续使用 X-Preferred-Language-Code 代替
	PreferredLanguageHeaderKey    = "X-Preferred-Language-Code" // 偏好语言
	SceneCodeKey                  = "SceneCode"                 // scene code
	WSCKey                        = "Route_wsc_val"
)

func GetToken(h transport.Header) string {
	return h.Get(TokenHeaderKey)
}

func GetClientIp(h transport.Header) string {
	//out := h.Get(XRealIpHeaderKey)
	//if out != "" {
	//	return out
	//}

	value := h.Get(ForwardForHeaderKey)
	splits := strings.Split(value, ",")
	if len(splits) > 0 {
		return splits[0]
	}
	return ""
}

func GetRouteWSC(h transport.Header) string {
	return h.Get(WSCKey)
}

func GetUserId(h transport.Header) string {
	return h.Get(UserIdHeaderKey)
}

func GetRequestId(h transport.Header) string {
	out := h.Get(RequestIdKey)
	if out != "" {
		return out
	}

	return h.Get(RequestIdKeyOld)
}

func GetCountryCode(h transport.Header) string {
	out := h.Get(CountryCodeHeaderKey)
	if out != "" {
		return out
	}

	return h.Get(CountryCodeHeaderKeyOld)
}

func GetLanguageCode(h transport.Header) string {
	out := h.Get(LanguageCodeHeaderKey)
	if out != "" {
		return languageMiddle(out)
	}

	return languageMiddle(h.Get(LanguageCodeHeaderKeyOld))
}

func GetPreferredLanguageCode(h transport.Header) string {
	out := h.Get(PreferredLanguageHeaderKey)
	if out != "" {
		return languageMiddle(out)
	}

	return languageMiddle(h.Get(PreferredLanguageHeaderKeyOld))
}
func GetWikiDataCenterRequestId(h transport.Header) string {
	value := h.Get(RequestIdKeyOld)
	if value != "" {
		return value
	}
	return h.Get(RequestIdKeyOld2)
}

func GetSceneCode(h transport.Header) string {
	return h.Get(SceneCodeKey)
}

func GetBasicData(h transport.Header) string {
	return h.Get(BasicDataHeaderKey)
}

func languageMiddle(in string) string {
	var (
		values    = strings.Split(in, ",")
		newValues = make([]string, 0, len(values))
	)

	for _, value := range values {
		value = strings.ToLower(value)
		if value == "zh" {
			value = "zh-hk"
		}
		newValues = append(newValues, value)
	}
	return strings.Join(newValues, ",")
}
