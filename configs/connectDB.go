package configs

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	JWTSecretKey                             = "My.Ultra.Secure.Password"
	JWTAccessExpirationMinutes time.Duration = 1440
	JWTRefreshExpirationDays   time.Duration = 7
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	// Construct DSN using Railway environment variables
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		UseConfig.PGHOST,
		UseConfig.PGPORT,
		UseConfig.PGUSER,
		UseConfig.PGPASSWORD,
		UseConfig.PGDATABASE,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database:", err)
	}

	fmt.Println("🚀 Connected Successfully to the Database")
}
