package repository

import (
	"context"
	"examples/pkg/adapter/infra"
	"examples/pkg/adapter/repository/model"
	"examples/pkg/code"
	"examples/pkg/entity"
	"examples/pkg/errors"
	"examples/pkg/logic/repository"
	"time"
)

type loginRepository struct {
	infra.SqlHandler
}

func NewLoginRepository(handler infra.SqlHandler) repository.LoginRepository {
	return &loginRepository{handler}
}

func (r *loginRepository) GetByID(ctx context.Context, loginID string) (*entity.Login, error) {
	query := "SELECT * FROM login WHERE login_id = ?"
	row, err := r.QueryRow(ctx, query, loginID)
	if err != nil {
		return nil, err
	}

	var login model.Login
	if err := row.Scan(&login.LoginID, &login.UserID, &login.LastSignedAt, &login.Password); err != nil {
		return nil, errors.Wrap(code.CodeNotFound, err)
	}
	return &entity.Login{
		LoginID:      login.LoginID,
		UserID:       login.UserID,
		LastSignedAt: login.LastSignedAt,
		Password:     login.Password,
	}, nil
}

func (r *loginRepository) GetByUserID(ctx context.Context, userID int) (*entity.Login, error) {
	query := "SELECT * FROM login WHERE user_id = ?"
	row, err := r.QueryRow(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var login model.Login
	if err := row.Scan(&login.LoginID, &login.UserID, &login.LastSignedAt, &login.Password); err != nil {
		return nil, errors.Wrap(code.CodeNotFound, err)
	}
	return &entity.Login{
		LoginID:      login.LoginID,
		UserID:       login.UserID,
		LastSignedAt: login.LastSignedAt,
		Password:     login.Password,
	}, nil
}

func (r *loginRepository) Store(ctx context.Context, login *entity.Login) error {
	e, err := getExecutor(ctx)
	if err != nil {
		return err
	}

	query := "INSERT INTO login (`login_id`, `user_id`, `password`) VALUES (?, ?, ?)"
	if _, err := e.Execute(ctx, query, login.LoginID, login.UserID, login.Password); err != nil {
		return err
	}
	return nil
}

func (r *loginRepository) ModifyLastSigned(ctx context.Context, loginID string) error {
	e, err := getExecutor(ctx)
	if err != nil {
		return err
	}

	query := "UPDATE login SET last_signed_at = ? WHERE login_id = ?"
	if _, err := e.Execute(ctx, query, time.Now().Format(time.DateTime), loginID); err != nil {
		return err
	}
	return nil
}

var _ repository.LoginRepository = (*loginRepository)(nil)
