package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SpotifyLink holds the schema definition for the SpotifyLink entity.
type SpotifyLink struct {
	ent.Schema
}

// Fields of the SpotifyLink.
func (SpotifyLink) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Immutable(),
		field.Int("user_id").Unique(),
		field.String("access_token"),
		field.String("refresh_token"),
	}
}

// Edges of the SpotifyLink.
func (SpotifyLink) Edges() []ent.Edge {
	return []ent.Edge{
		// O2O SpotifyLink <--> User(required)
		edge.From("user", User.Type).Ref("spotify_link").Field("user_id").Unique().
			// Required() to make edge required on creation;
			// i.e. SpotifyLink cannot be created without its linked User.
			Required(),
	}
}
