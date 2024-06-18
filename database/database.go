package database

import (
	"eskimoe-server/config"
	"fmt"
	"log"

	"gorm.io/datatypes"
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
		&ServerReaction{},
		&ServerInvite{},
		&ServerPermission{},
		&Member{},
		&Role{},
		&RoomCategory{},
		&Room{},
		&RoomPermission{},
		&Message{},
		&MessageAttachment{},
		&MessageReaction{},
	)

	// Setup the server if it doesn't exist
	if !SetupServer() {
		log.Fatal("Error Setting Up Server")
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

	// Create a new server
	newServer := Server{
		Name:    config.Name,
		Message: config.Message,
		Mode:    Open, // Default to open mode
	}

	// Create default category and room
	generalCategory := RoomCategory{
		Name:        "General",
		Description: "General discussion",
	}

	chatRoom := Room{
		Name:        "Chat",
		Description: "General chat room",
		Type:        Text,
	}

	// Create default role and permissions
	everyoneRole := Role{
		Name: "Everyone",
		Permissions: []ServerPermission{
			{Permission: ViewRoom},
			{Permission: SendMessage},
			{Permission: AddLink},
			{Permission: AddFile},
			{Permission: AddReaction},
			{Permission: ChangeName},
			{Permission: RunCommands},
			{Permission: ViewMessageHistory},
		},
	}

	// Create relationships
	generalCategory.Rooms = []Room{chatRoom}
	newServer.Roles = []Role{everyoneRole}
	newServer.RoomCategories = []RoomCategory{generalCategory}

	// Save the server
	if Database.Create(&newServer).Error != nil {
		return false
	}

	// Save the order
	newServer.CategoryOrder = datatypes.JSON([]byte(fmt.Sprintf("[%d]", generalCategory.ID)))
	newServer.RoleOrder = datatypes.JSON([]byte(fmt.Sprintf("[%d]", everyoneRole.ID)))
	generalCategory.RoomOrder = datatypes.JSON([]byte(fmt.Sprintf("[%d]", chatRoom.ID)))

	return Database.Save(&newServer).Error == nil
}
