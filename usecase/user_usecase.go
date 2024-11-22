package usecase

import (
	"Quera_webinar_bot/config"
	"Quera_webinar_bot/internal/models"
	"Quera_webinar_bot/internal/repository"
	"context"
)

type UserUsecase struct {
	base       *BaseUsecase[models.User, models.User, models.User, models.User]
	cfg        *config.Config
	repository repository.UserRepository
}

func NewUserUsecase(cfg *config.Config, userRepository repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		cfg:        cfg,
		repository: userRepository,
	}
}

// Create
func (u *UserUsecase) Create(ctx context.Context, req models.User) (models.User, error) {
	return u.base.Create(ctx, req)
}
