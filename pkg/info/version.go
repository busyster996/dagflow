package info

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	Version   = "v0.0.0"
	GitUrl    = "https://github.com/busyster996/dagflow.git"
	GitBranch = "main"
	GitCommit string
	BuildTime string
	UserName  = "busyster996"
	UserEmail = "busyster996@gmail.com"
)

func VersionInfo() string {
	var info = []string{
		fmt.Sprintf("Go Version: %s", runtime.Version()),
		fmt.Sprintf("OS: %s", runtime.GOOS),
		fmt.Sprintf("Arch: %s", runtime.GOARCH),
	}
	info = append(info, fmt.Sprintf("Version: %s", Version))
	info = append(info, fmt.Sprintf("Git URL: %s", GitUrl))
	info = append(info, fmt.Sprintf("Git Branch: %s", GitBranch))
	if GitCommit != "" {
		info = append(info, fmt.Sprintf("Git Commit: %s", GitCommit))
	}
	if BuildTime != "" {
		info = append(info, fmt.Sprintf("Build Time: %s", BuildTime))
	}
	info = append(info, fmt.Sprintf("Developer: %s", UserName))
	info = append(info, fmt.Sprintf("Mail: %s", UserEmail))
	return strings.Join(info, "\n")
}

func PrintHeadInfo() {
	fmt.Println(VersionInfo())
}
