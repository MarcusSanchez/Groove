package types

import "time"

type OAuthState struct {
	ID         int
	UserID     int
	State      string
	Expiration time.Time

	// relations (loaded or not)
	User *User `json:"user,omitempty"`
}

type OAuthStateService interface {
	FindByID(id int) (*OAuthState, error)
	FindByUserID(userID int) (*OAuthState, error)
}
