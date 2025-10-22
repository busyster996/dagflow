package utility

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

/*
// 测试函数
func runPathMatchTests() {
	tests := []struct {
		paths     []string
		whitelist []string
		blacklist []string
		desc      string
	}{
		{[]string{"/app/Logs/xxx"}, []string{"/app/Logs", "/bin"}, []string{"/app", "/etc"}, "场景一"},
		{[]string{"/app/Logs"}, []string{"/app/Logs", "/bin"}, []string{"/app", "/etc"}, "场景二"},
		{[]string{"/app/config"}, []string{"/app/Logs", "/bin"}, []string{"/app", "/etc"}, "场景三"},
		{[]string{"/app"}, []string{"/app/Logs", "/bin"}, []string{"/app", "/etc"}, "场景四"},
		{[]string{"/app/Logs"}, []string{"/app", "/bin"}, []string{"/app/Logs", "/etc"}, "场景五"},
		{[]string{"/app/Logs/xxx"}, []string{"/app/Logs", "/bin"}, []string{"/app/Logs/xxx", "/etc"}, "场景六"},
		{[]string{"/app"}, []string{"/app/Logs", "/bin"}, []string{"/app/config", "/etc"}, "场景七"},
		{[]string{"/app"}, []string{"/app/data/logs", "/bin"}, []string{"/app/data/config", "/etc"}, "场景八"},
		{[]string{"/app"}, []string{}, []string{}, "场景九"},
		{[]string{"/app/Logs/xxx"}, []string{}, []string{"/app/config"}, "场景十"},
		{[]string{"/app/Logs/xxx"}, []string{"/app/Logs"}, []string{}, "场景十一"},
		{[]string{"/app"}, []string{"/app"}, []string{"/app/config"}, "场景十二"},
		{[]string{"c:\\app\\..\\"}, []string{"C:\\app"}, []string{"c:\\app\\config"}, "场景十三 (Windows路径)"},
	}

	for _, test := range tests {
		pm := NewPathMatcher(test.whitelist, test.blacklist)
		allowed, blocked := pm.Match(test.paths)

		fmt.Println("==============================")
		fmt.Printf("%s:\n", test.desc)
		fmt.Printf("  请求: %v\n", test.paths)
		fmt.Printf("  白名单: %v\n", test.whitelist)
		fmt.Printf("  黑名单: %v\n", test.blacklist)
		fmt.Println("结果:")
		fmt.Printf("  白名单: %v\n", allowed)
		fmt.Printf("  黑名单: %v\n", blocked)
		fmt.Println()
	}
}

输出结果:
==============================
场景一:
  请求: [/app/Logs/xxx]
  白名单: [/app/Logs /bin]
  黑名单: [/app /etc]
结果:
  白名单: [/app/Logs/xxx]
  黑名单: []

==============================
场景二:
  请求: [/app/Logs]
  白名单: [/app/Logs /bin]
  黑名单: [/app /etc]
结果:
  白名单: [/app/Logs]
  黑名单: []

==============================
场景三:
  请求: [/app/config]
  白名单: [/app/Logs /bin]
  黑名单: [/app /etc]
结果:
  白名单: []
  黑名单: []

==============================
场景四:
  请求: [/app]
  白名单: [/app/Logs /bin]
  黑名单: [/app /etc]
结果:
  白名单: [/app/Logs]
  黑名单: []

==============================
场景五:
  请求: [/app/Logs]
  白名单: [/app /bin]
  黑名单: [/app/Logs /etc]
结果:
  白名单: []
  黑名单: []

==============================
场景六:
  请求: [/app/Logs/xxx]
  白名单: [/app/Logs /bin]
  黑名单: [/app/Logs/xxx /etc]
结果:
  白名单: []
  黑名单: []

==============================
场景七:
  请求: [/app]
  白名单: [/app/Logs /bin]
  黑名单: [/app/config /etc]
结果:
  白名单: [/app/Logs]
  黑名单: []

==============================
场景八:
  请求: [/app]
  白名单: [/app/data/logs /bin]
  黑名单: [/app/data/config /etc]
结果:
  白名单: [/app/data/logs]
  黑名单: []

==============================
场景九:
  请求: [/app]
  白名单: []
  黑名单: []
结果:
  白名单: [/app]
  黑名单: []

==============================
场景十:
  请求: [/app/Logs/xxx]
  白名单: []
  黑名单: [/app/config]
结果:
  白名单: [/app/Logs/xxx]
  黑名单: []

==============================
场景十一:
  请求: [/app/Logs/xxx]
  白名单: [/app/Logs]
  黑名单: []
结果:
  白名单: [/app/Logs/xxx]
  黑名单: []

==============================
场景十二:
  请求: [/app]
  白名单: [/app]
  黑名单: [/app/config]
结果:
  白名单: [/app]
  黑名单: [/app/config]

==============================
场景十三 (Windows路径):
  请求: [c:\app\..\]
  白名单: [C:\app]
  黑名单: [c:\app\config]
结果:
  白名单: [c:\app]
  黑名单: [c:\app\config]
*/

type Rule struct {
	Prefix    string
	IsWinPath bool
	Allow     bool
}

type PathMatcher struct {
	hasWhitelist bool
	hasBlacklist bool
	rules        []Rule
}

// NewPathMatcher 根据目录或文件的白名单和黑名单前缀构建匹配器
func NewPathMatcher(whitelist, blacklist []string) *PathMatcher {
	pm := &PathMatcher{
		hasWhitelist: len(whitelist) > 0,
		hasBlacklist: len(blacklist) > 0,
	}
	pm.rules = pm.buildRules(whitelist, blacklist)
	return pm
}

// IsAllowed 判断给定路径是否被允许
func (pm *PathMatcher) IsAllowed(path string) bool {
	cp := pm.normalizePath(path)
	//var defaultAllow bool
	//switch {
	//case !pm.hasWhitelist && !pm.hasBlacklist:
	//	// 如果白黑名单都没有，则允许所有路径
	//	return true
	//case pm.hasBlacklist && !pm.hasWhitelist:
	//	// 如果只有黑名单，则默认允许所有路径
	//	defaultAllow = true
	//case !pm.hasBlacklist && pm.hasWhitelist:
	//	// 如果只有白名单，则默认不允许所有路径
	//	defaultAllow = false
	//}

	if !pm.hasWhitelist && !pm.hasBlacklist {
		return true
	}
	defaultAllow := pm.hasBlacklist && !pm.hasWhitelist

	for _, r := range pm.rules {
		if isDirPrefix(cp, r.Prefix) {
			return r.Allow
		}
	}
	return defaultAllow
}

// GetAllowedSubpaths 返回 basePaths 下所有允许的前缀
func (pm *PathMatcher) GetAllowedSubpaths(basePaths []string) []string {
	var allowed []string
	seen := map[string]bool{}
	for _, basePath := range basePaths {
		base := pm.normalizePath(basePath)
		// 如果允许，包括基础本身
		if pm.IsAllowed(base) && !seen[base] {
			allowed = append(allowed, pm.formatPath(isWindowsPath(basePath), base))
			seen[base] = true
		}
		// 包含更深的允许前缀
		for _, r := range pm.rules {
			if !r.Allow {
				continue
			}
			// skip if same as base or not descendant
			if r.Prefix == base || !isDirPrefix(r.Prefix, base) {
				continue
			}
			if pm.IsAllowed(r.Prefix) && !seen[r.Prefix] {
				allowed = append(allowed, pm.formatPath(r.IsWinPath, r.Prefix))
				seen[r.Prefix] = true
			}
		}
	}
	return allowed
}

// GetBlockedPrefixes 返回 basePaths 下的所有黑名单前缀
func (pm *PathMatcher) GetBlockedPrefixes(basePaths []string) []string {
	seen := map[string]bool{}
	var blocked []string
	for _, basePath := range basePaths {
		base := pm.normalizePath(basePath)
		for _, r := range pm.rules {
			// 仅考虑黑名单规则
			if r.Allow {
				continue
			}
			// 如果是基类的后代，则包含
			if isDirPrefix(r.Prefix, base) && !seen[r.Prefix] {
				blocked = append(blocked, pm.formatPath(r.IsWinPath, r.Prefix))
				seen[r.Prefix] = true
			}
		}
	}
	return blocked
}

// Match 返回允许的子路径和阻止的前缀
func (pm *PathMatcher) Match(paths []string) (allowed, blocked []string) {
	allowed = pm.GetAllowedSubpaths(paths)
	blocked = pm.GetBlockedPrefixes(allowed)
	return
}

// buildRules 按前缀长度降序对规则进行规范化和排序
func (pm *PathMatcher) buildRules(whitelist, blacklist []string) []Rule {
	rules := make([]Rule, 0, len(whitelist)+len(blacklist))
	for _, p := range whitelist {
		n := pm.normalizePath(p)
		rules = append(rules, Rule{Prefix: n, IsWinPath: isWindowsPath(p), Allow: true})
	}
	for _, p := range blacklist {
		n := pm.normalizePath(p)
		rules = append(rules, Rule{Prefix: n, IsWinPath: isWindowsPath(p), Allow: false})
	}
	sort.Slice(rules, func(i, j int) bool {
		return len(rules[i].Prefix) > len(rules[j].Prefix)
	})
	return rules
}

// isDirPrefix 检查前缀是否与路径匹配或者是否是其目录祖先
func isDirPrefix(path, prefix string) bool {
	if path == prefix {
		return true
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	return strings.HasPrefix(path, prefix)
}

// formatPath 将规范化的路径转换回特定于操作系统的分隔符
func (pm *PathMatcher) formatPath(isWin bool, path string) string {
	if isWin {
		return strings.ReplaceAll(path, "/", "\\")
	}
	return path
}

// normalizePath 清除、斜线化和小写 Windows
func (pm *PathMatcher) normalizePath(path string) string {
	isWin := isWindowsPath(path)
	p := strings.ReplaceAll(path, "\\", "/")
	p = filepath.Clean(p)
	p = filepath.ToSlash(p)
	if isWin {
		p = strings.ToLower(p)
	}
	return p
}

var winDrive = regexp.MustCompile(`^[a-zA-Z]:`)

// isWindowsPath 报告路径是否类似于 Windows 路径
func isWindowsPath(path string) bool {
	if winDrive.MatchString(path) {
		return true
	}
	if strings.HasPrefix(path, `\\`) || strings.HasPrefix(path, "//") {
		return true
	}
	return false
}
