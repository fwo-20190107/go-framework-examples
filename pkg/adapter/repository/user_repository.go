package repository

import (
	"context"
	"database/sql"
	"examples/pkg/adapter/infra"
	"examples/pkg/adapter/repository/model"
	"examples/pkg/code"
	"examples/pkg/entity"
	"examples/pkg/errors"
	"examples/pkg/logic/repository"
)

type userRepository struct {
	infra.SqlHandler
}

func NewUserRepository(handler infra.SqlHandler) repository.UserRepository {
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
		c := code.CodeDatabase
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c = code.CodeNotFound
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
			return nil, errors.Wrap(code.CodeDatabase, err)
		}
		users = append(users, entity.User{
			UserID:    user.UserID,
			Name:      user.Name,
			Authority: user.Authority,
		})
	}
	return users, nil
}

func (r *userRepository) Store(ctx context.Context, user *entity.User) (int64, error) {
	e, err := getExecutor(ctx)
	if err != nil {
		return 0, err
	}

	query := "INSERT INTO user (`name`, `authority`) VALUES (?, ?)"
	result, err := e.Execute(ctx, query, user.Name, user.Authority)
	if err != nil {
		return 0, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrap(code.CodeDatabase, err)
	}
	return userID, nil
}

func (r *userRepository) ModifyAuthority(ctx context.Context, userID int, authority int8) error {
	e, err := getExecutor(ctx)
	if err != nil {
		return err
	}

	query := "UPDATE user SET authority = ? WHERE user_id = ?"
	if _, err := e.Execute(ctx, query, authority, userID); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) ModifyName(ctx context.Context, userID int, name string) error {
	e, err := getExecutor(ctx)
	if err != nil {
		return err
	}

	query := "UPDATE user SET name = ? WHERE user_id = ?"
	if _, err := e.Execute(ctx, query, name, userID); err != nil {
		return err
	}
	return nil
}

var _ repository.UserRepository = (*userRepository)(nil)
