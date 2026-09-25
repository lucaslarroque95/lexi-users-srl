package db

import (
	"errors"
	"fmt"
	"os"

	"lexi-users-srl/internal/lib/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(buildDSN()), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database.")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic("Could not configure database connection pool.")
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	if err := DB.AutoMigrate(&User{}, &Role{}); err != nil {
		panic("Could not migrate database schema.")
	}

	if err := seedAdminUser(); err != nil {
		panic("Could not seed admin user.")
	}

	if err := seedDefaultRole(); err != nil {
		panic("Could not seed default role.")
	}
}

func seedAdminUser() error {
	adminEmail := getEnv("ADMIN_EMAIL", "admin@lexi.com")

	var existing User
	err := DB.Where("email = ?", adminEmail).First(&existing).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var adminRole Role
	if err := DB.FirstOrCreate(&adminRole, Role{Name: "admin"}).Error; err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword("admin")
	if err != nil {
		return err
	}

	admin := User{
		Email:    adminEmail,
		Password: hashedPassword,
		Roles:    []Role{adminRole},
	}
	return DB.Create(&admin).Error
}

func seedDefaultRole() error {
	var existing Role
	err := DB.Where("name = ?", "default").First(&existing).Error
	if err == nil {
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	role := Role{
		Name: "default",
	}
	return DB.Create(&role).Error
}

func buildDSN() string {
	host := getEnv("POSTGRES_SERVER", "localhost")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "postgres")
	password := getEnv("POSTGRES_PASSWORD", "postgres")
	name := getEnv("POSTGRES_DB", "lexi_users")
	sslMode := getEnv("POSTGRES_SSLMODE", "disable")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, name, sslMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
