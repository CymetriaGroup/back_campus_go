package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"hexagonal-go-backend/internal/platform/identifier"
)

// CourseVersions holds the schema definition for the CourseVersions entity.
type CourseVersions struct {
	ent.Schema
}

// Fields of the CourseVersions.
func (CourseVersions) Fields() []ent.Field {
	return []ent.Field{
		field.String("template_id").MinLen(identifier.Length).MaxLen(identifier.Length).Validate(identifier.Validate),
		field.String("version_tag").NotEmpty().MaxLen(100),
		field.String("status").NotEmpty().MaxLen(100),
		field.Int("estimated_hours").Default(0),
	}
}

// Edges of the CourseVersions.
func (CourseVersions) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("template", CourseTemplates.Type).
			Ref("versions").
			Field("template_id").
			Unique().
			Required(),
		edge.To("syllabi", Syllabi.Type),
	}
}

func (CourseVersions) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
