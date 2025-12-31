package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Step holds the schema definition for the Step entity.
type Step struct {
	ent.Schema
}

// Mixin of the Step.
func (Step) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Step.
func (Step) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("名称"),
		field.Text("desc").
			Optional().
			Nillable().
			Comment("描述"),
		field.String("kind").
			Optional().
			Nillable().
			Comment("类型"),
		field.Text("content").
			Optional().
			Nillable().
			Comment("内容"),
		field.Int64("timeout").
			Optional().
			Nillable().
			Default(86400000000000).
			Comment("超时时间(毫秒)"),
		field.String("action").
			Optional().
			Nillable().
			Comment("动作"),
		field.String("rule").
			Optional().
			Nillable().
			Comment("规则"),
		field.JSON("retry_policy", &RetryPolicy{}).
			Default(&RetryPolicy{
				Interval:    1000,
				MaxInterval: 60000,
				MaxAttempts: 15,
				Multiplier:  0.3,
			}).
			Optional().
			Comment("重试策略"),
		field.JSON("metadata", map[string]any{}).
			Optional().
			Comment("元数据"),
	}
}

// Edges of the Step.
func (Step) Edges() []ent.Edge {
	return []ent.Edge{
		// 步骤与任务的多对多关系（通过 TaskStep）
		edge.From("tasks", Task.Type).
			Ref("steps").
			Through("task_steps", TaskStep.Type),

		// 步骤与参数的多对多关系（通过 StepParam）
		edge.To("params", Param.Type).
			Through("step_params", StepParam.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),

		// 步骤依赖关系（通过 StepDepend）
		edge.To("dependencies", Step.Type).
			Through("step_depends", StepDepend.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.From("dependents", Step.Type).
			Ref("dependencies").
			Through("step_dependents", StepDepend.Type),
	}
}

// Indexes of the Step.
func (Step) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("kind"),
		index.Fields("action"),
	}
}
