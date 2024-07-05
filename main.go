package main

import (
	initi "easyauthapi/configs"
	gen "easyauthapi/docs"
	"easyauthapi/models/migration"
	"easyauthapi/routes"
	"flag"
	"fmt"
	// "easyauthapi/services"

	// "easyauthapi/tools"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//===============================================================

// @title GoLang Rest API Starter Doc
// @version 1.0
// @description GoLang - Gin - RESTful - MongoDB - Redis
// @termsOfService https://swagger.io/terms/

// @contact.name Ebubekir Yiğit
// @contact.url https://github.com/ebubekiryigit
// @contact.email ebubekiryigit6@gmail.com

// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Access-Token
func serverStart(configFile *string) {
	// Load configuration from file
	// initi.LoadConfigViper("./", configFile)
	initi.LoadConfigGodotenv()

	// Connect to PostgreSQL database
	initi.ConnectDB()

	// Migrate database tables if they don't exist
	data := migration.AdmissionMigration{}
	data.Migrate()
	data2 := migration.TokenMigration{}
	data2.Migrate()
	data3 := migration.UserMigration{}
	data3.Migrate()

	fmt.Println("Database migrations completed")

	// Initialize Gin router and server
	routes.InitGin()
	router := routes.New()

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Handler:      router,
	}

	// Start HTTP server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Server listen error: %s\n", err)
		}
	}()

	// Handle graceful shutdown on interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}
	log.Println("Server exited")
}

//===============================================================

func main() {
	execute := flag.String("exec", "", "name of the exec to run")
	configFile := flag.String("config", "app.env", "Name of the config file (without extension)")
	flag.Parse()

	fmt.Printf("Using config file: %s\n", *configFile)

	switch *execute {
	case "CreateDocs":
		gen.CreateDocs()
	case "":
		serverStart(configFile)
	default:
		fmt.Println("Fungsi tidak dikenali. Gunakan 'CreateDocs' atau biarkan kosong untuk menjalankan serverStart.")
	}
}

//===============================================================
