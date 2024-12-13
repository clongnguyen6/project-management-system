package main

import (
	"context"
	"errors"
	"example/project-management-system/internal/config"
	"example/project-management-system/internal/database"
	"example/project-management-system/internal/server"
	"example/project-management-system/pkg/logger"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "example/project-management-system/docs"
)

//	@title			Project Management System API
//	@version		1.0
//	@description	API for managing projects, tasks, teams, and users
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
func main() {
	// Load config env
	cfg := config.LoadEnvConfigs()
	fmt.Println(cfg)

	// Initialize logger
	appLogger := logger.NewLogger()

	// Initialize database
	db := database.NewPostgresConnection(cfg, appLogger)

	// Initialize server
	server, err := server.NewHTTPServer(db.DB, cfg)
	if err != nil {
		log.Fatal(err)
	}


	log.Printf("API server listening on %s", server.Addr)

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("API server closed: err: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("got shutdown signal. shutting down server...")

	localCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(localCtx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("server shutdown complete")

}
