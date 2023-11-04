package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"regexp"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Immutable(),
		field.String("username").Unique().MinLen(4).MaxLen(16),
		field.String("password").Sensitive(), // sensitive won't print in logs/stack-traces.
		field.String("email").Sensitive().Unique().MinLen(4).MaxLen(320).Match(
			regexp.MustCompile("^[a-zA-Z0-9+_.-]+@[a-zA-Z0-9.-]+$"), // regex to validate email
		),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// O2O User <--> SpotifyLink(optional)
		edge.To("spotify_link", SpotifyLink.Type).Unique().
			// When User is deleted, cascade SpotifyLink referencing it.
			Annotations(entsql.OnDelete(entsql.Cascade)),
		// O2M User <--> Session
		edge.To("session", Session.Type).Annotations(entsql.OnDelete("CASCADE")),
	}
}
