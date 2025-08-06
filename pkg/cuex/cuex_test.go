package cuex

import "testing"

const testFirstCue = `
input: {
	aaa: "123456"
	data: json.Marshal({a: math.Sqrt(7)})
}

output: {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: 	input.name
	spec: {
		selector:
			matchLabels:
				app: input.name
		template: {
			metadata:
				labels:
					app: input.name
			spec: containers: [{
				piPlusOne: input.aaa
				image: input.image
				name:  input.name
				env:   input.env
				ports: [{
					containerPort: input.port
					protocol:      "TCP"
					name:          "default"
				}]
				if input["cpu"] != _|_ {
					resources: {
						limits:
							cpu: input.cpu
						requests:
							cpu: input.cpu
					}
				}
			}]
		}
	}
}
`

func TestParseCue(t *testing.T) {
	cue, err := Parse(testFirstCue, map[string]any{
		"name":  "test",
		"image": "nginx:latest",
		"env": []map[string]any{
			{
				"name":  "ENV",
				"value": "dev",
			},
		},
		"port": 80,
		"cpu":  "500m",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(string(cue))
}

func TestParseCueYaml(t *testing.T) {
	cue, err := ParseYaml(testFirstCue, map[string]any{
		"name":  "test",
		"image": "nginx:latest",
		"env": []map[string]any{
			{
				"name":  "ENV",
				"value": "dev",
			},
		},
		"port": 80,
		"cpu":  "500m",
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(string(cue))
}
