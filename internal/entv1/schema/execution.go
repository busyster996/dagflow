package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Execution holds the schema definition for the Execution entity.
type Execution struct {
	ent.Schema
}

// Mixin of the Execution.
func (Execution) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		StatusMixin{},
	}
}

// Fields of the Execution.
func (Execution) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("task_id").
			Comment("任务ID"),
		field.Uint64("step_id").
			Comment("步骤ID"),
		field.Int64("seq_no").
			Nillable().
			Optional().
			Default(0).
			Comment("执行序号"),
		field.Int64("code").
			Optional().
			Nillable().
			Default(0).
			Comment("退出码"),
	}
}

// Edges of the Execution.
func (Execution) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联到任务
		edge.To("task", Task.Type).
			Unique().
			Required().
			Field("task_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),

		// 关联到步骤
		edge.To("step", Step.Type).
			Unique().
			Required().
			Field("step_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),

		// 步骤与输出的一对多关系
		edge.To("outputs", Output.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Indexes of the Execution.
func (Execution) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("state"),
		index.Fields("seq_no"),
		index.Fields("start_time"),
		index.Fields("task_id", "seq_no"),
		index.Fields("task_id", "step_id", "state"),
		index.Fields("task_id", "step_id").Unique(),
		index.Fields("task_id", "step_id", "state", "start_time", "end_time"),
	}
}
