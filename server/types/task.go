package types

import "time"

type STaskCreateRes struct {
	Name  string `json:"name" yaml:"name"`
	Count int64  `json:"count,omitempty" yaml:"count,omitempty"`
}

type STaskListDetailRes struct {
	Page  *SPageRes `json:"page" yaml:"page"`
	Tasks STasksRes `json:"tasks" yaml:"tasks"`
}

type STasksRes []*STaskRes

type STaskRes struct {
	Kind    string        `json:"kind" yaml:"kind"`
	Name    string        `json:"name" yaml:"name"`
	State   string        `json:"state" yaml:"state"`
	Count   int64         `json:"count,omitempty" yaml:"count,omitempty"`
	Desc    string        `json:"desc,omitempty" yaml:"desc,omitempty"`
	Node    string        `json:"node,omitempty" yaml:"node,omitempty"`
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Disable bool          `json:"disable,omitempty" yaml:"disable,omitempty"`
	Message string        `json:"message" yaml:"message"`
	Env     SEnvs         `json:"env,omitempty" yaml:"env,omitempty"`
	Time    *STimeRes     `json:"time,omitempty" yaml:"time,omitempty"`
}

type STaskReq struct {
	Delayed time.Time     `json:"delayed,omitempty" form:"delayed" yaml:"delayed,omitempty"`
	Kind    string        `json:"kind,omitempty" form:"kind" yaml:"kind,omitempty"`
	Name    string        `json:"name,omitempty" form:"name" yaml:"name,omitempty"`
	Desc    string        `json:"desc,omitempty" form:"desc" yaml:"desc,omitempty"`
	Node    string        `json:"node,omitempty" form:"node" yaml:"node,omitempty"`
	Disable bool          `json:"disable,omitempty" form:"disable" yaml:"disable,omitempty"`
	Timeout time.Duration `json:"timeout,omitempty" form:"timeout,omitempty" yaml:"timeout,omitempty"`
	Env     SEnvs         `json:"env,omitempty" form:"env" yaml:"env,omitempty"`
	Step    SStepsReq     `json:"step,omitempty" form:"step" yaml:"step,omitempty" binding:"required"`
}
