package testutil

import (
	"errors"
	"strconv"

	"lexi-users-srl/internal/lib/models"
	"lexi-users-srl/internal/lib/repositories"
	"lexi-users-srl/internal/lib/utils"
)

var (
	_ repositories.UserRepository = (*FakeUserRepository)(nil)
	_ repositories.RoleRepository = (*FakeRoleRepository)(nil)
)

// FakeUserRepository is an in-memory repositories.UserRepository for tests, no database required.
type FakeUserRepository struct {
	Users  map[string]models.User
	NextID int

	CreateErr       error
	GetErr          error
	ValidateErr     error
	GetAllErr       error
	GetAllByRoleErr error
	UpdateErr       error
	DeleteErr       error

	CreateCalls int
}

func NewFakeUserRepository() *FakeUserRepository {
	return &FakeUserRepository{Users: make(map[string]models.User)}
}

func (f *FakeUserRepository) Create(user models.User) (models.User, error) {
	f.CreateCalls++
	if f.CreateErr != nil {
		return models.User{}, f.CreateErr
	}

	f.NextID++
	user.ID = strconv.Itoa(f.NextID)
	f.Users[user.ID] = user
	return user, nil
}

func (f *FakeUserRepository) Get(id string) (models.User, error) {
	if f.GetErr != nil {
		return models.User{}, f.GetErr
	}

	user, ok := f.Users[id]
	if !ok {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (f *FakeUserRepository) ValidateUser(user models.User) (models.User, error) {
	if f.ValidateErr != nil {
		return models.User{}, f.ValidateErr
	}

	for _, existing := range f.Users {
		if existing.Email != user.Email {
			continue
		}
		if !utils.CheckPasswordHash(user.Password, existing.Password) {
			return models.User{}, errors.New("invalid password")
		}
		return existing, nil
	}
	return models.User{}, errors.New("invalid email")
}

func (f *FakeUserRepository) GetAll() ([]models.User, error) {
	if f.GetAllErr != nil {
		return nil, f.GetAllErr
	}

	result := make([]models.User, 0, len(f.Users))
	for _, user := range f.Users {
		result = append(result, user)
	}
	return result, nil
}

func (f *FakeUserRepository) GetAllByRole(roleName string) ([]models.User, error) {
	if f.GetAllByRoleErr != nil {
		return nil, f.GetAllByRoleErr
	}

	var result []models.User
	for _, user := range f.Users {
		for _, role := range user.Roles {
			if role.Name == roleName {
				result = append(result, user)
				break
			}
		}
	}
	return result, nil
}

func (f *FakeUserRepository) Update(id string, user models.User) (models.User, error) {
	if f.UpdateErr != nil {
		return models.User{}, f.UpdateErr
	}
	if _, ok := f.Users[id]; !ok {
		return models.User{}, errors.New("user not found")
	}

	user.ID = id
	f.Users[id] = user
	return user, nil
}

func (f *FakeUserRepository) Delete(id string) error {
	if f.DeleteErr != nil {
		return f.DeleteErr
	}
	delete(f.Users, id)
	return nil
}

// FakeRoleRepository is an in-memory repositories.RoleRepository for tests, no database required.
type FakeRoleRepository struct {
	Roles  map[int64]models.Role
	NextID int64
}

func NewFakeRoleRepository() *FakeRoleRepository {
	return &FakeRoleRepository{Roles: make(map[int64]models.Role)}
}

func (f *FakeRoleRepository) Create(role models.Role) (models.Role, error) {
	f.NextID++
	role.ID = f.NextID
	f.Roles[role.ID] = role
	return role, nil
}

func (f *FakeRoleRepository) Get(id int64) (models.Role, error) {
	role, ok := f.Roles[id]
	if !ok {
		return models.Role{}, errors.New("role not found")
	}
	return role, nil
}

func (f *FakeRoleRepository) GetAll() ([]models.Role, error) {
	result := make([]models.Role, 0, len(f.Roles))
	for _, role := range f.Roles {
		result = append(result, role)
	}
	return result, nil
}

func (f *FakeRoleRepository) GetByName(name string) (models.Role, error) {
	for _, role := range f.Roles {
		if role.Name == name {
			return role, nil
		}
	}
	return models.Role{}, errors.New("role not found")
}

func (f *FakeRoleRepository) Update(id int64, role models.Role) (models.Role, error) {
	if _, ok := f.Roles[id]; !ok {
		return models.Role{}, errors.New("role not found")
	}

	role.ID = id
	f.Roles[id] = role
	return role, nil
}

func (f *FakeRoleRepository) Delete(id int64) error {
	delete(f.Roles, id)
	return nil
}
