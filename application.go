package main

import (
    "starterProject/DB"
    "starterProject/customUser"
    "starterProject/templates"
    "database/sql"
    "os"
    "net/http"
    "log"
)

func index(w http.ResponseWriter, r *http.Request){
    landing_page, _ := os.ReadFile("static/index.html")
    w.Write(landing_page)
}

func about(w http.ResponseWriter, r *http.Request){
    data := templates.AboutData{UserID: "user id", UserImage: "/static/images/test.jpg"}
    templates.AboutTemplate.Execute(w, data)
}


func initDB(){
    // Connect to the postgres db
    RDS_USERNAME := os.Getenv("RDS_USERNAME")
    RDS_PASSWORD := os.Getenv("RDS_PASSWORD")
    RDS_DB_NAME := os.Getenv("RDS_DB_NAME")
    RDS_HOSTNAME := os.Getenv("RDS_HOSTNAME")
    RDS_PORT := os.Getenv("RDS_PORT")

    connStr := "user=" + RDS_USERNAME + " dbname=" + RDS_DB_NAME + " sslmode=disable" + " host=" + RDS_HOSTNAME + " port=" + RDS_PORT + " password=" + RDS_PASSWORD

    var err error
    DB.DBCon, err = sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
}

func main() {
    // endpoints
    http.HandleFunc("/", index)
    http.HandleFunc("/about", about)
    http.HandleFunc("/login", customUser.Login)
    http.HandleFunc("/signup", customUser.Signup)
    http.HandleFunc("/logout", customUser.Logout)
    http.HandleFunc("/secret", customUser.Secret)

    // initialize database connection
    initDB()

    // initialize user options
    customUser.InitUser()

    // initialize templates
    templates.InitTemplates()

    // start the server on given $PORT
    PORT := ":" + os.Getenv("PORT")
    log.Print("started app on port ", PORT)
    log.Fatal(http.ListenAndServe(PORT, nil))
}
