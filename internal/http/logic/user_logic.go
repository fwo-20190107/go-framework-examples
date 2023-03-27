package logic

import (
	"context"
	"examples/internal/http/entity"
	"examples/internal/http/logic/iodata"
	"examples/internal/http/logic/repository"
	"examples/internal/http/util"
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
	user := &entity.User{
		Name:      input.Name,
		Authority: defaultAutority,
	}
	userID, err := l.userRepository.Store(ctx, user)
	if err != nil {
		return err
	}

	login := &entity.Login{
		UserID:   int(userID),
		Password: input.Password,
	}
	if err := l.loginRepository.Store(ctx, login); err != nil {
		return err
	}
	return nil
}

func (l *userLogic) ModifyAuthority(ctx context.Context, userID int, authority int8) error {
	return l.userRepository.ModifyAuthority(ctx, userID, authority)
}

func (l *userLogic) ModifyName(ctx context.Context, userID int, name string) error {
	return l.userRepository.ModifyName(ctx, userID, name)
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
