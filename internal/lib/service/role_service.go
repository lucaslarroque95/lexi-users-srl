package service

import (
	"lexi-users-srl/internal/lib/models"
	"lexi-users-srl/internal/lib/repositories"
)

type RoleService struct {
	repository repositories.RoleRepository
}

func NewRoleService(repository repositories.RoleRepository) *RoleService {
	return &RoleService{repository: repository}
}

func (rs *RoleService) CreateRole(role models.Role) (models.Role, error) {
	return rs.repository.Create(role)
}

func (rs *RoleService) GetRole(roleID int64) (models.Role, error) {
	return rs.repository.Get(roleID)
}

func (rs *RoleService) ListRoles() ([]models.Role, error) {
	return rs.repository.GetAll()
}

func (rs *RoleService) UpdateRole(roleID int64, role models.Role) (models.Role, error) {
	return rs.repository.Update(roleID, role)
}

func (rs *RoleService) DeleteRole(roleID int64) error {
	return rs.repository.Delete(roleID)
}
