package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Tenants holds the schema definition for the Tenants entity.
type Tenants struct {
	ent.Schema
}

// Fields of the Tenants.
func (Tenants) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").StorageKey("id").Immutable(),
		field.String("name").NotEmpty().MaxLen(255),
		field.JSON("branding", map[string]any{}).Optional(),
		field.JSON("domains", []string{}).Optional(),
		field.JSON("limits", map[string]any{}).Optional(),
		field.JSON("features", map[string]any{}).Optional(),
		field.Enum("status").Values("ACTIVE", "INACTIVE", "PENDING").Default("ACTIVE"),
	}
}

// Edges of the Tenants.
func (Tenants) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("settings", TenantSettings.Type).Unique(),
	}
}

// Mixin of the Tenants.
func (Tenants) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Indexes of the Tenants.
func (Tenants) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("status"),
	}
}
