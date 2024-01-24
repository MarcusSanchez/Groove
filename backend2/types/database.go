package types

import "context"

type DatabaseService interface {
	Open(uri string) error
	Close() error
	Migrate(ctx context.Context) error

	Users() UserService
}
