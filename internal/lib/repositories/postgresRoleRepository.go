package repositories

import (
	"lexi-users-srl/internal/lib/db"
	"lexi-users-srl/internal/lib/models"

	"gorm.io/gorm"
)

type PostgresRoleRepository struct {
	gormDB *gorm.DB
}

func NewPostgresRoleRepository(gormDB *gorm.DB) *PostgresRoleRepository {
	return &PostgresRoleRepository{gormDB: gormDB}
}

func (prr *PostgresRoleRepository) Create(role models.Role) (models.Role, error) {
	ormRole := toRoleORM(role)
	if err := prr.gormDB.Create(&ormRole).Error; err != nil {
		return models.Role{}, err
	}
	return toRoleDomain(ormRole), nil
}

func (prr *PostgresRoleRepository) Get(id int64) (models.Role, error) {
	var ormRole db.Role
	if err := prr.gormDB.First(&ormRole, id).Error; err != nil {
		return models.Role{}, err
	}
	return toRoleDomain(ormRole), nil
}

func (prr *PostgresRoleRepository) GetByName(name string) (models.Role, error) {
	var ormRole db.Role
	err := prr.gormDB.Where("name = ?", name).First(&ormRole).Error

	if err != nil {
		return models.Role{}, err
	}
	return toRoleDomain(ormRole), nil
}

func (prr *PostgresRoleRepository) GetAll() ([]models.Role, error) {
	var ormRoles []db.Role
	if err := prr.gormDB.Find(&ormRoles).Error; err != nil {
		return nil, err
	}

	roles := make([]models.Role, len(ormRoles))
	for i, ormRole := range ormRoles {
		roles[i] = toRoleDomain(ormRole)
	}
	return roles, nil
}

func (prr *PostgresRoleRepository) Update(id int64, role models.Role) (models.Role, error) {
	ormRole := toRoleORM(role)
	ormRole.ID = id

	if err := prr.gormDB.Model(&db.Role{}).Where("id = ?", id).Updates(&ormRole).Error; err != nil {
		return models.Role{}, err
	}
	return prr.Get(id)
}

func (prr *PostgresRoleRepository) Delete(id int64) error {
	return prr.gormDB.Delete(&db.Role{}, id).Error
}

func toRoleDomain(ormRole db.Role) models.Role {
	return models.Role{
		ID:   ormRole.ID,
		Name: ormRole.Name,
	}
}

func toRoleORM(role models.Role) db.Role {
	return db.Role{
		ID:   role.ID,
		Name: role.Name,
	}
}
