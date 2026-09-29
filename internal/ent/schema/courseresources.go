package schema

import (
	"encoding/json"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"hexagonal-go-backend/internal/platform/identifier"
)

// CourseResources holds the schema definition for the CourseResources entity.
type CourseResources struct {
	ent.Schema
}

func (CourseResources) Fields() []ent.Field {
	return []ent.Field{
		field.String("activity_id").MinLen(identifier.Length).MaxLen(identifier.Length).Validate(identifier.Validate),
		field.String("type").NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("url_storage_key").NotEmpty(),
		field.String("mime_type").NotEmpty(),
		field.Int64("size_bytes").NonNegative(),
		field.Int("position").NonNegative(),
		field.JSON("metadata", json.RawMessage{}).Optional(),
	}
}

func (CourseResources) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("activity", Activities.Type).
			Ref("resources").
			Field("activity_id").
			Unique().
			Required(),
	}
}

func (CourseResources) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("activity_id", "position").Unique(),
	}
}

func (CourseResources) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
