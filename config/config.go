package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
}

// Load membaca variabel environment dari .env atau OS
func Load() *Config {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("⚠️ .env file not found, using system environment variables")
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}

// ConnectDB membuat koneksi ke database PostgreSQL
func ConnectDB(cfg *Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Failed to open DB:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("❌ Failed to ping DB:", err)
	}

	log.Println("✅ Connected to the database successfully")
	return db
}
