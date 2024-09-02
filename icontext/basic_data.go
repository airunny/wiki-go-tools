package icontext

import (
	"context"
	"strings"

	"github.com/airunny/wiki-go-tools/country"
)

func fromBasicData(ctx context.Context, index int) (string, bool) {
	basicData, ok := BasicDataFrom(ctx)
	if !ok {
		return "", false
	}

	splits := strings.Split(basicData, ",")
	if index >= len(splits) {
		return "", false
	}

	return splits[index], true
}

// device id

func DeviceIdFrom(ctx context.Context) (string, bool) {
	return fromBasicData(ctx, 5)
}

// app版本

func AppVersionFrom(ctx context.Context) (string, bool) {
	return fromBasicData(ctx, 3)
}

// app 平台

func PlatformFrom(ctx context.Context) (Platform, bool) {
	plat, ok := fromBasicData(ctx, 0)
	if !ok {
		return "", false
	}

	switch plat {
	case "0":
		return IOS, true
	case "1":
		return Android, true
	case "3":
		return PC, true
	case "999":
		return Web, true
	}
	return Platform(plat), true
}

// appid

func AppIdFrom(ctx context.Context) (string, bool) {
	return fromBasicData(ctx, 2)
}

// area code

func AreaCodeFrom(ctx context.Context) (string, bool) {
	countryCode, ok := CountryCodeFrom(ctx)
	if !ok {
		return "", false
	}

	return country.GetAreaCodeByCode(countryCode), true
}

func AllLanguageCodeFrom(ctx context.Context) []string {
	var (
		languageCode, _          = LanguageCodeFrom(ctx)
		preferredLanguageCode, _ = PreferredLanguageCodeFrom(ctx)
	)

	var (
		languages = strings.Split(preferredLanguageCode, ",")
		out       = make([]string, 0, len(languages)+1)
	)

	for _, language := range languages {
		language = strings.ToLower(strings.TrimSpace(language))
		if language == "" {
			continue
		}

		if language == strings.ToLower(languageCode) {
			continue
		}
		out = append(out, language)
	}
	out = append(out, languageCode)
	return out
}
