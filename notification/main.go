// main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	r "github.com/rethinkdb/rethinkdb-go"

	"notification/api"
	"notification/rethinkstore"
)

func main() {
	// Set up signal handling for graceful shutdown

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Get RethinkDB address and database name from environment variable or use default
	rethinkAddr := getEnv("RETHINKDB_ADDR", "localhost:28015")
	dbName := getEnv("DB_NAME", "notifdb")

	// Connect to RethinkDB
	log.Printf("Connecting to RethinkDB at %s...", rethinkAddr)
	session, err := r.Connect(r.ConnectOpts{
		Address:  rethinkAddr,
		Database: dbName,
		Timeout:  10 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to RethinkDB: %v", err)
	}
	defer session.Close()
	log.Println("Successfully connected to RethinkDB")

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

	// Initialize repository
	notifRepo := rethinkstore.NewRethinkNotificationRepo(session)

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
		log.Println("Notification service running on :3000")
		log.Println("Available endpoints:")
		log.Println("- POST /api/notifications - Create a notification")
		log.Println("- GET /api/notifications?userId=xxx - Get notifications for a user")
		log.Println("- POST /api/notifications/:id/read - Mark a notification as read")
		log.Println("- DELETE /api/notifications/:id - Delete a notification")
		log.Println("- POST /api/notifications/search - Search notifications")
		log.Println("- GET /api/notifications/subscribe?userId=xxx - Subscribe to notification updates")

		if err := router.Run(":3000"); err != nil {
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
