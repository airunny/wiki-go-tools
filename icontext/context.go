package icontext

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
)

const (
	clientIP                   = "X-Real-Ip"                 // 客户端IP
	userIdKey                  = "x-user-id"                 // 用户ID
	languageCodeKey            = "Languagecode"              // 语言code
	countryCodeKey             = "Countrycode"               // 国家code
	preferredLanguageCodeKey   = "Preferredlanguagecode"     // 偏好语言
	requestIdKey               = "X-Request-Id"              // req id
	basicDataKey               = "Basicdata"                 // basic data
	wikiDataCenterRequestIdKey = "Wikidatacenter-Request-Id" // req id
	sceneCodeKey               = "SceneCode"                 // scene code
	wikiChannelKey             = "wikichannel"               // wiki channel
	wscKey                     = "wsc"                       // wsc
	apphpgverKey               = "apphpgver"                 // app version
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

// wiki channel

func WithWikiChannel(ctx context.Context, in string) context.Context {
	return withValue(ctx, wikiChannelKey, in)
}

func WikiChannelFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, wikiChannelKey)
}

// wsc

func WithWSC(ctx context.Context, in string) context.Context {
	return withValue(ctx, wscKey, in)
}

func WSCFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, wscKey)
}

// app hpgver

func WithAPPHPGVer(ctx context.Context, in string) context.Context {
	return withValue(ctx, apphpgverKey, in)
}

func AppHPGVerFrom(ctx context.Context) (string, bool) {
	return fromValue(ctx, apphpgverKey)
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
	}
}
