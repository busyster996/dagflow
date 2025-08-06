package cuex

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/encoding/yaml"
)

func ParseYaml(template string, data map[string]any) ([]byte, error) {
	cuex := New("input", "output")
	res, err := cuex.Parse([]byte(template), data)
	if err != nil {
		return nil, err
	}
	return yaml.Encode(res)
}

func Parse(template string, data map[string]any) ([]byte, error) {
	cuex := New("input", "output")
	res, err := cuex.Parse([]byte(template), data)
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

func (c *Cuex) Parse(template []byte, data map[string]any) (cue.Value, error) {
	content, err := c.autoImportLib(template)
	if err != nil {
		return cue.Value{}, err
	}
	f, err := parser.ParseFile(
		"-",
		content,
		parser.ParseComments,
		parser.ParseFuncs,
		parser.Trace,
		parser.DeclarationErrors,
		parser.AllErrors,
		parser.AllowPartial,
	)
	if err != nil {
		return cue.Value{}, err
	}

	inst := cuecontext.New().BuildFile(f).
		FillPath(cue.ParsePath(c.inputKey), data)
	if inst.Err() != nil {
		return cue.Value{}, inst.Err()
	}

	output := inst.LookupPath(cue.ParsePath(c.outputKey))
	if output.Err() != nil {
		return cue.Value{}, output.Err()
	}

	return output.Value(), output.Err()
}
