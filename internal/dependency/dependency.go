package dependency

import (
	"Quera_webinar_bot/config"
	infraRepository "Quera_webinar_bot/internal/persistence/repository"
	contractRepository "Quera_webinar_bot/internal/repository"
)

func GetUserRepository(cfg *config.Config) contractRepository.UserRepository {
	return infraRepository.NewUserRepository(cfg)
}
