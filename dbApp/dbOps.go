package dbApp

import (
	"fmt"
	"github.com/joho/godotenv" // Import for loading environment variables
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func (db *GormDB) Connect() error {
	err := godotenv.Load() // Load environment variables from .env file (optional)
	if err != nil {
		return fmt.Errorf("error loading environment variables: %w", err)
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Build the connection string using environment variables
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUsername, dbPassword, dbName, dbPort)

	// Configure connection options (optional)
	config := &gorm.Config{}
	// Add other configuration options that can be made:
	// config.Logger = logger.Default // Set a custom logger

	dbClient, err := gorm.Open(postgres.Open(connectionString), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set automigrations mode (optional)
	err = dbClient.AutoMigrate(&User{}, &Organization{}, &Machine{}, &Transactions{})
	if err != nil {
		fmt.Println("Error occured during migration: ", err)
	}

	db.DB = dbClient
	return nil
}
