package storage

import (
	"BE-HA001/pkg"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	*gorm.DB
	cache Cache
}

func NewStorage() (*Storage, error) {
	cfg := pkg.LoadConfig()

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName:           "pgx",
		DSN:                  cfg.DSN(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Error connect: %v", err)
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return &Storage{
		DB:    db,
		cache: &CacheMock{},
	}, nil
}

func (s *Storage) Close() {
	if s.DB != nil {
		sqlDB, _ := s.DB.DB()
		sqlDB.Close()
	}
}
