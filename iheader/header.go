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
	RequestIdKeyOld               = "Wikidatacenter-Request-Id" // kong传递下来的request_id
	CountryCodeHeaderKeyOld       = "Countrycode"               // 国家code，后续使用X-Country-Code" 代替
	CountryCodeHeaderKey          = "X-Country-Code"            // 国家code
	LanguageCodeHeaderKeyOld      = "Languagecode"              // 语言code，使用X-Language-Code 代替
	LanguageCodeHeaderKey         = "X-Language-Code"           // 语言code
	BasicDataHeaderKey            = "Basicdata"                 // 其他信息
	PreferredLanguageHeaderKeyOld = "Preferredlanguagecode"     // 偏好语言，后续使用 X-Preferred-Language-Code 代替
	PreferredLanguageHeaderKey    = "X-Preferred-Language-Code" // 偏好语言
	VersionHeaderKey              = "X-Version"                 // 版本号
	PlatformHeaderKey             = "X-Platform"                // 平台
	AuthorizationHeaderKey        = "Authorization"             // 权限验证key
	AppIdHeaderKey                = "X-App-Id"                  // app_id
	ProjectTypeHeaderKey          = "X-Project-Type"            // 项目类型
)

func GetToken(h transport.Header) string {
	return h.Get(TokenHeaderKey)
}

func GetClientIp(h transport.Header) string {
	out := h.Get(XRealIpHeaderKey)
	if out != "" {
		return out
	}

	value := h.Get(ForwardForHeaderKey)
	splits := strings.Split(value, ",")
	if len(splits) > 0 {
		return splits[0]
	}

	return ""
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
		return out
	}

	return h.Get(LanguageCodeHeaderKeyOld)
}

func GetPreferredLanguageCode(h transport.Header) string {
	out := h.Get(PreferredLanguageHeaderKey)
	if out != "" {
		return out
	}

	return h.Get(PreferredLanguageHeaderKeyOld)
}

func GetAuthorization(h transport.Header) string {
	return h.Get(AuthorizationHeaderKey)
}

func GetBasicData(h transport.Header) string {
	return h.Get(BasicDataHeaderKey)
}

func ParseBasicData(h transport.Header) func(key string) string {
	out := h.Get(BasicDataHeaderKey)
	if out == "" {
		return func(key string) string {
			return h.Get(key)
		}
	}

	var (
		splits  = strings.Split(out, ",")
		mapping = make(map[string]string, len(splits))
	)

	if len(splits) > 0 {
		mapping[PlatformHeaderKey] = splits[0]
	}

	if len(splits) > 1 {
		mapping[AppIdHeaderKey] = splits[1]
	}

	if len(splits) > 2 {
		mapping[ProjectTypeHeaderKey] = splits[2]
	}

	if len(splits) > 3 {
		mapping[VersionHeaderKey] = splits[3]
	}

	return func(key string) string {
		value := h.Get(key)
		if value != "" {
			return value
		}

		return mapping[key]
	}
}

func GetPlatform(h transport.Header) string {
	out := h.Get(PlatformHeaderKey)
	if out != "" {
		return out
	}

	out = h.Get(BasicDataHeaderKey)
	if out == "" {
		return ""
	}

	splits := strings.Split(out, ",")
	if len(splits) > 0 {
		return splits[0]
	}

	return ""
}

func GetAppId(h transport.Header) string {
	out := h.Get(AppIdHeaderKey)
	if out != "" {
		return out
	}

	out = h.Get(BasicDataHeaderKey)
	if out == "" {
		return ""
	}

	splits := strings.Split(out, ",")
	if len(splits) > 1 {
		return splits[1]
	}

	return ""
}

func GetProjectType(h transport.Header) string {
	out := h.Get(ProjectTypeHeaderKey)
	if out != "" {
		return out
	}

	out = h.Get(BasicDataHeaderKey)
	if out == "" {
		return ""
	}

	splits := strings.Split(out, ",")
	if len(splits) > 2 {
		return splits[2]
	}

	return ""
}

func GetVersion(h transport.Header) string {
	out := h.Get(VersionHeaderKey)
	if out != "" {
		return out
	}

	out = h.Get(BasicDataHeaderKey)
	if out == "" {
		return ""
	}

	splits := strings.Split(out, ",")
	if len(splits) > 3 {
		return splits[3]
	}

	return ""
}
