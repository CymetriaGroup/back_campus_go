package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CourseModules holds the schema definition for the CourseModules entity.
type CourseModules struct {
	ent.Schema
}

// Fields of the CourseModules.
func (CourseModules) Fields() []ent.Field {
	return []ent.Field{
		field.String("version_id").Unique(),
		field.String("title").NotEmpty(),
		field.Int("sequence_order").Positive().NonNegative(),
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

func (CourseModules) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
