package repository

import (
	"context"
	"database/sql"
	"examples/code"
	"examples/errors"
	"examples/internal/http/entity"
	"examples/internal/http/interface/infra"
	"examples/internal/http/interface/repository/model"
	"examples/internal/http/logic/repository"
)

type userRepository struct {
	infra.SqlHandler
}

func NewUserRepository(handler infra.SqlHandler) *userRepository {
	return &userRepository{handler}
}

func (r *userRepository) GetByID(ctx context.Context, userID int) (*entity.User, error) {
	query := "SELECT * FROM user WHERE user_id = ?"
	row, err := r.QueryRow(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := row.Scan(&user.UserID, &user.Name, &user.Authority); err != nil {
		c := code.ErrDatabase
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c = code.ErrNotFound
		}
		return nil, errors.Wrap(c, err)
	}
	return &entity.User{
		UserID:    user.UserID,
		Name:      user.Name,
		Authority: user.Authority,
	}, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]entity.User, error) {
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
			return nil, errors.Wrap(code.ErrDatabase, err)
		}
		users = append(users, entity.User{
			UserID:    user.UserID,
			Name:      user.Name,
			Authority: user.Authority,
		})
	}
	return users, nil
}

func (r *userRepository) Store(ctx context.Context, user entity.User) (int64, error) {
	query := "INSERT INTO user (`name`, `authority`) VALUES (?, ?, ?)"
	result, err := r.Execute(ctx, query, user.Name, user.Authority)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *userRepository) ModifyAuthority(ctx context.Context, userID, authority int) error {
	query := "UPDATE user SET authority = ? WHERE user_id = ?"
	if _, err := r.Execute(ctx, query, authority, userID); err != nil {
		return err
	}
	return nil
}

var _ repository.UserRepository = (*userRepository)(nil)
