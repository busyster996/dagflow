package cuex

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/parser"
	_ "cuelang.org/go/pkg"
)

// list of std packages as of CUE
var stdPackages = map[string]string{
	"crypto":    "crypto",
	"ed25519":   "crypto/ed25519",
	"hmac":      "crypto/hmac",
	"md5":       "crypto/md5",
	"sha1":      "crypto/sha1",
	"sha256":    "crypto/sha256",
	"sha512":    "crypto/sha512",
	"base64":    "encoding/base64",
	"encoding":  "encoding",
	"csv":       "encoding/csv",
	"hex":       "encoding/hex",
	"json":      "encoding/json",
	"yaml":      "encoding/yaml",
	"html":      "html",
	"list":      "list",
	"math":      "math",
	"bits":      "math/bits",
	"net":       "net",
	"path":      "path",
	"regexp":    "regexp",
	"strconv":   "strconv",
	"strings":   "strings",
	"struct":    "struct",
	"text":      "text",
	"tabwriter": "text/tabwriter",
	"template":  "text/template",
	"time":      "time",
	"tool":      "tool",
	"cli":       "tool/cli",
	"exec":      "tool/exec",
	"file":      "tool/file",
	"http":      "tool/http",
	"os":        "tool/os",
	"uuid":      "uuid",
}

// Import reads the given cue file and updates the import statements,
// adding missing ones and removing unused ones.
// If content is nil, the file is read from disk,
// otherwise the content is used without reading the file.
// It returns the update file content.
func (c *Cuex) autoImportLib(content []byte) ([]byte, error) {
	if content == nil {
		return nil, nil
	}

	opt := []parser.Option{
		parser.ParseComments,
		parser.AllowPartial,
	}
	// ParseFile is too strict and does not allow passing a nil byte slice
	f, err := parser.ParseFile("-", content, opt...)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	unresolved := make(map[string][]string, len(f.Unresolved))
	// get a list of all unresolved identifiers
	ast.Walk(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			xIdent, ok := x.X.(*ast.Ident)
			if !ok {
				return true
			}
			xSel, ok := x.Sel.(*ast.Ident)
			if !ok {
				return true
			}

			for _, u := range f.Unresolved {
				if u.Name == xIdent.Name {
					unresolved[u.Name] = append(unresolved[u.Name], xSel.Name)
				}
			}
		}

		return true
	}, nil)

	if len(unresolved) == 0 {
		// nothing to do
		return c.insertImports(f, nil)
	}

	// resolve imports
	resolved, err := c.resolve(unresolved)
	if err != nil {
		return nil, err
	}

	// insert resolved imports
	return c.insertImports(f, resolved)
}

func (c *Cuex) resolve(unresolved map[string][]string) (map[string]string, error) {
	resolved := make(map[string]string)

	if len(unresolved) == 0 {
		return resolved, nil
	}

	// resolve using the stdlib
	c.resolveInStdlib(unresolved, resolved)

	return resolved, nil
}

func (c *Cuex) resolveInStdlib(unresolved map[string][]string, resolved map[string]string) {
	for n := range unresolved {
		if p, ok := stdPackages[n]; ok {
			resolved[n] = p
			delete(unresolved, n)
		}
	}
}

func (c *Cuex) insertImports(f *ast.File, resolved map[string]string) ([]byte, error) {
	if resolved == nil {
		resolved = map[string]string{}
	}

	var modAst ast.Node

	var err error
	if len(f.Imports) != 0 {
		// filter out unused imports
		ast.Walk(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.SelectorExpr:
				xx, ok := x.X.(*ast.Ident)
				if !ok {
					return true
				}
				for _, i := range f.Imports {
					var p string
					p, err = strconv.Unquote(i.Path.Value)
					if err != nil {
						return false
					}
					lastElem := filepath.Base(p)
					// handle package names with a different name
					// i.e. "test.com/dimensions:alt
					if idx := strings.Index(lastElem, ":"); idx != -1 {
						lastElem = lastElem[idx+1:]
					}
					if xx.Name == lastElem {
						resolved[xx.Name] = p
					}
				}
			}

			return true
		}, nil)
		if err != nil {
			return nil, err
		}
	}

	// remove all import statements
	modAst = astutil.Apply(f, func(c astutil.Cursor) bool {
		switch c.Node().(type) {
		case *ast.ImportDecl:
			c.Delete()
			return false
		}
		return true
	}, nil)

	if len(resolved) == 0 {
		return format.Node(modAst)
	}

	// sort the imports by group
	// 1. standard library

	var std []string
	for _, p := range resolved {
		// ensure the package actually belongs to the stdlib
		// and not simply ends with the same name
		if fullName, ok := stdPackages[filepath.Base(p)]; ok && fullName == p {
			std = append(std, p)
		}
	}
	sort.Strings(std)

	var idecl ast.ImportDecl

	for _, r := range std {
		idecl.Specs = append(idecl.Specs, ast.NewImport(nil, r))
	}

	var inserted bool
	// add single import statements with all resolved imports
	modAst = astutil.Apply(modAst, func(c astutil.Cursor) bool {
		switch c.Node().(type) {
		case *ast.File:
			return true
		case *ast.Package:
			if inserted {
				return false
			}
			inserted = true
			c.InsertAfter(&idecl)
			return false
		default:
			if inserted {
				return false
			}
			inserted = true
			c.InsertBefore(&idecl)
			return false
		}
	}, nil)

	return format.Node(modAst)
}
