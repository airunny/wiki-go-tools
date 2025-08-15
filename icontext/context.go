package icontext

import (
	"context"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
)

const (
	wikiDataCenterRequestIdKey = "WikidataCenter-Request-ID" // req id
	clientIP                   = "X-Forwarded-For"           // 客户端IP
	basicDataKey               = "BasicData"                 // basic data
	languageCodeKey            = "LanguageCode"              // 语言code
	countryCodeKey             = "CountryCode"               // 国家code
	preferredLanguageCodeKey   = "PreferredLanguageCode"     // 偏好语言
	clientPort                 = "ClientPort"                // 客户端设备端口号
	clientMac                  = "ClientMac"                 // 客户端设备Mac地址
	sceneCodeKey               = "SceneCode"                 // scene code
	requestIdKey               = "RequestId"                 // req id
	userIdKey                  = "X-User-Id"                 // 用户ID
	wscKey                     = "Route-Wsc-Val"             // wsc
	sessionAppIdKey            = "AppId"                     // 应用ID
	xPWX                       = "X-Pwa"
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
	return metadata.NewServerContext(ctx, md)
	//return metadata.AppendToClientContext(ctx, key, value)
}

func fromValue(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return "", false
	}

	out := md.Get(key)
	return out, out != ""
}

// wsc

func WithWSC(ctx context.Context, in string) context.Context {
	return withValue(ctx, wscKey, in)
}

func WSCFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, wscKey)
}

// basic data

func WithBasicData(ctx context.Context, in string) context.Context {
	return withValue(ctx, basicDataKey, in)
}

func BasicDataFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, basicDataKey)
}

// 客户端ip

func WithClientIP(ctx context.Context, in string) context.Context {
	return withValue(ctx, clientIP, in)
}

func ClientIPFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, clientIP)
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

// request id

func WithRequestId(ctx context.Context, in string) context.Context {
	return withValue(ctx, requestIdKey, in)
}

func RequestIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, requestIdKey)
}

// wiki data center request-id

func WithWikiDataCenterRequestId(ctx context.Context, in string) context.Context {
	return withValue(ctx, wikiDataCenterRequestIdKey, in)
}

func WikiDataCenterRequestIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, wikiDataCenterRequestIdKey)
}

// scene code

func WithSceneCode(ctx context.Context, in string) context.Context {
	return withValue(ctx, sceneCodeKey, in)
}

func SceneCodeFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, sceneCodeKey)
}

// clientPort

func WithClientPort(ctx context.Context, in string) context.Context {
	return withValue(ctx, clientPort, in)
}

func ClientPortFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, clientPort)
}

// clientMac

func WithClientMac(ctx context.Context, in string) context.Context {
	return withValue(ctx, clientMac, in)
}

func ClientMacFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, clientMac)
}

// appId

func WithSessionAppId(ctx context.Context, in string) context.Context {
	return withValue(ctx, sessionAppIdKey, in)
}

func SessionAppIdFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, sessionAppIdKey)
}

// x-pwa

func WithXPWA(ctx context.Context, in string) context.Context {
	return withValue(ctx, xPWX, in)
}

func XPWAFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, xPWX)
}

// context

func LoggerValues() []interface{} {
	return []interface{}{
		"user_id", log.Valuer(func(ctx context.Context) interface{} {
			userId, _ := UserIdFrom(ctx)
			return userId
		}),
		"request_id", log.Valuer(func(ctx context.Context) interface{} {
			reqId, _ := RequestIdFrom(ctx)
			return reqId
		}),
		"area_code", log.Valuer(func(ctx context.Context) interface{} {
			areaCode, _ := AreaCodeFrom(ctx)
			return areaCode
		}),
		"language_code", log.Valuer(func(ctx context.Context) interface{} {
			langCode, _ := LanguageCodeFrom(ctx)
			return langCode
		}),
		"country_code", log.Valuer(func(ctx context.Context) interface{} {
			countryCode, _ := CountryCodeFrom(ctx)
			return countryCode
		}),
		"preferred_language_code", log.Valuer(func(ctx context.Context) interface{} {
			preferredLanguageCode, _ := PreferredLanguageCodeFrom(ctx)
			return preferredLanguageCode
		}),
		"basic_data", log.Valuer(func(ctx context.Context) interface{} {
			basicData, _ := BasicDataFrom(ctx)
			return basicData
		}),
		"platform", log.Valuer(func(ctx context.Context) interface{} {
			platform, _ := PlatformFrom(ctx)
			return platform
		}),
		"scene_code", log.Valuer(func(ctx context.Context) interface{} {
			sceneCode, _ := SceneCodeFrom(ctx)
			return sceneCode
		}),
		"wikidatacenter-request-id", log.Valuer(func(ctx context.Context) interface{} {
			id, _ := WikiDataCenterRequestIdFrom(ctx)
			return id
		}),
		"wsc", log.Valuer(func(ctx context.Context) interface{} {
			value, _ := WSCFrom(ctx)
			return value
		}),
		"client_ip", log.Valuer(func(ctx context.Context) interface{} {
			clientIp, _ := ClientIPFrom(ctx)
			return clientIp
		}),
		"namespace", log.Valuer(func(ctx context.Context) interface{} {
			return os.Getenv("NAMESPACE")
		}),
	}
}
