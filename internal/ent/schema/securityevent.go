package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SecurityEvent defines the schema for audit logs of security events.
type SecurityEvent struct {
	ent.Schema
}

func (SecurityEvent) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("event_type").Values(
			"SUSPICIOUS_ACCESS",
			"OUT_OF_SCOPE_ATTEMPT",
			"UNAUTHORIZED_RESOURCE_ACCESS",
			"ADMIN_ACTION",
		),
		field.String("ip_address").MaxLen(45).Optional().Nillable(),
		field.JSON("details", map[string]any{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

func (SecurityEvent) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
	}
}

func (SecurityEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("security_events").Unique(),
	}
}

func (SecurityEvent) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("event_type"),
		index.Fields("created_at"),
	}
}
