package models

type SPipeline struct {
	SBase
	Name string `json:"name,omitempty" gorm:"size:256;uniqueIndex;not null;comment:名称"`
	SPipelineUpdate
}

func (p *SPipeline) TableName() string {
	return "t_pipeline"
}

type SPipelineUpdate struct {
	Desc    string `json:"desc,omitempty" gorm:"comment:描述"`
	Disable *bool  `json:"disable,omitempty" gorm:"not null;default:false;comment:禁用"`
	TplType string `json:"tpl_type,omitempty" gorm:"size:256;index;not null;comment:模板类型"`
	Params  string `json:"params,omitempty" gorm:"comment:参数"`
	Content string `json:"content,omitempty" gorm:"type:text;comment:内容"`
}

type SPipelines []*SPipeline
