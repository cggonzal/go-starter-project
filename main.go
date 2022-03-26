package main

import (
    "database/sql"
    "starterProject/DB"
    "starterProject/customUser"
    "os"
    "net/http"
    "log"
)

func index(w http.ResponseWriter, r *http.Request){
    landing_page, _ := os.ReadFile("index.html")
    w.Write(landing_page)
}

func main() {
	// endpoints
    http.HandleFunc("/login", customUser.Login)
	http.HandleFunc("/signup", customUser.Signup)
    http.HandleFunc("/logout", customUser.Logout)
    http.HandleFunc("/secret", customUser.Secret)
    http.HandleFunc("/", index)

    // initialize database connection
	initDB()

    // initialize user options
    customUser.InitUser()

	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func initDB(){
	var err error
	// Connect to the postgres db
    DB.DBCon, err = sql.Open("postgres", "user=cgg dbname=mytestdb sslmode=disable") // TODO: change to use environment variables
	if err != nil {
		panic(err)
	}
}
