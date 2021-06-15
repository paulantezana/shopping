package provider

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GetConnection get connection database
func GetConnection() *gorm.DB {
	c := GetConfig()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", c.Database.Server, c.Database.User, c.Database.Pass, c.Database.Database, c.Database.Port)
		// dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.Database.User, c.Database.Pass, c.Database.Server, c.Database.Port, c.Database.Database)
	}

	// db, err := gorm.Open("postgres", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
