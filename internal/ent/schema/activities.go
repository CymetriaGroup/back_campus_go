package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"hexagonal-go-backend/internal/platform/identifier"
)

// Activities holds the schema definition for the Activities entity.
type Activities struct {
	ent.Schema
}

func (Activities) Fields() []ent.Field {
	return []ent.Field{
		field.String("lesson_id").MinLen(identifier.Length).MaxLen(identifier.Length).Validate(identifier.Validate),
		field.String("title").NotEmpty(),
		field.String("type").NotEmpty(),
		field.Bool("is_required").Default(true),
		field.Int("sequence_order").NonNegative(),
	}
}

func (Activities) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("lesson", Lessons.Type).
			Ref("activities").
			Field("lesson_id").
			Unique().
			Required(),
		edge.To("resources", CourseResources.Type),
	}
}

func (Activities) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("lesson_id", "sequence_order").Unique(),
	}
}

func (Activities) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
