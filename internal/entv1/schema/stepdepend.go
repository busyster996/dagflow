package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// StepDepend holds the schema definition for the StepDepend entity.
type StepDepend struct {
	ent.Schema
}

// Mixin of the StepDepend.
func (StepDepend) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the StepDepend.
func (StepDepend) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("step_id").
			Comment("步骤ID"),
		field.Uint64("dependent_id").
			Comment("依赖步骤ID"),
	}
}

// Edges of the StepDepend.
func (StepDepend) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联到当前步骤（依赖关系的主体）
		edge.To("step", Step.Type).
			Unique().
			Required().
			Field("step_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),

		// 关联到被依赖的步骤（依赖关系的目标）
		edge.To("depend_step", Step.Type).
			Unique().
			Required().
			Field("dependent_id").
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Indexes of the StepDepend.
func (StepDepend) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("step_id"),
		index.Fields("dependent_id"),
		index.Fields("step_id", "dependent_id").Unique(),
	}
}
