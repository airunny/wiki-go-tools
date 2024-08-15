package icontext

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
)

const (
	accountIDKey             = "x-md-global-account-id"
	accountNameKey           = "x-md-global-account-name"
	clientIP                 = "x-md-global-client-ip"
	appVersionKey            = "x-md-global-app-version"
	platformKey              = "x-md-global-platform"
	userIdKey                = "x-md-global-user-id"
	languageCodeKey          = "x-md-global-language-code"
	countryCodeKey           = "x-md-global-country-code"
	preferredLanguageCodeKey = "x-md-global-preferred-language-code"
	appIdKey                 = "x-md-global-app-id"
	projectTypeKey           = "x-md-global-project-type"
	areaCodeKey              = "x-md-global-area-code"
	twoAreaCodeKey           = "x-md-global-two-area-code"
	requestIdKey             = "x-md-global-request-id"
	basicDataKey             = "x-md-global-basic-data"
	deviceIdKey              = "x-md-global-basic-data"
)

type Platform string

const (
	IOS     Platform = "iOS"
	Android Platform = "Android"
	PC      Platform = "PC"
	Web     Platform = "Web"
)

func withValue(ctx context.Context, key, value string) context.Context {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		md = metadata.Metadata{}
	}
	md.Set(key, value)
	ctx = metadata.NewServerContext(ctx, md)

	return metadata.AppendToClientContext(ctx, key, value)
}

func fromValue(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return "", false
	}

	out := md.Get(key)
	return out, out != ""
}

// basic data

func WithBasicData(ctx context.Context, in string) context.Context {
	return withValue(ctx, basicDataKey, in)
}

func BasicDataFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, basicDataKey)
}

// device id

func WithDeviceId(ctx context.Context, in string) context.Context {
	return withValue(ctx, deviceIdKey, in)
}

func DeviceIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, deviceIdKey)
}

// 后台用户ID

func WithAccountID(ctx context.Context, in string) context.Context {
	return withValue(ctx, accountIDKey, in)
}

func AccountIDFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, accountIDKey)
}

// 后台用户名称

func WithAccountName(ctx context.Context, in string) context.Context {
	return withValue(ctx, accountNameKey, in)
}

func AccountNameFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, accountNameKey)
}

// 获取项目类型

func WithProjectType(ctx context.Context, in string) context.Context {
	return withValue(ctx, projectTypeKey, in)
}

func ProjectTypeFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, projectTypeKey)
}

// 客户端ip

func WithClientIP(ctx context.Context, in string) context.Context {
	return withValue(ctx, clientIP, in)
}

func ClientIPFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, clientIP)
}

// app版本

func WithAppVersion(ctx context.Context, in string) context.Context {
	return withValue(ctx, appVersionKey, in)
}

func AppVersionFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, appVersionKey)
}

// app 平台

func WithAppPlatform(ctx context.Context, in Platform) context.Context {
	return withValue(ctx, platformKey, string(in))
}

func PlatformFrom(ctx context.Context) (Platform, bool) {
	plat, ok := fromValue(ctx, platformKey)
	if !ok {
		return "", false
	}
	return Platform(plat), true
}

// app 用户ID

func WithUserId(ctx context.Context, in string) context.Context {
	return withValue(ctx, userIdKey, in)
}

func UserIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, userIdKey)
}

// 语言码

func WithLanguageCode(ctx context.Context, in string) context.Context {
	return withValue(ctx, languageCodeKey, in)
}

func LanguageCodeFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, languageCodeKey)
}

// 偏好语言

func WithPreferredLanguageCode(ctx context.Context, in string) context.Context {
	return withValue(ctx, preferredLanguageCodeKey, in)
}

func PreferredLanguageCodeFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, preferredLanguageCodeKey)
}

// 城市码

func WithCountryCode(ctx context.Context, in string) context.Context {
	return withValue(ctx, countryCodeKey, in)
}

func CountryCodeFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, countryCodeKey)
}

// appid

func WithAppId(ctx context.Context, in string) context.Context {
	return withValue(ctx, appIdKey, in)
}

func AppIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, appIdKey)
}

// area code

func WithAreaCode(ctx context.Context, in string) context.Context {
	return withValue(ctx, areaCodeKey, in)
}

func AreaCodeFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, areaCodeKey)
}

// two area code

func WithTwoAreaCode(ctx context.Context, in string) context.Context {
	return withValue(ctx, twoAreaCodeKey, in)
}

func TwoAreaCodeFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, twoAreaCodeKey)
}

// request id

func WithRequestId(ctx context.Context, in string) context.Context {
	return withValue(ctx, requestIdKey, in)
}

func RequestIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, requestIdKey)
}

func LoggerValues() []interface{} {
	return []interface{}{
		"request_id", log.Valuer(func(ctx context.Context) interface{} {
			reqId, _ := RequestIdFrom(ctx)
			return reqId
		}),
		"two_area_code", log.Valuer(func(ctx context.Context) interface{} {
			twoAreaCode, _ := TwoAreaCodeFrom(ctx)
			return twoAreaCode
		}),
		"language_code", log.Valuer(func(ctx context.Context) interface{} {
			langCode, _ := LanguageCodeFrom(ctx)
			return langCode
		}),
		"country_code", log.Valuer(func(ctx context.Context) interface{} {
			countryCode, _ := CountryCodeFrom(ctx)
			return countryCode
		}),
	}
}
