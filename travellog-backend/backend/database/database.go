package database

import (
	"fmt"
	"log"
	"os"

	"github.com/travellog/backend/config"
	"github.com/travellog/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var err error
	var db *gorm.DB

	// Check if we should use SQLite (for development without Docker)
	useSQLite := os.Getenv("USE_SQLITE") == "true"

	if useSQLite {
		log.Println("Using SQLite database for development")
		db, err = gorm.Open(sqlite.Open("travellog.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		log.Println("Using PostgreSQL database")
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Hotel{}, &models.Flight{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	DB = db
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
