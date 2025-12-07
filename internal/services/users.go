package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/7ngg/bread/internal/db"
)

type IUserRepository interface {
	GetUserByPhone(ctx context.Context, phone string) (db.User, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
}

type UserService struct {
	userRepository IUserRepository
	logger         *slog.Logger
}

func NewUserService(userRepository IUserRepository, logger *slog.Logger) *UserService {
	return &UserService{
		userRepository: userRepository,
		logger:         logger.With("service", "user"),
	}
}

func (us *UserService) GetUserByPhone(ctx context.Context, phone string) (db.User, error) {
	return us.userRepository.GetUserByPhone(ctx, phone)
}

func (us *UserService) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return us.userRepository.CreateUser(ctx, arg)
}

func (us *UserService) EnsureUserExists(ctx context.Context, phone, name string) (db.User, error) {
	user, err := us.GetUserByPhone(ctx, phone)

	if errors.Is(err, sql.ErrNoRows) {
		return us.CreateUser(ctx, db.CreateUserParams{
			Name:  name,
			Phone: phone,
		})
	} else if err != nil {
		return db.User{}, nil
	}

	return user, nil
}
