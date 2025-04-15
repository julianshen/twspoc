// main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	r "github.com/rethinkdb/rethinkdb-go"

	"notification/api"
	"notification/mongostore"
	"notification/notificationrepo"
	"notification/rethinkstore"
)

func main() {
	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Get database configuration from environment variables
	dbType := strings.ToLower(getEnv("DB_TYPE", "rethink")) // "rethink" or "mongo"
	dbName := getEnv("DB_NAME", "notifdb")

	// Initialize repository based on DB_TYPE
	var notifRepo notificationrepo.NotificationRepository
	var cleanup func()

	if dbType == "mongo" {
		// MongoDB setup
		mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
		log.Printf("Using MongoDB at %s...", mongoURI)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		repo, err := mongostore.NewMongoNotificationRepository(ctx, mongoURI, dbName, "notifications")
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		cleanup = func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := repo.Close(ctx); err != nil {
				log.Printf("Error closing MongoDB connection: %v", err)
			}
		}

		notifRepo = repo
		log.Println("Successfully connected to MongoDB")
	} else {
		// RethinkDB setup (default)
		rethinkAddr := getEnv("RETHINKDB_ADDR", "localhost:28015")
		log.Printf("Using RethinkDB at %s...", rethinkAddr)

		session, err := r.Connect(r.ConnectOpts{
			Address:  rethinkAddr,
			Database: dbName,
			Timeout:  10 * time.Second,
		})
		if err != nil {
			log.Fatalf("Failed to connect to RethinkDB: %v", err)
		}

		cleanup = func() {
			session.Close()
		}

		// Ensure database exists
		err = ensureDatabase(session, dbName)
		if err != nil {
			log.Fatalf("Failed to ensure database: %v", err)
		}

		// Ensure notifications table exists
		err = ensureTable(session, dbName, "notifications")
		if err != nil {
			log.Fatalf("Failed to ensure notifications table: %v", err)
		}

		notifRepo = rethinkstore.NewRethinkNotificationRepo(session)
		log.Println("Successfully connected to RethinkDB")
	}

	// Ensure cleanup happens when the program exits
	defer cleanup()

	// Set Gin to release mode for production, or comment out for debug
	// gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Register API handlers (handlers must be refactored to Gin style)
	router.POST("/api/notifications", api.GinCreateNotificationHandler(notifRepo))
	router.GET("/api/notifications", api.GinListNotificationsHandler(notifRepo))
	router.POST("/api/notifications/:id/read", api.GinMarkAsReadHandler(notifRepo))
	router.DELETE("/api/notifications/:id", api.GinDeleteNotificationHandler(notifRepo))
	// router.POST("/api/notifications/search", api.GinSearchHandler(notifRepo))
	router.GET("/api/notifications/subscribe", api.GinSSEHandler(notifRepo))

	// Start HTTP server in a goroutine
	go func() {
		port := getEnv("PORT", "3000")
		log.Printf("Notification service running on :%s (using %s)", port, dbType)
		log.Println("Available endpoints:")
		log.Println("- POST /api/notifications - Create a notification")
		log.Println("- GET /api/notifications?userId=xxx - Get notifications for a user")
		log.Println("- POST /api/notifications/:id/read - Mark a notification as read")
		log.Println("- DELETE /api/notifications/:id - Delete a notification")
		log.Println("- POST /api/notifications/search - Search notifications")
		log.Println("- GET /api/notifications/subscribe?userId=xxx - Subscribe to notification updates")

		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Gin server error: %v", err)
		}
	}()

	// Wait for termination signal
	<-signalChan
	log.Println("Shutdown signal received, initiating graceful shutdown...")

	// Gin does not provide a built-in shutdown in this pattern, but context cancellation will stop CDC and other background tasks.

	log.Println("Notification service shutdown complete")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ensureDatabase(session *r.Session, dbName string) error {
	res, err := r.DBList().Run(session)
	if err != nil {
		return err
	}
	defer res.Close()

	var dbs []string
	if err := res.All(&dbs); err != nil {
		return err
	}
	for _, db := range dbs {
		if db == dbName {
			return nil
		}
	}
	// Database does not exist, create it
	_, err = r.DBCreate(dbName).RunWrite(session)
	return err
}

// ensureTable checks if a table exists, and creates it if not
func ensureTable(session *r.Session, dbName, tableName string) error {
	res, err := r.DB(dbName).TableList().Run(session)
	if err != nil {
		return err
	}
	defer res.Close()

	var tables []string
	if err := res.All(&tables); err != nil {
		return err
	}
	for _, t := range tables {
		if t == tableName {
			return nil
		}
	}
	// Table does not exist, create it
	_, err = r.DB(dbName).TableCreate(tableName).RunWrite(session)
	return err
}
