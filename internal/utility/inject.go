package utility

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

var customPlaceholder = regexp.MustCompile(`%([A-Za-z0-9_-]+)%`) // %VAR% style

// Inject replaces both custom %VAR%% placeholders and shell-style $VAR or ${VAR} in args.
func Inject(arg string, envs map[string]any) string {
	// First, replace all %VAR% placeholders
	arg = customPlaceholder.ReplaceAllStringFunc(arg, func(match string) string {
		// extract VAR name without the percent signs
		name := customPlaceholder.FindStringSubmatch(match)[1]
		// lookup in provided envs map, fallback to real env
		if val, ok := envs[name]; ok {
			return anyToJSONStr(val)
		}
		if lookupEnv, exists := os.LookupEnv(name); exists {
			return lookupEnv
		}
		return ""
	})
	// Then, replace shell-style placeholders using os.Expand
	arg = os.Expand(arg, func(name string) string {
		if val, ok := envs[name]; ok {
			return anyToJSONStr(val)
		}
		if lookupEnv, exists := os.LookupEnv(name); exists {
			return lookupEnv
		}
		return ""
	})
	return arg
}

func anyToJSONStr(val any) string {
	if val == nil {
		return ""
	}
	// 判断基本类型
	switch v := val.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case bool, float32, float64, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprint(v)
	default:
		// 对复杂类型（slice, map, struct, etc），尝试用 JSON 编码
		b, err := json.Marshal(v)
		if err != nil {
			// 序列化失败的话，退回到 fmt.Sprint
			return fmt.Sprint(v)
		}
		return string(b)
	}
}
