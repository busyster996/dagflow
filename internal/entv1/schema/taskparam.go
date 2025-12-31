package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TaskParam holds the schema definition for the TaskParam entity.
type TaskParam struct {
	ent.Schema
}

// Mixin of the TaskParam.
func (TaskParam) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the TaskParam.
func (TaskParam) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("task_id").
			Comment("任务ID"),
		field.Uint64("param_id").
			Comment("参数ID"),
	}
}

// Edges of the TaskParam.
func (TaskParam) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联到任务
		edge.To("task", Task.Type).
			Unique().
			Required().
			Field("task_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),

		// 关联到参数
		edge.To("param", Param.Type).
			Unique().
			Required().
			Field("param_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Indexes of the TaskParam.
func (TaskParam) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("task_id", "param_id").Unique(),
	}
}
