package types

type SVersion struct {
	BuildTime string        `json:"buildTime" yaml:"buildTime"`
	Version   string        `json:"version" yaml:"version"`
	Git       *SVersionGit  `json:"git" yaml:"git"`
	Go        *SVersionGO   `json:"go" yaml:"go"`
	User      *SVersionUser `json:"user" yaml:"user"`
}

type SVersionGit struct {
	Branch string `json:"branch" yaml:"branch"`
	Commit string `json:"commit" yaml:"commit"`
	URL    string `json:"url" yaml:"url"`
}

type SVersionGO struct {
	Arch    string `json:"arch" yaml:"arch"`
	OS      string `json:"os" yaml:"os"`
	Version string `json:"version" yaml:"version"`
}

type SVersionUser struct {
	Email string `json:"email" yaml:"email"`
	Name  string `json:"name" yaml:"name"`
}
