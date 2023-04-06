package logic

import (
	"context"
	"examples/pkg/entity"
	"examples/pkg/logic/iodata"
	"examples/pkg/logic/repository"
	"examples/pkg/util"
)

const defaultAutority = 99

type UserLogic interface {
	GetByID(ctx context.Context, userID int) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Signup(ctx context.Context, input *iodata.SignupInput) error
	ModifyAuthority(ctx context.Context, userID int, authority int8) error
	ModifyName(ctx context.Context, userID int, name string) error
	Authorization(ctx context.Context, required int8) (bool, error)
}

type userLogic struct {
	userRepository  repository.UserRepository
	loginRepository repository.LoginRepository
	transaction     repository.Transaction
}

func NewUserLogic(
	userRepository repository.UserRepository,
	loginRepository repository.LoginRepository,
	transaction repository.Transaction,
) UserLogic {
	return &userLogic{
		userRepository:  userRepository,
		loginRepository: loginRepository,
		transaction:     transaction,
	}
}

func (l *userLogic) GetByID(ctx context.Context, userID int) (*entity.User, error) {
	user, err := l.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (l *userLogic) GetAll(ctx context.Context) ([]entity.User, error) {
	users, err := l.userRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (l *userLogic) Signup(ctx context.Context, input *iodata.SignupInput) error {
	if _, err := l.transaction.Do(ctx, func(ctx context.Context) (any, error) {
		user := &entity.User{
			Name:      input.Name,
			Authority: defaultAutority,
		}
		userID, err := l.userRepository.Store(ctx, user)
		if err != nil {
			return nil, err
		}

		login := &entity.Login{
			LoginID:  input.LoginID,
			UserID:   int(userID),
			Password: input.Password,
		}
		if err := l.loginRepository.Store(ctx, login); err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

func (l *userLogic) ModifyAuthority(ctx context.Context, userID int, authority int8) error {
	if _, err := l.transaction.Do(ctx, func(ctx context.Context) (interface{}, error) {
		if err := l.userRepository.ModifyAuthority(ctx, userID, authority); err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

func (l *userLogic) ModifyName(ctx context.Context, userID int, name string) error {
	if _, err := l.transaction.Do(ctx, func(ctx context.Context) (interface{}, error) {
		if err := l.userRepository.ModifyName(ctx, userID, name); err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

func (l *userLogic) Authorization(ctx context.Context, required int8) (bool, error) {
	userID, err := util.GetUserID(ctx)
	if err != nil {
		return true, err
	}
	user, err := l.userRepository.GetByID(ctx, userID)
	if err != nil {
		return true, err
	}
	return user.Authority >= required, nil
}

var _ UserLogic = (*userLogic)(nil)
