package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type BaseMixin struct {
	mixin.Schema
}

func (BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Unique().
			Nillable().
			Immutable(),
		field.Time("created_at").
			Optional().
			Nillable().
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Optional().
			Nillable().
			Default(time.Now).
			UpdateDefault(time.Now),
		field.String("created_by").
			Nillable().
			Optional(),
		field.String("updated_by").
			Nillable().
			Optional(),
	}
}

type StatusMixin struct {
	mixin.Schema
}

func (StatusMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("disabled").
			Optional().
			Nillable().
			Default(false),
		field.Text("message").
			Optional().
			Nillable(),
		field.Enum("state").
			Values(
				Unknown,
				Stopped,
				Running,
				Failed,
				Pending,
				Paused,
				Skipped,
			).
			Optional().
			Nillable().
			Default(Pending).
			Comment("状态"),
		field.Enum("previous_state").
			Values(
				Unknown,
				Stopped,
				Running,
				Failed,
				Pending,
				Paused,
				Skipped,
			).
			Optional().
			Nillable().
			Default(Pending).
			Comment("状态"),
		field.Time("start_time").
			Optional().
			Nillable(),
		field.Time("end_time").
			Optional().
			Nillable(),
	}
}

type RetryPolicy struct {
	Interval    time.Duration `json:"interval" yaml:"interval" description:"间隔时间"`
	MaxInterval time.Duration `json:"maxInterval" yaml:"maxInterval" description:"最大间隔时间"`
	MaxAttempts int           `json:"maxAttempts" yaml:"maxAttempts" description:"最大尝试次数"`
	Multiplier  float64       `json:"multiplier" yaml:"multiplier" description:"乘数"`
}

const (
	Unknown = "unknown"
	Stopped = "stopped"
	Running = "running"
	Failed  = "failed"
	Pending = "pending"
	Paused  = "paused"
	Skipped = "skipped"
)
