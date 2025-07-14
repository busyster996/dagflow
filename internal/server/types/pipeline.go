package types

type SPipelineListRes struct {
	Page      *SPageRes     `json:"page" yaml:"page"`
	Pipelines SPipelinesRes `json:"pipelines" yaml:"pipelines"`
}

type SPipelinesRes []*SPipelineRes

type SPipelineRes struct {
	Name    string `json:"name" yaml:"name"`
	Desc    string `json:"desc,omitempty" yaml:"desc,omitempty"`
	Disable bool   `json:"disable,omitempty" yaml:"disable,omitempty"`
	TplType string `json:"tplType" yaml:"tplType"`
	Params  string `json:"params,omitempty" yaml:"params,omitempty"`
	Content string `json:"content,omitempty" yaml:"content,omitempty"`
}

type SPipelineCreateReq struct {
	Name    string `json:"name" yaml:"name" binding:"required"`
	Desc    string `json:"desc,omitempty" yaml:"desc,omitempty"`
	Disable *bool  `json:"disable" yaml:"disable"`
	TplType string `json:"tplType" yaml:"tplType" binding:"required" example:"jinja2"`
	Params  string `json:"params,omitempty" yaml:"params,omitempty"`
	Content string `json:"content" yaml:"content" binding:"required"`
}

type SPipelineUpdateReq struct {
	Desc    string `json:"desc,omitempty" yaml:"desc,omitempty"`
	Disable *bool  `json:"disable" yaml:"disable"`
	TplType string `json:"tplType" yaml:"tplType" binding:"required" example:"jinja2"`
	Params  string `json:"params,omitempty" yaml:"params,omitempty"`
	Content string `json:"content" yaml:"content" binding:"required"`
}

type SPipelineBuildRes struct {
	PipelineName string    `json:"pipelineName" yaml:"pipelineName"`
	TaskName     string    `json:"taskName" yaml:"taskName"`
	Params       string    `json:"params,omitempty" yaml:"params,omitempty"`
	State        string    `json:"state" yaml:"state"`
	Message      string    `json:"message" yaml:"message"`
	Time         *STimeRes `json:"time" yaml:"time"`
}

type SPipelineBuildReq struct {
	Params map[string]any `json:"params,omitempty" yaml:"params,omitempty"`
}

type SPipelineBuildListRes struct {
	Page  *SPageRes          `json:"page" yaml:"page"`
	Tasks SPipelineBuildsRes `json:"tasks" yaml:"tasks"`
}

type SPipelineBuildsRes []*SPipelineBuildRes
