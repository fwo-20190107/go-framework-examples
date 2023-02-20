package repository

import (
	"context"
	"examples/logic/repository"
	"examples/model"
)

type userRepository struct {
	SqlHandler
}

func NewUserRepository(handler SqlHandler) *userRepository {
	return &userRepository{handler}
}

func (r *userRepository) GetUserByID(ctx context.Context, userID int) (*model.User, error) {
	query := "SELECT * FROM user WHERE user_id = ?"
	row, err := r.QueryRow(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var user *model.User
	if err := row.Scan(&user.UserID, &user.Name, &user.Authority); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	query := "SELECT * FROM user"
	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.UserID,
			&user.Name,
			&user.Authority,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) GetLoginByID(ctx context.Context, loginID string) (*model.Login, error) {
	query := "SELECT * FROM login WHERE login_id = ?"
	row, err := r.QueryRow(ctx, query, loginID)
	if err != nil {
		return nil, err
	}

	var login *model.Login
	if err := row.Scan(&login.LoginID, &login.UserID, &login.Password); err != nil {
		return nil, err
	}
	return login, nil
}

func (r *userRepository) ModifyAuthority(ctx context.Context, userID, authority int) error {
	query := "UPDATE user SET authority = ? WHERE user_id = ?"
	if _, err := r.Execute(ctx, query, authority, userID); err != nil {
		return err
	}
	return nil
}

var _ repository.UserRepository = (*userRepository)(nil)
