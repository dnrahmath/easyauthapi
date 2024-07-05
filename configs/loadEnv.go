package configs

import (
	"easyauthapi/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var UseConfig *models.Config

//=============================================================================

func LoadConfigGodotenv() {
	err := godotenv.Load() // Load variabel lingkungan dari file .env
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	UseConfig = &models.Config{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		PGHOST:       os.Getenv("PGHOST"),
		PGUSER:       os.Getenv("PGUSER"),
		PGPASSWORD:   os.Getenv("PGPASSWORD"),
		PGDATABASE:   os.Getenv("PGDATABASE"),
		PGPORT:       os.Getenv("PGPORT"),
	}
}

//=============================================================================

func LoadConfigViper(path string, configFile *string) {
	viper := viper.New()
	viper.AutomaticEnv()
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("MODE", "debug")
	viper.SetConfigType("env")
	viper.SetConfigName(*configFile) // Gunakan nama file yang diberikan
	viper.AddConfigPath(path)        // Gunakan path yang diberikan

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&UseConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	if err := UseConfig.Validate(); err != nil {
		log.Fatalf("Error validating configuration: %v", err)
	}
}

//=============================================================================
