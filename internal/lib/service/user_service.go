package service

import (
	"errors"
	"lexi-users-srl/internal/lib/models"
	"lexi-users-srl/internal/lib/repositories"
	"lexi-users-srl/internal/lib/utils"
)

type UserService struct {
	repository     repositories.UserRepository
	roleRepository repositories.RoleRepository
	keys           utils.Keys
}

func NewUserService(repository repositories.UserRepository, roleRepository repositories.RoleRepository, keys utils.Keys) *UserService {
	return &UserService{repository: repository, roleRepository: roleRepository, keys: keys}
}

func (us *UserService) SignUp(user models.User) (models.User, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return models.User{}, errors.New("Failed to hash password")
	}

	role, err := us.roleRepository.GetByName("default")
	if err != nil {
		return models.User{}, err
	}

	user.Password = hashedPassword
	user.Roles = append(user.Roles, role)
	return us.repository.Create(user)
}

func (us *UserService) GetUser(userID string) (models.User, error) {
	return us.repository.Get(userID)
}

func (us *UserService) ListUsers() ([]models.User, error) {
	return us.repository.GetAll()
}

func (us *UserService) ListUsersByRole(roleName string) ([]models.User, error) {
	return us.repository.GetAllByRole(roleName)
}

func (us *UserService) UpdateUser(userID string, user models.User) (models.User, error) {
	return us.repository.Update(userID, user)
}

func (us *UserService) ResolveRoles(names []string) ([]models.Role, error) {
	roles := make([]models.Role, len(names))
	for i, name := range names {
		role, err := us.roleRepository.GetByName(name)
		if err != nil {
			return nil, err
		}
		roles[i] = role
	}
	return roles, nil
}

func (us *UserService) LogIn(user models.User) (string, error) {
	token := ""
	user, err := us.repository.ValidateUser(user)
	if err != nil {
		return "", err
	}

	roles := make([]string, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = r.Name
	}

	token, err = us.keys.GenerateToken(user.ID, user.Email, roles)
	if err != nil {
		return "", errors.New("Failed to generate token")
	}

	return token, err
}

func (us *UserService) DeleteUser(userID string) error {
	return us.repository.Delete(userID)
}
