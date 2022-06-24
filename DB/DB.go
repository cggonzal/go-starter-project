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
	PGUSER := os.Getenv("PGUSER")
	PGPASSWORD := os.Getenv("PGPASSWORD")
	PGDATABASE := os.Getenv("PGDATABASE")
	PGHOST := os.Getenv("PGHOST")
	PGPORT := os.Getenv("PGPORT")

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable host=%s port=%s password=%s",
		PGUSER, PGDATABASE, PGHOST, PGPORT, PGPASSWORD)

	var err error
	DBCon, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}
