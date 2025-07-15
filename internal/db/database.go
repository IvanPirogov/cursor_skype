package db

import (
	"fmt"
	"log"
	"messenger/internal/config"
	"messenger/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(config *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.Host,
		config.Username,
		config.Password,
		config.Database,
		config.Port,
		config.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
	log.Println("Running database migrations...")
	
	err := d.DB.AutoMigrate(
		&models.User{},
		&models.UserSession{},
		&models.Contact{},
		&models.Chat{},
		&models.ChatMember{},
		&models.ChatSettings{},
		&models.ChatInvite{},
		&models.Message{},
		&models.MessageRead{},
		&models.File{},
		&models.Reaction{},
		&models.Call{},
		&models.CallParticipant{},
		&models.CallSettings{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) Ping() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}