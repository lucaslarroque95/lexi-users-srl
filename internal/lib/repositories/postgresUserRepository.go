package repositories

import (
	"errors"
	"lexi-users-srl/internal/lib/db"
	"lexi-users-srl/internal/lib/models"
	"lexi-users-srl/internal/lib/utils"

	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	gormDB *gorm.DB
}

func NewPostgresUserRepository(gormDB *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{gormDB: gormDB}
}

func (pur *PostgresUserRepository) Create(user models.User) (models.User, error) {
	ormUser := toUserORM(user)
	if err := pur.gormDB.Create(&ormUser).Error; err != nil {
		return models.User{}, err
	}
	return toUserDomain(ormUser), nil
}

func (pur *PostgresUserRepository) Get(id string) (models.User, error) {
	var ormUser db.User
	if err := pur.gormDB.Preload("Roles").Where("id = ?", id).First(&ormUser).Error; err != nil {
		return models.User{}, err
	}
	return toUserDomain(ormUser), nil
}

func (pur *PostgresUserRepository) ValidateUser(user models.User) (models.User, error) {
	var ormUser db.User
	err := pur.gormDB.Preload("Roles").Where("email = ?", user.Email).First(&ormUser).Error
	if err != nil {
		return models.User{}, errors.New("invalid email")
	}
	isValidUser := utils.CheckPasswordHash(user.Password, ormUser.Password)
	if !isValidUser {
		return models.User{}, errors.New("invalid password")
	}
	return toUserDomain(ormUser), nil
}

func (pur *PostgresUserRepository) GetAll() ([]models.User, error) {
	var ormUsers []db.User
	if err := pur.gormDB.Find(&ormUsers).Error; err != nil {
		return nil, err
	}

	users := make([]models.User, len(ormUsers))
	for i, ormUser := range ormUsers {
		users[i] = toUserDomain(ormUser)
	}
	return users, nil
}

func (pur *PostgresUserRepository) GetAllByRole(roleName string) ([]models.User, error) {
	var ormUsers []db.User
	err := pur.gormDB.
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.name = ?", roleName).
		Find(&ormUsers).Error
	if err != nil {
		return nil, err
	}

	users := make([]models.User, len(ormUsers))
	for i, ormUser := range ormUsers {
		users[i] = toUserDomain(ormUser)
	}
	return users, nil
}

func (pur *PostgresUserRepository) Update(id string, user models.User) (models.User, error) {
	ormUser := toUserORM(user)
	ormUser.ID = id

	if err := pur.gormDB.Model(&db.User{}).Where("id = ?", id).Updates(&ormUser).Error; err != nil {
		return models.User{}, err
	}
	if err := pur.gormDB.Model(&db.User{ID: id}).Association("Roles").Replace(ormUser.Roles); err != nil {
		return models.User{}, err
	}
	return pur.Get(id)
}

func (pur *PostgresUserRepository) Delete(id string) error {
	return pur.gormDB.Where("id = ?", id).Delete(&db.User{}).Error
}

func toUserDomain(ormUser db.User) models.User {
	roles := make([]models.Role, len(ormUser.Roles))
	for i, ormRole := range ormUser.Roles {
		roles[i] = toRoleDomain(ormRole)
	}

	return models.User{
		ID:       ormUser.ID,
		Email:    ormUser.Email,
		Password: ormUser.Password,
		Roles:    roles,
	}
}

func toUserORM(user models.User) db.User {
	roles := make([]db.Role, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = toRoleORM(role)
	}

	return db.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Roles:    roles,
	}
}
