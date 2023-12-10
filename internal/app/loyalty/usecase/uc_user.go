package usecase

import (
	"context"

	"github.com/Kotletta-TT/bonus-service/config"
	"github.com/Kotletta-TT/bonus-service/internal/app/loyalty/entity"
	"github.com/Kotletta-TT/bonus-service/internal/app/loyalty/utils"
)

// Contracts
type UserManagmentUseCase interface {
	RegisterUser(ctx context.Context, login, password string) (string, error)
	LoginUser(ctx context.Context, login, password string) (string, error)
}

type UserRepo interface {
	AddUser(ctx context.Context, user *entity.User) error
	GetUserByLogin(ctx context.Context, login string) (*entity.User, error)
}

// Realization contract
type UserUseCase struct {
	repo   UserRepo
	config *config.Config
}

func New(repo UserRepo, config *config.Config) *UserUseCase {
	return &UserUseCase{repo: repo, config: config}
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, login, password string) (string, error) {
	newUser := entity.User{Login: login, Password: password}
	if err := uc.repo.AddUser(ctx, &newUser); err != nil {
		return "", err
	}
	return utils.GenerateToken(login, uc.config.SecretKey)
}

func (uc *UserUseCase) LoginUser(ctx context.Context, login, password string) (string, error) {
	user, err := uc.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	if err := utils.VerifyPassword(password, user.Password); err != nil {
		return "", err
	}
	return utils.GenerateToken(login, uc.config.SecretKey)
}
