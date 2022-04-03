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

	// start the server on given $PORT
    PORT := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func initDB(){
	var err error
	// Connect to the postgres db
    DB_USER := os.Getenv("DB_USER")
    DB_NAME := os.Getenv("DB_NAME")

    db_con_string := "user=" + DB_USER + " dbname=" + DB_NAME + " sslmode=disable"

    DB.DBCon, err = sql.Open("postgres", db_con_string)
	if err != nil {
		panic(err)
	}
}
