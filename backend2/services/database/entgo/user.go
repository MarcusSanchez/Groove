package entgo

import (
	"context"
	"errors"
	"groove/services/database/entgo/ent"
	User "groove/services/database/entgo/ent/user"
	"groove/types"
)

var _ types.UserService = (*UserService)(nil)

type UserService struct {
	client *ent.Client
}

func NewUserService(client *ent.Client) *UserService {
	return &UserService{client: client}
}

func (u *UserService) FindByID(id int, ctx context.Context) (*types.User, error) {
	user, err := u.client.User.
		Query().
		Where(User.IDEQ(id)).
		First(ctx)
	if err != nil {
		return nil, errors.New("error finding user by id")
	}

	return &types.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}, nil
}

func (u *UserService) FindByUsername(username string, ctx context.Context) (*types.User, error) {
	user, err := u.client.User.
		Query().
		Where(User.UsernameEQ(username)).
		First(ctx)
	if err != nil {
		return nil, errors.New("error finding user by username")
	}

	return &types.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}, nil
}

func (u *UserService) FindByEmail(email string, ctx context.Context) (*types.User, error) {
	user, err := u.client.User.
		Query().
		Where(User.EmailEQ(email)).
		First(ctx)
	if err != nil {
		return nil, errors.New("error finding user by email")
	}

	return &types.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}, nil
}

func (u *UserService) ExistsByID(id int, ctx context.Context) (bool, error) {
	exists, err := u.client.User.
		Query().
		Where(User.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		return false, errors.New("error checking id")
	}
	return exists, nil
}

func (u *UserService) ExistsByUsername(username string, ctx context.Context) (bool, error) {
	exists, err := u.client.User.
		Query().
		Where(User.UsernameEQ(username)).
		Exist(ctx)
	if err != nil {
		return false, errors.New("error checking username")
	}
	return exists, nil
}

func (u *UserService) ExistsByEmail(email string, ctx context.Context) (bool, error) {
	exists, err := u.client.User.
		Query().
		Where(User.EmailEQ(email)).
		Exist(ctx)
	if err != nil {
		return false, errors.New("error checking email")
	}
	return exists, nil
}

func (u *UserService) Insert(user *types.User, ctx context.Context) error {
	_, err := u.client.User.
		Create().
		SetUsername(user.Username).
		SetPassword(user.Password).
		SetEmail(user.Email).
		Save(ctx)
	if err != nil {
		return errors.New("error inserting user")
	}
	return nil
}

func (u *UserService) Update(user *types.UserUpdate, ctx context.Context) error {
	query := u.client.User.Update()
	switch {
	case user.Username != nil:
		query.SetUsername(*user.Username)
	case user.Email != nil:
		query.SetEmail(*user.Email)
	case user.Password != nil:
		query.SetPassword(*user.Password)
	}

	_, err := query.Save(ctx)
	if err != nil {
		return errors.New("error updating user")
	}

	return nil
}

func (u *UserService) Delete(user *types.User, ctx context.Context) error {
	_, err := u.client.User.
		Delete().
		Where(User.IDEQ(user.ID)).
		Exec(ctx)
	if err != nil {
		return errors.New("error deleting user")
	}
	return nil
}
