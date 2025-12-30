package xexec

import (
	"fmt"
	"os"
	"regexp"
)

var (
	customPlaceholder = regexp.MustCompile(`%([A-Za-z0-9_-]+)%`) // %VAR% style
)

// Inject replaces both custom %VAR% placeholders and shell-style $VAR or ${VAR} in args.
func Inject(content string, envs map[string]interface{}) string {
	// First, replace all %VAR% placeholders
	content = customPlaceholder.ReplaceAllStringFunc(content, func(match string) string {
		// extract VAR name without the percent signs
		submatches := customPlaceholder.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		name := submatches[1]
		// lookup in provided envs map, fallback to real env
		if val, ok := envs[name]; ok {
			return fmt.Sprintf("%v", val)
		}
		return lookupEnv(name)
	})
	// Then, replace shell-style placeholders using os.Expand
	content = os.Expand(content, func(name string) string {
		if val, ok := envs[name]; ok {
			return fmt.Sprintf("%v", val)
		}
		return lookupEnv(name)
	})
	return content
}

func lookupEnv(name string) string {
	if val, exists := os.LookupEnv(name); exists {
		return val
	}
	return ""
}
