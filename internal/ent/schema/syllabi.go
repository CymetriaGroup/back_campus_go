package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"hexagonal-go-backend/internal/platform/identifier"
)

// Syllabi holds the schema definition for the Syllabi entity.
type Syllabi struct {
	ent.Schema
}

// Fields of the Syllabi.
func (Syllabi) Fields() []ent.Field {
	return []ent.Field{
		field.String("version_id").MinLen(identifier.Length).MaxLen(identifier.Length).Validate(identifier.Validate),
		field.Text("objectives").Default("").Optional(),
		field.Text("entry_profile").Default("").Optional(),
		field.Text("exit_profile").Default("").Optional(),
		field.String("methodology").Default("").Optional(),
		field.Int("durations_hours").Default(0),
	}
}

// Edges of the Syllabi.
func (Syllabi) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("version", CourseVersions.Type).
			Ref("syllabi").
			Field("version_id").
			Unique().
			Required(),
	}
}

func (Syllabi) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
