package schema

import (
	"entgo.io/ent"

	"entgo.io/ent/schema/field"

	"hexagonal-go-backend/internal/platform/identifier"
)

// CourseCategories holds the schema definition for the CourseCategories entity.
type CourseCategories struct {
	ent.Schema
}

// Fields of the CourseCategories.
func (CourseCategories) Fields() []ent.Field {
	return []ent.Field{

		field.String("code").NotEmpty().MaxLen(100),
		field.String("name").NotEmpty().MaxLen(250),
		field.Text("description").Default("").Optional(),
		field.String("parent_id").MinLen(identifier.Length).MaxLen(identifier.Length).Validate(identifier.Validate).Immutable(),
	}
}

// Edges of the CourseCategories.
func (CourseCategories) Edges() []ent.Edge {
	return nil
}

func (CourseCategories) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TenantMixin{},
		TimeMixin{},
	}
}
