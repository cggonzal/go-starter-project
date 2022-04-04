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
    RDS_USERNAME := os.Getenv("RDS_USERNAME")
    RDS_PASSWORD := os.Getenv("RDS_PASSWORD")
    RDS_DB_NAME := os.Getenv("RDS_DB_NAME")
    RDS_HOSTNAME := os.Getenv("RDS_HOSTNAME")
    RDS_PORT := os.Getenv("RDS_PORT")

    connStr := "user=" + RDS_USERNAME + " dbname=" + RDS_DB_NAME + " sslmode=disable" + " host=" + RDS_HOSTNAME + " port=" + RDS_PORT + " password=" + RDS_PASSWORD

    DB.DBCon, err = sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
}
