package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Mixin of the Task.
func (Task) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
		StatusMixin{},
	}
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique().
			Comment("名称"),
		field.Text("desc").
			Optional().
			Nillable().
			Comment("描述"),
		field.String("kind").
			Optional().
			Nillable().
			Comment("类型"),
		field.String("node").
			Optional().
			Nillable().
			Comment("节点"),
		field.Int64("timeout").
			Optional().
			Nillable().
			Default(86400000000000).
			Comment("超时(毫秒)"),
		field.JSON("retry_policy", &RetryPolicy{}).
			Default(&RetryPolicy{
				Interval:    1000,
				MaxInterval: 60000,
				MaxAttempts: 15,
				Multiplier:  0.3,
			}).
			Optional().
			Comment("重试策略"),
		field.Bool("is_tpl").
			Optional().
			Nillable().
			Default(false).
			Comment("是否为模板(默认否)"),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		// 任务与步骤的多对多关系（通过 TaskStep）
		edge.To("steps", Step.Type).
			Through("task_steps", TaskStep.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		// 任务与参数的多对多关系（通过 TaskParam）
		edge.To("params", Param.Type).
			Through("task_params", TaskParam.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Indexes of the Task.
func (Task) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
		index.Fields("state"),
		index.Fields("is_tpl"),
		index.Fields("start_time"),
		index.Fields("end_time"),
		index.Fields("state", "start_time"),
	}
}
