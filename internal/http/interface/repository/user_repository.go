package repository

import (
	"context"
	"examples/internal/http/entity"
	"examples/internal/http/interface/repository/model"
	"examples/internal/http/logic/repository"
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

	var user model.User
	if err := row.Scan(&user.UserID, &user.Name, &user.Authority); err != nil {
		return nil, err
	}
	return &entity.User{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	}, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	query := "SELECT * FROM user"
	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var users []entity.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.UserID,
			&user.Name,
			&user.Authority,
		); err != nil {
			return nil, err
		}
		users = append(users, entity.User{
			UserID:    user.UserID,
			Name:      user.Name,
			Authority: user.Authority,
		})
	}
	return users, nil
}

func (r *userRepository) GetLoginByID(ctx context.Context, loginID string) (*entity.Login, error) {
	query := "SELECT * FROM login WHERE login_id = ?"
	row, err := r.QueryRow(ctx, query, loginID)
	if err != nil {
		return nil, err
	}

	var login model.Login
	if err := row.Scan(&login.LoginID, &login.UserID, &login.Password); err != nil {
		return nil, err
	}
	return &entity.Login{
		LoginID:  login.LoginID,
		UserID:   login.UserID,
		Password: login.Password,
	}, nil
}

func (r *userRepository) StoreUser(ctx context.Context, user entity.User) (int64, error) {
	query := "INSERT INTO user (`user_id`, `name`, `authority`) VALUES (?, ?, ?)"
	result, err := r.Execute(ctx, query, user.UserID, user.Name, user.Authority)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *userRepository) StoreLogin(ctx context.Context, login entity.Login) error {
	query := "INSERT INTO login (`login_id`, `user_id`, `password`) VALUES (?, ?, ?)"
	if _, err := r.Execute(ctx, query, login.LoginID, login.UserID, login.Password); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) ModifyAuthority(ctx context.Context, userID, authority int) error {
	query := "UPDATE user SET authority = ? WHERE user_id = ?"
	if _, err := r.Execute(ctx, query, authority, userID); err != nil {
		return err
	}
	return nil
}

var _ repository.UserRepository = (*userRepository)(nil)
