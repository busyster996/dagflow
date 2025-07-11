package types

type SFileRes struct {
	Name    string `json:"name" yaml:"name"`
	Path    string `json:"path" yaml:"path"`
	Size    int64  `json:"size" yaml:"size"`
	Mode    string `json:"mode" yaml:"mode"`
	ModTime int64  `json:"modTime" yaml:"modTime"`
	IsDir   bool   `json:"isDir" yaml:"isDir"`
}

type SFileListRes struct {
	Total int       `json:"total" yaml:"total"`
	Files SFilesRes `json:"files" yaml:"files"`
}

type SFilesRes []*SFileRes
