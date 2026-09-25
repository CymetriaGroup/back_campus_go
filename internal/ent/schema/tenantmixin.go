package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// TenantMixing holds the schema definition for the TenantMixing entity.
type TenantMixin struct {
	ent.Schema
}

// Fields of the TenantMixing.
func (TenantMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenant_id").NotEmpty().Immutable(),
	}
}

// Edges of the TenantMixing.
func (TenantMixin) Edges() []ent.Edge {
	return nil
}
