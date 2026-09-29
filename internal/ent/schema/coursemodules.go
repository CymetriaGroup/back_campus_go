package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"hexagonal-go-backend/internal/platform/identifier"
)

// CourseModules holds the schema definition for the CourseModules entity.
type CourseModules struct {
	ent.Schema
}

// Fields of the CourseModules.
func (CourseModules) Fields() []ent.Field {
	return []ent.Field{
		field.String("version_id").MinLen(identifier.Length).MaxLen(identifier.Length).Validate(identifier.Validate),
		field.String("title").NotEmpty(),
		field.Int("sequence_order").NonNegative(),
	}
}

// Edges of the CourseModules.
func (CourseModules) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("version", CourseVersions.Type).
			Ref("modules").
			Field("version_id").
			Unique().
			Required(),
		edge.To("lessons", Lessons.Type),
	}
}

func (CourseModules) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("version_id", "sequence_order").Unique(),
	}
}

func (CourseModules) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
