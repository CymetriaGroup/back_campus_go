package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TimeMixing holds the schema definition for the TimeMixing entity.
type TimeMixin struct {
	mixin.Schema
}

// Fields of the TimeMixing.
func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now()).UpdateDefault(time.Now),
	}
}

// Edges of the TimeMixing.
func (TimeMixin) Edges() []ent.Edge {
	return nil
}
