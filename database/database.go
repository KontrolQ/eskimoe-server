package database

import (
	"eskimoe-server/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Initialize() {
	driver := config.DatabaseDriver
	dsn := config.DSN
	var error error
	switch driver {
	case "sqlite":
		Database, error = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if error != nil {
			log.Fatal("Error Connecting to SQLite Database")
		}
	case "mysql":
		Database, error = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if error != nil {
			log.Fatal("Error Connecting to MySQL Database")
		}
	case "postgres":
		Database, error = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if error != nil {
			log.Fatal("Error Connecting to PostgreSQL Database")
		}
	case "mssql":
		Database, error = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if error != nil {
			log.Fatal("Error Connecting to MSSQL Database")
		}
	default:
		log.Fatal("Unsupported Database Driver")
	}

	log.Default().Println("Connected to Database")

	// Migrate the Schema
	Database.AutoMigrate(&Member{}, &Role{})
}
