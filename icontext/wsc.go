package icontext

import (
	"context"

	"github.com/buger/jsonparser"
)

func GetWSCValue(ctx context.Context, name string) (string, bool) {
	value, ok := WSCFrom(ctx)
	if !ok {
		return "", false
	}

	out, err := jsonparser.GetString([]byte(value), name)
	if err != nil {
		return "", false
	}
	return out, true
}
