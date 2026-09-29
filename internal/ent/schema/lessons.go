package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"hexagonal-go-backend/internal/platform/identifier"
)

// Lessons holds the schema definition for the Lessons entity.
type Lessons struct {
	ent.Schema
}

// Fields of the Lessons.
func (Lessons) Fields() []ent.Field {
	return []ent.Field{
		field.String("module_id").MinLen(identifier.Length).MaxLen(identifier.Length).Validate(identifier.Validate),
		field.String("title").NotEmpty().MaxLen(500),
		field.Int("sequence_order").NonNegative(),
	}
}

// Edges of the Lessons.
func (Lessons) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("module", CourseModules.Type).
			Ref("lessons").
			Field("module_id").
			Unique().
			Required(),
		edge.To("activities", Activities.Type),
	}
}

func (Lessons) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("module_id", "sequence_order").Unique(),
	}
}

// Mixin of the Lessons.
func (Lessons) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
