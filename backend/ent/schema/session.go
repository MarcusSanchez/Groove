package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Immutable(),
		field.Int("user_id"),
		field.String("token"),
		field.Time("expires_at"),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		// O2O Session <--> User(required)
		edge.From("user", User.Type).Ref("session").Field("user_id").Unique().
			// Required() to make edge required on creation;
			// i.e. Session cannot be created without its linked User
			Required(),
	}
}
