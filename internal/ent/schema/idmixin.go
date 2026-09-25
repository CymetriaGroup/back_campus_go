package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"

	"hexagonal-go-backend/internal/platform/identifier"
)

// IDMixin provides every entity with an immutable, generated ULID primary key.
type IDMixin struct {
	mixin.Schema
}

func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			StorageKey("id").
			MinLen(identifier.Length).
			MaxLen(identifier.Length).
			Validate(identifier.Validate).
			DefaultFunc(identifier.New).
			Immutable(),
	}
}
