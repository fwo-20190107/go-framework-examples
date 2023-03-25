package logic

import (
	"context"
	"examples/internal/http/entity"
	"examples/internal/http/logic/iodata"
	"examples/internal/http/logic/repository"
)

const defaultAutority = 99

type UserLogic interface {
	GetByID(ctx context.Context, userID int) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Signup(ctx context.Context, input *iodata.SignupInput) error
}

type userLogic struct {
	userRepository  repository.UserRepository
	loginRepository repository.LoginRepository
}

func NewUserLogic(userRepository repository.UserRepository, loginRepository repository.LoginRepository) *userLogic {
	return &userLogic{
		userRepository:  userRepository,
		loginRepository: loginRepository,
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
	user := entity.User{
		Name:      input.Name,
		Authority: defaultAutority,
	}
	userID, err := l.userRepository.Store(ctx, user)
	if err != nil {
		return err
	}

	login := entity.Login{
		UserID:   int(userID),
		Password: input.Password,
	}
	if err := l.loginRepository.Store(ctx, login); err != nil {
		return err
	}
	return nil
}

var _ UserLogic = (*userLogic)(nil)
