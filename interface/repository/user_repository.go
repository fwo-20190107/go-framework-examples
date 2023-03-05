package repository

import (
	"context"
	"examples/entity"
	"examples/logic/repository"
)

type userRepository struct {
	SqlHandler
}

func NewUserRepository(handler SqlHandler) *userRepository {
	return &userRepository{handler}
}

func (r *userRepository) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	query := "SELECT * FROM user WHERE user_id = ?"
	row, err := r.QueryRow(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := row.Scan(&user.UserID, &user.Name, &user.Authority); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	query := "SELECT * FROM user"
	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var users []entity.User
	for rows.Next() {
		var user entity.User
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

func (r *userRepository) GetLoginByID(ctx context.Context, loginID string) (*entity.Login, error) {
	query := "SELECT * FROM login WHERE login_id = ?"
	row, err := r.QueryRow(ctx, query, loginID)
	if err != nil {
		return nil, err
	}

	var login entity.Login
	if err := row.Scan(&login.LoginID, &login.UserID, &login.Password); err != nil {
		return nil, err
	}
	return &login, nil
}

func (r *userRepository) ModifyAuthority(ctx context.Context, userID, authority int) error {
	query := "UPDATE user SET authority = ? WHERE user_id = ?"
	if _, err := r.Execute(ctx, query, authority, userID); err != nil {
		return err
	}
	return nil
}

var _ repository.UserRepository = (*userRepository)(nil)
