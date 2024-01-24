package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

/*
 * OAuthState is a temporary store for OAuth state. State is a random 16 character string generated
 * by the backend and forwarded to Spotify. Then Spotify, after authentication, will send the state back to the
 * backend during a redirect. The backend will then verify that the state is the same as the one it sent.
 * This is to prevent CSRF attacks. A temporary table is required to store the state to keep the backend stateless.
 */

// OAuthState holds the schema definition for the OAuthState entity.
type OAuthState struct {
	ent.Schema
}

// Fields of the OAuthState.
func (OAuthState) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Immutable(),
		field.Int("user_id").Unique(),
		field.String("state").MinLen(16).MaxLen(16),
		field.Time("expiration"),
	}
}

// Edges of the OAuthState.
func (OAuthState) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("oauth_state").Field("user_id").Unique().
			// Required() to make edge required on creation;
			// i.e. state cannot be created without its linked User
			Required(),
	}
}
