package DB

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

var (
	DBCon *sql.DB
)

func handleMigrations() {
	// DB contains a table called "migration" which contains one row that holds the number of the last applied migration

	// if migration table doesn't exist, create it and set migration number to 0
	_, err := DBCon.Exec("CREATE TABLE IF NOT EXISTS migration (last_applied_migration INTEGER)")
	var lastMigration int
	result := DBCon.QueryRow("SELECT last_applied_migration from migration")
	err = result.Scan(&lastMigration)
	if err != nil {
		// migrations table is empty, create the single row in the table and set its value to 0
		if err == sql.ErrNoRows {
			_, err = DBCon.Exec("INSERT INTO migration (last_applied_migration) VALUES (0)")
			if err != nil {
				log.Fatal("Could not insert into the migrations table... Exiting...", err)
			}
		} else {
			log.Fatal("Error with the migrations table in the DB... Exiting...", err)
		}

	}

	// apply new migrations, if needed
	var files []os.DirEntry
	files, err = os.ReadDir("./migrations")
	if err != nil {
		log.Fatal("Error trying to read migrations directory... Exiting...", err)
	}
	numMigrationFiles := len(files)
	if numMigrationFiles > lastMigration {
		for i := 0; i < numMigrationFiles-lastMigration; i += 1 {
			migrationFileToApply := "migration_" + strconv.Itoa(lastMigration+i+1) + ".sql"
			migrationStatement, err := os.ReadFile(migrationFileToApply)
			if err != nil {
				log.Fatal("Error reading migration file... Exiting...", string(migrationStatement))
			}
			_, err = DBCon.Exec(string(migrationStatement))
		}
	}

	// update last_applied_migration value
	_, err = DBCon.Exec("UPDATE migration SET last_applied_migration = $1", numMigrationFiles)
	if err != nil {
		log.Fatal("Error updating last_applied_migration value... Exiting...", err)
	}

	// verify that there is only one row in the migrations table
	var numRows int
	result = DBCon.QueryRow("SELECT COUNT(*) from migration")
	err = result.Scan(&numRows)
	if err != nil {
		log.Fatal("Error verifying number of rows in migration table... Exiting...", err)
	}
	if numRows != 1 {
		log.Fatal("ERROR... Migration table does not contain exactly 1 row... Exiting...", err)
	}
}

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
		log.Fatal("Error connecting to the DB... Exiting...")
	}

	handleMigrations()

}
