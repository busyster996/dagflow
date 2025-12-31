package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// StepParam holds the schema definition for the StepParam entity.
type StepParam struct {
	ent.Schema
}

// Mixin for the StepParam.
func (StepParam) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the StepParam.
func (StepParam) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("step_id").
			Comment("步骤ID"),
		field.Uint64("param_id").
			Comment("参数ID"),
	}
}

// Edges of the StepParam.
func (StepParam) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联到步骤
		edge.To("step", Step.Type).
			Unique().
			Required().
			Field("step_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),

		// 关联到参数
		edge.To("param", Param.Type).
			Unique().
			Required().
			Field("param_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Indexes of the StepParam.
func (StepParam) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("step_id", "param_id").Unique(),
	}
}
