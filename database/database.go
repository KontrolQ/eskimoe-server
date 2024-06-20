package database

import (
	"eskimoe-server/config"
	"log"
	"os"

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

	// Migrate the schema
	Database.AutoMigrate(
		&Server{},
		&Category{},
		&Room{},
		&Message{},
		&MessageReaction{},
		&MessageAttachment{},
		&ServerReaction{},
		&Invite{},
		&Role{},
		&Event{},
		&Log{},
		&Member{},
	)

	// Setup the server if it doesn't exist
	if !SetupServer() {
		// Delete Database
		Database.Migrator().DropTable(
			&Server{},
			&Category{},
			&Room{},
			&Message{},
			&MessageReaction{},
			&MessageAttachment{},
			&ServerReaction{},
			&Invite{},
			&Role{},
			&Event{},
			&Log{},
			&Member{},
		)

		// If Database is sqlite, delete the file
		if driver == "sqlite" {
			if error = os.Remove(dsn); error != nil {
				log.Fatal("Error Removing Database")
			}
		}

		log.Fatal("Error Setting Up Server. Removing Database")
	}

	log.Default().Println("Server Setup Complete")

}

// Setups a default server for first time use - returns false if a server already exists
func SetupServer() bool {
	var server Server
	if Database.First(&server).Error == nil {
		server.Name = config.Name
		server.Message = config.Message

		return Database.Save(&server).Error == nil
	}

	newServer := Server{
		Name:    config.Name,
		Message: config.Message,
		Mode:    Open,
	}

	if Database.Create(&newServer).Error != nil {
		return false
	}

	// Update Server with default values

	likeReaction := ServerReaction{
		Reaction: "LIKE",
		ServerID: newServer.ID,
	}

	if Database.Create(&likeReaction).Error != nil {
		return false
	}

	// Create General Category
	generalCategory := Category{
		Name:     "General",
		ServerID: newServer.ID,
	}

	if Database.Create(&generalCategory).Error != nil {
		return false
	}

	generalCategoryID := generalCategory.ID

	// Create General Chat Room
	generalChatRoom := Room{
		Name:        "General",
		Description: "General Chat Room",
		Type:        Text,
		CategoryID:  generalCategoryID,
	}

	if Database.Create(&generalChatRoom).Error != nil {
		return false
	}

	// Set Room Order, Category Order
	generalCategory.RoomOrder = []int{generalChatRoom.ID}
	newServer.CategoryOrder = []int{generalCategory.ID}

	if Database.Save(&generalCategory).Error != nil {
		return false
	}

	// Create Everyone Role
	everyoneRole := Role{
		Name:        "everyone",
		Permissions: []Permission{SendMessage, AddLink, AddFile, AddReaction, RunCommands, ViewMessageHistory, GenerateInvites},
		SystemRole:  true,
	}

	if Database.Create(&everyoneRole).Error != nil {
		return false
	}

	// Set Role Order
	newServer.RoleOrder = []int{everyoneRole.ID}

	return Database.Save(&newServer).Error == nil
}
