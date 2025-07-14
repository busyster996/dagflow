package types

type SHealthyz struct {
	Server string `json:"server" yaml:"server"`
	Client string `json:"client" yaml:"client"`
	State  string `json:"state" yaml:"state"`
}
