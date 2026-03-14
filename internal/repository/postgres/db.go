package postgres

import (
	"fmt"
	"log"

	"github.com/akshay/productiv-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB creates a new GORM database connection and runs auto-migrations.
func NewDB(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	log.Println("database migrations completed")
	return db, nil
}

// runMigrations auto-migrates all domain models.
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Pillar{},
		&domain.FastingSession{},
		&domain.GymSession{},
		&domain.MeditationSession{},
		&domain.RetentionStreak{},
	)
}

// Seed inserts the default user and pillars if they don't exist.
func Seed(db *gorm.DB) error {
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	if userCount > 0 {
		return nil
	}

	user := domain.User{
		Name:     "Akshay",
		Timezone: "Asia/Kolkata",
	}
	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("seeding user: %w", err)
	}

	pillars := []domain.Pillar{
		{UserID: user.ID, Type: domain.PillarFasting, Name: "Fasting", Icon: "🍽️", Color: "#e94560", IsActive: true, DisplayOrder: 1},
		{UserID: user.ID, Type: domain.PillarGym, Name: "Gym", Icon: "💪", Color: "#4a9eff", IsActive: true, DisplayOrder: 2},
		{UserID: user.ID, Type: domain.PillarMeditation, Name: "Meditation", Icon: "🧘", Color: "#9b59b6", IsActive: true, DisplayOrder: 3},
		{UserID: user.ID, Type: domain.PillarRetention, Name: "Retention", Icon: "🔥", Color: "#2ecc71", IsActive: true, DisplayOrder: 4},
	}
	if err := db.Create(&pillars).Error; err != nil {
		return fmt.Errorf("seeding pillars: %w", err)
	}

	log.Println("database seeded with default user and pillars")
	return nil
}
