package pkg

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect establishes a connection to the database with optimized settings
// and returns a *gorm.DB instance.
func DatabaseConnect(dbUser, dbPassword, dbName, dbHost, dbPort string) (*gorm.DB, error) {
	// Connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open a GORM connection to MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Retrieve the underlying sql.DB to configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("could not get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)               // Adjust based on workload
	sqlDB.SetMaxOpenConns(100)              // Adjust based on workload
	sqlDB.SetConnMaxLifetime(1 * time.Hour) // Prevent stale connections

	log.Printf("Database connection established")

	return db, nil
}
