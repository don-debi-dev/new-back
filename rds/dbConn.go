package rds

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type settings struct {
	dbHost         string
	dbUser         string
	dbPassword     string
	dbName         string
	maxConnections int
	maxIdleConns   int
}

// GetDB connects and returns a RDS connection based on environment variables.
func GetDB() (*sql.DB, error) {
	sett, err := getSettings()
	if err != nil {
		log.Printf("Error getting env variables: %s\n", err.Error())
		return nil, err
	}

	return GetDBFromSettings(sett)
}

// GetDBFromInfo connects and returns a RDS connection based on given information.
func GetDBFromInfo(hostAndPort, user, password, DBName string) (*sql.DB, error) {
	sett := settings{
		dbHost:     hostAndPort,
		dbUser:     user,
		dbPassword: password,
		dbName:     DBName,
	}
	return GetDBFromSettings(sett)
}

// GetDBFromSettings connects and returns a RDS connection based on settings.
func GetDBFromSettings(sett settings) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", sett.dbUser, sett.dbPassword, sett.dbHost, sett.dbName)

	log.Print(connStr)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Printf("Error connecting to Db Error: %s\n", err.Error())
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging DB Error: %s\n", err.Error())
		return nil, err
	}

	return db, nil
}

func getSettings() (settings, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return settings{}, err
	}

	return settings{
		dbHost:         os.Getenv("DB_HOST"),
		dbUser:         os.Getenv("DB_USER"),
		dbPassword:     os.Getenv("DB_PASSWORD"),
		dbName:         os.Getenv("DB_NAME"),
		maxConnections: 3,
		maxIdleConns:   1,
	}, nil
}
