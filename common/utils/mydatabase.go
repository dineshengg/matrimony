package utils

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MyDatabase struct {
	Database     *gorm.DB
	dbSourceName string
	dbErrCount   int
}

var db *MyDatabase
var mtx sync.Mutex
var closing bool

func IsPostgresInit() bool {
	mtx.Lock()
	defer mtx.Unlock()
	if db != nil && db.Database != nil {
		return true
	}
	return false
}

// NewDatabaseClient initializes a new GORM database client
func NewDatabaseClient(datasourceName string) *MyDatabase {
	mtx.Lock()
	defer mtx.Unlock()
	db = &MyDatabase{
		Database:     nil,
		dbSourceName: datasourceName,
		dbErrCount:   0,
	}
	//Initialize the database connection
	if db.initializeDB() != nil {
		log.Error("failed to connect to postgresql database")
		return nil
	}
	return db
}

// GetDB returns the GORM database instance
func GetDB() *gorm.DB {
	return db.GetDB()
}

// GetDB ensures the database connection is initialized and returns it
func (db *MyDatabase) GetDB() *gorm.DB {

	//dont need to lock and unlock every time most of the time DB handle will be valid only during closing need to take care of check DB handle validity
	if closing == true {
		defer mtx.Unlock()
		mtx.Lock()
	}

	if db.Database != nil {
		return db.Database
	}

	if db.dbErrCount <= 5 {
		defer mtx.Unlock()
		mtx.Lock()
		if err := db.initializeDB(); err != nil {
			db.dbErrCount++
			fmt.Println("failed to connect to PostgreSQL database")
			return nil
		}
	} else {
		fmt.Println("max retry occurred - failed to connect to PostgreSQL database")
	}

	return db.Database
}

// InitializeDB initializes the PostgreSQL connection using GORM
func (db *MyDatabase) initializeDB() error {
	var err error

	if db.Database != nil {
		return nil
	}
	log.Println("dsn - ", db.dbSourceName)
	// Configure GORM with PostgreSQL driver
	db.Database, err = gorm.Open(postgres.Open(db.dbSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable detailed logging
	})
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	// Configure connection pooling
	sqlDB, err := db.Database.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(20)                  // Set maximum number of open connections
	sqlDB.SetMaxIdleConns(15)                  // Set maximum number of idle connections
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Set maximum connection lifetime

	// Test the connection
	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully using GORM!")
	return nil
}

// CloseDB closes the database connection
func (db *MyDatabase) CloseDB() {
	closing = true
	mtx.Lock()
	defer mtx.Unlock()
	if db != nil && db.Database != nil {
		sqlDB, err := db.Database.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("PostgreSQL connection closed")
		}
	}
}
