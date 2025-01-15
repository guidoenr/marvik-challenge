package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// User defines the data for a User, and the `gorm`|`json` tags to map in the table later and retrieve info from the api.
type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:50;not null" json:"name"`
	Surname       string         `gorm:"size:50;not null" json:"surname"`
	Email         string         `gorm:"size:100;unique;not null" json:"email"`
	Organizations []Organization `gorm:"many2many:user_organizations;" json:"organizations"`
}

// Organization defines the data for the organization and the `gorm`|`json` tags
type Organization struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"size:100;not null" json:"name"`
	Users []User `gorm:"many2many:user_organizations;" json:"users"`
}

// connectToDb connects to the postgresql database
func connectToDb() {
	log.Debug().Msg("connecting to PostgreSQL db ...")
	var err error

	// first get the ENV vars with the POSTGRESQL
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// endpoint to connect to the DB
	pgconn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable", user, password, dbname)

	// create the postgresql dialector to connect
	db, err = gorm.Open(postgres.Open(pgconn))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	// run the migrations and log the result
	err = db.AutoMigrate(&User{}, &Organization{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	} else {
		log.Info().Msg("migrations done successfully")
	}
}
