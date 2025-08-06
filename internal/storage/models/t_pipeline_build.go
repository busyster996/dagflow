package models

type SPipelineBuild struct {
	SBase
	PipelineName string `json:"pipeline_name,omitempty" gorm:"size:256;uniqueIndex:idx_pipeline_task_name;not null;comment:流水线名称"`
	TaskName     string `json:"task_name,omitempty" gorm:"size:256;uniqueIndex:idx_pipeline_task_name;not null;comment:任务名称"`
	Params       string `json:"params,omitempty" gorm:"comment:参数"`
}

func (p *SPipelineBuild) TableName() string {
	return "t_pipeline_build"
}

type SPipelineBuilds []*SPipelineBuildRes

type SPipelineBuildRes struct {
	SPipelineBuild
	// 链表查询使用
	STaskUpdate
}
