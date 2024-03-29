package DB

import (
	"database/sql"
	"fmt"
	"os"
	"starterProject/customLogger"
	"strconv"

	_ "github.com/lib/pq"
)

var (
	dbHandle *sql.DB
)

// returns the database handle
func GetDB() *sql.DB {
	return dbHandle
}

func handleMigrations() {
	logger := customLogger.GetLogger()
	db := GetDB()

	// DB contains a table called "migration" which contains one row that holds the number of the last applied migration

	// if migration table doesn't exist, create it and set migration number to 0
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS migration (last_applied_migration INTEGER)")
	if err != nil {
		logger.Fatal("Error creating table when handling migrations...", err)
	}
	var lastMigration int
	result := db.QueryRow("SELECT last_applied_migration from migration")
	err = result.Scan(&lastMigration)
	if err != nil {
		// migrations table is empty, create the single row in the table and set its value to 0
		if err == sql.ErrNoRows {
			lastMigration = 0
			_, err = db.Exec("INSERT INTO migration (last_applied_migration) VALUES (0)")
			if err != nil {
				logger.Fatal("Could not insert into the migrations table... Exiting...", err)
			}
		} else {
			logger.Fatal("Error with the migrations table in the DB... Exiting...", err)
		}

	}

	// apply new migrations, if needed
	files, err := os.ReadDir("./DB/migrations")
	if err != nil {
		logger.Fatal("Error trying to read migrations directory... Exiting...", err)
	}
	numMigrationFiles := len(files)
	if numMigrationFiles > lastMigration {
		// start transaction
		tx, err := db.Begin()

		// defer rollback in case anything fails. If commit happens successfully, rollback does nothing
		defer tx.Rollback()

		if err != nil {
			logger.Fatal("Error starting migration transaction... Exiting...", err)
		}

		// put all new sql statements into the transaction
		for i := 0; i < numMigrationFiles-lastMigration; i += 1 {
			migrationFileToApply := "./DB/migrations/migration_" + strconv.Itoa(lastMigration+i+1) + ".sql"
			migrationStatement, err := os.ReadFile(migrationFileToApply)
			if err != nil {
				logger.Fatal("Error reading migration file... Exiting...", err)
			}
			// add statement to transaction
			_, err = tx.Exec(string(migrationStatement))
			if err != nil {
				logger.Fatal("Error executing migration statement... Exiting...", err)
			}
		}

		// update last_applied_migration_value
		result, err := tx.Exec("UPDATE migration SET last_applied_migration = $1", numMigrationFiles)
		if err != nil {
			logger.Fatal("Error updating last_applied_migration value...", err)
		}
		numRowsAffected, err := result.RowsAffected()
		if numRowsAffected != 1 {
			logger.Fatal("Error, UPDATE statement affected more than 1 row... Exiting...", err)
		} else if err != nil {
			logger.Fatal("Error with UPDATE statement... Exiting...", err)
		}

		// Commit changes
		tx.Commit()
	}

	// verify that there is only one row in the migrations table
	var numRows int
	result = db.QueryRow("SELECT COUNT(*) from migration")
	err = result.Scan(&numRows)
	if err != nil {
		logger.Fatal("Error verifying number of rows in migration table... Exiting...", err)
	}
	if numRows != 1 {
		logger.Fatal("ERROR... Migration table does not contain exactly 1 row... Exiting...")
	}
}

func InitDB() {
	logger := customLogger.GetLogger()

	// Connect to the postgres db
	PGUSER := os.Getenv("PGUSER")
	PGPASSWORD := os.Getenv("PGPASSWORD")
	PGDATABASE := os.Getenv("PGDATABASE")
	PGHOST := os.Getenv("PGHOST")
	PGPORT := os.Getenv("PGPORT")

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable host=%s port=%s password=%s",
		PGUSER, PGDATABASE, PGHOST, PGPORT, PGPASSWORD)

	// cannot use walrus operator on DBCon and err declaration below since that creates a local version of DBCon which
	// does not set the global variable
	var err error
	dbHandle, err = sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("Error opening the DB... Exiting...")
	}

	// check that the database can be connected to
	err = dbHandle.Ping()
	if err != nil {
		logger.Fatal("Error pinging the DB... Exiting...", err)
	}

	handleMigrations()

}
