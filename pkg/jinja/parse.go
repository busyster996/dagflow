package jinja

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

func Parse(template string, data map[string]any) (string, error) {
	gonjaConfig := &config.Config{
		BlockStartString:    "{%",
		BlockEndString:      "%}",
		VariableStartString: "{{",
		VariableEndString:   "}}",
		CommentStartString:  "{#",
		CommentEndString:    "#}",
		AutoEscape:          false,
		StrictUndefined:     false,
		TrimBlocks:          false,
		LeftStripBlocks:     false,
	}
	environment := gonja.DefaultEnvironment

	environment.Filters.Update(filters)
	environment.Tests.Update(tests)
	environment.Context.Update(globals)

	loader, err := loaders.NewFileSystemLoader("")
	if err != nil {
		return "", fmt.Errorf("failed to create a file system loader: %v", err)
	}

	sha := sha256.New()
	if _, err = sha.Write([]byte(template)); err != nil {
		return "", fmt.Errorf("failed to compute sha256 from root template")
	}
	rootID := fmt.Sprintf("root-%s", hex.EncodeToString(sha.Sum(nil)))

	shiftedLoader, err := loaders.NewShiftedLoader(rootID, bytes.NewBufferString(template), loader)
	if err != nil {
		return "", fmt.Errorf("failed to create a shifted loader: %v", err)
	}
	tpl, err := exec.NewTemplate(rootID, gonjaConfig, shiftedLoader, environment)
	if err != nil {
		return "", fmt.Errorf("failed to create a template: %v", err)
	}
	return tpl.ExecuteToString(exec.NewContext(data))
}
