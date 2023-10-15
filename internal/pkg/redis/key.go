package redis

import (
	"strings"

	"iam-permission/internal/pkg/constant"
)

func FormatKey(args ...string) string {
	return strings.Join(args, constant.RedisKeyDelimiter)
}
