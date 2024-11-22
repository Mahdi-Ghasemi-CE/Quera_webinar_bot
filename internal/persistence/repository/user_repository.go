package repository

import (
	"Quera_webinar_bot/config"
	"Quera_webinar_bot/internal/models"
	"context"
)

const userFilterExp string = "mobile = ?"
const countFilterExp string = "count(*) > 0"

type UserRepository struct {
	*BaseRepository[models.User]
}

func NewUserRepository(cfg *config.Config) *UserRepository {
	var preloads []string
	return &UserRepository{BaseRepository: NewBaseRepository[models.User](cfg, preloads)}
}

func (r *UserRepository) CreateUser(ctx context.Context, u models.User) (models.User, error) {
	tx := r.database.WithContext(ctx).Begin()
	err := tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		return u, err
	}
	tx.Commit()
	return u, nil
}

func (r *UserRepository) FetchUserInfo(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := r.database.WithContext(ctx).
		Model(&models.User{}).
		Where(userFilterExp, username).
		Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
