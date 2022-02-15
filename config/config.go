package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/memandip/go-examples/mysql"

	"github.com/joho/godotenv"
)

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return err
}

func GetDBCredentials() mysql.DatabaseCredentials {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	return mysql.DatabaseCredentials{
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_HOST"),
	}
}

func Dbpool(dbCreds mysql.DatabaseCredentials) (*sql.DB, error) {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	var pool *sql.DB

	if dbCreds != (mysql.DatabaseCredentials{}) {
		pool, err = mysql.Connect(dbCreds)
	} else {
		pool, err = mysql.Connect(GetDBCredentials())
	}

	return pool, err
}
