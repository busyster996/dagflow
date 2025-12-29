//go:build !windows

package runner

func (c *sCmd) selfScriptSuffix() string {
	switch c.shell {
	case "ash":
		return ".ash"
	case "bash":
		return ".bash"
	case "csh":
		return ".csh"
	case "dash":
		return ".dash"
	case "ksh":
		return ".ksh"
	case "shell", "sh":
		return ".sh"
	case "tcsh":
		return ".tcsh"
	case "zsh":
		return ".zsh"
	default:
		return ".bash"
	}
}
