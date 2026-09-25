package schema

import "entgo.io/ent"

// CourseModules holds the schema definition for the CourseModules entity.
type CourseModules struct {
	ent.Schema
}

// Fields of the CourseModules.
func (CourseModules) Fields() []ent.Field {
	return nil
}

// Edges of the CourseModules.
func (CourseModules) Edges() []ent.Edge {
	return nil
}

func (CourseModules) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
	}
}
