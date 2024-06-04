package config

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

var Name string
var Message string
var Owner string
var OwnerToken string
var Version string
var Port string
var DatabaseDriver string
var DSN string

func init() {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading Environment Variables")
	}

	log.Default().Println("Loaded Environment Variables")

	Name = os.Getenv("NAME")
	if Name == "" {
		Name = "Eskimoe"
	}
	if len(Name) > 64 {
		log.Fatal("Server Name is too long. Maximum 64 characters")
	}
	if !isAlpha(Name) {
		log.Fatal("Server Name must be alphanumeric")
	}
	log.Default().Println("Server Name:", Name)

	Message = os.Getenv("MESSAGE")
	if Message == "" {
		Message = "Eskimoe Chat Server"
	}
	if len(Message) > 256 {
		log.Fatal("Server Message is too long. Maximum 256 characters")
	}
	log.Default().Println("Server Message:", Message)

	Version = "0.1.0"

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8000"
	}
	log.Default().Println("Server Port:", Port)

	Owner = os.Getenv("OWNER_ID")
	if Owner == "" {
		log.Fatal("Owner ID not found in Environment Variables")
	}

	OwnerToken = os.Getenv("OWNER_TOKEN")
	if OwnerToken == "" {
		log.Fatal("Owner Token not found in Environment Variables")
	}

	DatabaseDriver = os.Getenv("DATABASE_DRIVER")
	if DatabaseDriver == "" {
		DatabaseDriver = "sqlite"
	}

	DSN = os.Getenv("DSN")
	if DSN == "" {
		log.Fatal("DSN not found in Environment Variables")
	}
}
