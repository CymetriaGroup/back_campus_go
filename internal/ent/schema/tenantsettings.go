package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TenantSettings holds the schema definition for the TenantSettings entity.
type TenantSettings struct {
	ent.Schema
}

// Fields of the TenantSettings.
func (TenantSettings) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("id").Immutable(),
		field.String("timezone").Default("UTC").MaxLen(100),
		field.String("language").Default("es").MaxLen(10),
		field.String("date_format").Default("YYYY-MM-DD").MaxLen(50),
		field.String("institutional_email").Default("").MaxLen(255),
		field.JSON("policies", map[string]any{}).Optional(),
		field.JSON("feature_flags", map[string]any{}).Optional(),
	}
}

// Edges of the TenantSettings.
func (TenantSettings) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tenant", Tenants.Type).Ref("settings").Unique(),
	}
}

// Mixin of the TenantSettings.
func (TenantSettings) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
