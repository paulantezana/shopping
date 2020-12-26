package provider

import (
	"encoding/json"
	"log"
	"os"
)

// Database models
type Database struct {
	Server   string
	Port     string
	User     string
	Pass     string
	Database string
}

// Email config
type Email struct {
	Name     string
	From     string
	Password string
	Server   string
	Host     string
}

// Global config
type Global struct {
	PageLimit uint
    QueryApi string
}

// Server condif
type Server struct {
	Port string
	Key  string
}

// Config config
type Config struct {
    Database Database
    Email    Email
    Server   Server
    Global   Global
}

// GetConfig return configuration from database json
func GetConfig() Config {
	var c Config

	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}
