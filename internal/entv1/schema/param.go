package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Param holds the schema definition for the Param entity.
type Param struct {
	ent.Schema
}

// Mixin of the Param.
func (Param) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Param.
func (Param) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("名称"),
		field.String("value").
			Optional().
			Nillable().
			Comment("值"),
	}
}

// Edges of the Param.
func (Param) Edges() []ent.Edge {
	return []ent.Edge{
		// 参数与任务的多对多关系（通过 TaskParam）
		edge.From("tasks", Task.Type).
			Ref("params").
			Through("task_params", TaskParam.Type),

		// 参数与步骤的多对多关系（通过 StepParam）
		edge.From("steps", Step.Type).
			Ref("params").
			Through("step_params", StepParam.Type),
	}
}

// Indexes of the Param.
func (Param) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
	}
}
