package repositories

import "lexi-users-srl/internal/lib/models"

type UserRepository interface {
	Create(user models.User) (models.User, error)
	Get(id string) (models.User, error)
	ValidateUser(user models.User) (models.User, error)
	GetAll() ([]models.User, error)
	GetAllByRole(roleName string) ([]models.User, error)
	Update(id string, user models.User) (models.User, error)
	Delete(id string) error
}
