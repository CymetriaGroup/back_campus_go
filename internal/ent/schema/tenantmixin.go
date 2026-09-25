package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"hexagonal-go-backend/internal/platform/identifier"
)

// TenantMixing holds the schema definition for the TenantMixing entity.
type TenantMixin struct {
	mixin.Schema
}

// Fields of the TenantMixing.
func (TenantMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenant_id").MinLen(identifier.Length).MaxLen(identifier.Length).Validate(identifier.Validate).Immutable(),
	}
}

// Edges of the TenantMixing.
func (TenantMixin) Edges() []ent.Edge {
	return nil
}
