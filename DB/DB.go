package DB

import (
	"database/sql"
	"fmt"
	"os"
)

var (
	DBCon *sql.DB
)

func InitDB() {
	// Connect to the postgres db
	RDS_USERNAME := os.Getenv("RDS_USERNAME")
	RDS_PASSWORD := os.Getenv("RDS_PASSWORD")
	RDS_DB_NAME := os.Getenv("RDS_DB_NAME")
	RDS_HOSTNAME := os.Getenv("RDS_HOSTNAME")
	RDS_PORT := os.Getenv("RDS_PORT")

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable host=%s port=%s password=%s",
		RDS_USERNAME, RDS_DB_NAME, RDS_HOSTNAME, RDS_PORT, RDS_PASSWORD)

	var err error
	DBCon, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}
