package config

import (
	"database/sql"
	"log"
	"mandip/go-examples/mysql"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return err
}

func Dbpool() (*sql.DB, error) {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}
	pool, err := mysql.Connect(mysql.DatabaseCredentials{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_HOST"),
	})

	return pool, err
}
