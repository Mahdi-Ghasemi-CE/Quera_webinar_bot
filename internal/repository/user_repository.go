package repository

import (
	"Quera_webinar_bot/internal/models"
)

type UserRepository interface {
	BaseRepository[models.User]
}
