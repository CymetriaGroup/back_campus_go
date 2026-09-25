package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CourseTemplates holds the schema definition for the CourseTemplates entity.
type CourseTemplates struct {
	ent.Schema
}

// Fields of the CourseTemplates.
func (CourseTemplates) Fields() []ent.Field {
	return []ent.Field{

		field.String("code").Unique().NotEmpty().MaxLen(200),
		field.String("title").NotEmpty().MaxLen(200),
		field.Text("description").Default("").Optional().NotEmpty(),
	}
}

// Edges of the CourseTemplates.
func (CourseTemplates) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("versions", CourseVersions.Type),
	}
}

func (CourseTemplates) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
