package types

import "time"

type Session struct {
	ID         int
	UserID     int
	Token      string
	Csrf       string
	Expiration time.Time

	// relations (loaded or not)
	User *User `json:"user,omitempty"`
}

type SessionService interface {
	FindByID(id int) (*Session, error)
	FindByToken(token string) (*Session, error)
}
