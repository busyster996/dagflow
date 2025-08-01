package cuex

import (
	"fmt"
	"reflect"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"
	"cuelang.org/go/encoding/yaml"
)

func ParseYaml(template string, data map[string]any) ([]byte, error) {
	cuex := New("input", "output")
	res, err := cuex.Parse(template, data)
	if err != nil {
		return nil, err
	}
	return yaml.Encode(res)
}

func Parse(template string, data map[string]any) ([]byte, error) {
	cuex := New("input", "output")
	res, err := cuex.Parse(template, data)
	if err != nil {
		return nil, err
	}
	return res.MarshalJSON()
}

type Cuex struct {
	inputKey  string
	outputKey string
}

func New(intputKey, outputKey string) *Cuex {
	return &Cuex{
		inputKey:  intputKey,
		outputKey: outputKey,
	}
}

func (c *Cuex) Parse(template string, data map[string]any) (cue.Value, error) {
	template = fmt.Sprintf("%s: {\n%s\n}", c.outputKey, template)

	// fill the cue.Value before resolve
	intTpl, err := c.extraFillPath(c.inputKey, data)
	if err != nil {
		return cue.Value{}, err
	}
	template = strings.Join([]string{template, intTpl}, "\n")
	bi := build.NewContext().NewInstance("", nil)
	f, err := parser.ParseFile("-", template, parser.ParseComments)
	if err != nil {
		return cue.Value{}, err
	}
	if err = bi.AddSyntax(f); err != nil {
		return cue.Value{}, err
	}
	inst := cuecontext.New().BuildInstance(bi)
	if inst.Err() != nil {
		return cue.Value{}, inst.Err()
	}
	res := inst.LookupPath(cue.ParsePath(c.outputKey))
	return res.Value(), nil
}

func (c *Cuex) extraFillPath(key string, value any) (string, error) {
	val, path := cuecontext.New().CompileString(""), cue.ParsePath(key)
	if c.isNil(value) {
		val = val.FillPath(path, struct{}{})
	} else {
		val = val.FillPath(path, value)
	}
	return c.toString(val)
}

func (c *Cuex) isNil(i any) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	case reflect.String:
		return i == ""
	default:
		return false
	}
}

// toString stringify cue.Value with reference resolved
func (c *Cuex) toString(v cue.Value, opts ...cue.Option) (string, error) {
	opts = append([]cue.Option{cue.Final(), cue.Docs(true), cue.All()}, opts...)
	bs, err := format.Node(c.format(v.Syntax(opts...)), format.Simplify())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

func (c *Cuex) format(n ast.Node) ast.Node {
	switch x := n.(type) {
	case *ast.StructLit:
		var decls []ast.Decl
		for _, elt := range x.Elts {
			if _, ok := elt.(*ast.Ellipsis); ok {
				continue
			}
			decls = append(decls, elt)
		}
		return &ast.File{Decls: decls}
	case ast.Expr:
		ast.SetRelPos(x, token.NoSpace)
		return &ast.File{Decls: []ast.Decl{&ast.EmbedDecl{Expr: x}}}
	default:
		return x
	}
}
