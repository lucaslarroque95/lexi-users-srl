package db

type User struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email    string `gorm:"uniqueIndex;not null" binding:"required"`
	Password string `gorm:"not null" binding:"required"`
	Roles    []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	ID    int64  `gorm:"primaryKey"`
	Name  string `gorm:"uniqueIndex;not null" binding:"required"`
	Users []User `gorm:"many2many:user_roles;"`
}
