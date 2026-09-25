package repositories

import "lexi-users-srl/internal/lib/models"

type RoleRepository interface {
	Create(role models.Role) (models.Role, error)
	Get(id int64) (models.Role, error)
	GetAll() ([]models.Role, error)
	GetByName(name string) (models.Role, error)
	Update(id int64, role models.Role) (models.Role, error)
	Delete(id int64) error
}
