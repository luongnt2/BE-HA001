package pkg

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

func LoadConfig() DBConfig {
	err := godotenv.Load("config.env") // Đọc file .env
	if err != nil {
		log.Println("env not found")
	}
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "yourpassword"),
		DBName:   getEnv("DB_NAME", "yourdb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
		TimeZone: getEnv("DB_TIMEZONE", "Asia/Ho_Chi_Minh"),
	}
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode, c.TimeZone,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}
