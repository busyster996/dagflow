package model

type Output struct {
	Base
	ExecutionID uint64 `gorm:"index:,unique,composite:key;not null;comment:执行ID"`
	Timestamp   int64  `gorm:"not null;comment:时间戳"`
	Content     string `gorm:"comment:内容"`

	// 关联关系
	Execution *Execution `gorm:"foreignKey:ExecutionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (o *Output) TableName() string {
	return "outputs"
}
