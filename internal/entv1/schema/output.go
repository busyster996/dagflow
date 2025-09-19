package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Output holds the schema definition for the Output entity.
type Output struct {
	ent.Schema
}

// Mixin for the Output.
func (Output) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Output.
func (Output) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("exec_id").
			Comment("步骤ID"),
		field.Int64("timestamp").
			DefaultFunc(func() int64 {
				return time.Now().UnixNano()
			}).
			Optional().
			Nillable().
			Comment("时间戳"),
		field.Text("content").
			Optional().
			Nillable().
			Comment("内容"),
	}
}

// Edges of the Output.
func (Output) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联到步骤
		edge.From("execution", Execution.Type).
			Ref("outputs").
			Unique().
			Required().
			Field("exec_id"),
	}
}

// Indexes of the Output.
func (Output) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("timestamp"),
		index.Fields("exec_id", "timestamp"),
	}
}
