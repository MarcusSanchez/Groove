package types

import "context"

type User struct {
	ID       int
	Username string
	Password string
	Email    string

	// relations (loaded or not)
	//SpotifyLink *SpotifyLink
	//Session     []*Session
	//OauthState  *OAuthState
}

type UserService interface {
	FindByID(id int, ctx context.Context) (*User, error)
	FindByUsername(username string, ctx context.Context) (*User, error)
	FindByEmail(email string, ctx context.Context) (*User, error)

	ExistsByID(id int, ctx context.Context) (bool, error)
	ExistsByUsername(username string, ctx context.Context) (bool, error)
	ExistsByEmail(email string, ctx context.Context) (bool, error)

	Insert(user *User, ctx context.Context) error
	Update(user *UserUpdate, ctx context.Context) error
	Delete(user *User, ctx context.Context) error
}

type UserUpdate struct {
	Username *string
	Password *string
	Email    *string
}
