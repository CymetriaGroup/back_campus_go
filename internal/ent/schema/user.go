package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User defines the database schema for application users.
type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{

		field.String("name").MaxLen(100).NotEmpty(),
		field.String("email").MaxLen(255).NotEmpty(),
		field.String("password_hash").NotEmpty().Sensitive(),
		field.Enum("role").Values("admin", "user").Default("user"),
		field.Bool("active").Default(true),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("created_at"),
	}
}
