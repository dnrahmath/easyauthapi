package configs

import (
	"errors"
	"fmt"
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

func ConnectDB() error {
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
		return errors.New("Failed to connect to the Database: " + err.Error())
	}

	fmt.Println("🚀 Connected Successfully to the Database")
	return nil
}
