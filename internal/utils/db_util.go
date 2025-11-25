package utils

import (
	"agnos-middleware/internal/configs"
	"agnos-middleware/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(config *configs.ApplicationConfig) (*gorm.DB, error) {
	var dsn string
	if config.Database.Password == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s dbname=%s port=%s sslmode=disable",
			config.Database.Host,
			config.Database.User,
			config.Database.DBName,
			config.Database.Port,
		)
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			config.Database.Host,
			config.Database.User,
			config.Database.Password,
			config.Database.DBName,
			config.Database.Port,
		)
	}

	log.Printf("Connecting to database: host=%s user=%s dbname=%s",
		config.Database.Host, config.Database.User, config.Database.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(
		&models.Staff{},
		&models.Patient{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("✅ Database connected and migrated successfully")
	return db, nil
}
