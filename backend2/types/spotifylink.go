package types

import "time"

type SpotifyLink struct {
	ID                    int
	UserID                int
	AccessToken           string
	AccessTokenExpiration time.Time
	RefreshToken          string

	// relations (loaded or not)
	User *User `json:"user,omitempty"`
}

type SpotifyLinkService interface {
	FindByID(id int) (*SpotifyLink, error)
	FindByUserID(userID int) (*SpotifyLink, error)
	FindByAccessToken(accessToken string) (*SpotifyLink, error)

	ExistsByID(id int) (bool, error)
	ExistsByUserID(userID int) (bool, error)
	ExistsByAccessToken(accessToken string) (bool, error)
}
