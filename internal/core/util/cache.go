package util

import (
	"fmt"
)

func GenerateCacheKey(prefix string, params any) string {
	return fmt.Sprintf("%s:%v", prefix, params)
}

func GenerateCacheKeyParams(params ...any) string {
	var str string
	last := len(params) - 1
	for i, param := range params {
		str += fmt.Sprintf("%v", param)

		if i != last {
			str += ":"
		}
	}

	return str
}
