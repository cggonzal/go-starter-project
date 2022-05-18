package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"starterProject/DB"
	"starterProject/customUser"
	"starterProject/templates"
)

func about(w http.ResponseWriter, r *http.Request) {
	data := templates.AboutData{UserID: "user id", UserImage: "/static/images/test.jpg"}
	templates.AboutTemplate.Execute(w, data)
}

func initDB() {
	// Connect to the postgres db
	RDS_USERNAME := os.Getenv("RDS_USERNAME")
	RDS_PASSWORD := os.Getenv("RDS_PASSWORD")
	RDS_DB_NAME := os.Getenv("RDS_DB_NAME")
	RDS_HOSTNAME := os.Getenv("RDS_HOSTNAME")
	RDS_PORT := os.Getenv("RDS_PORT")

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable host=%s port=%s password=%s",
		RDS_USERNAME, RDS_DB_NAME, RDS_HOSTNAME, RDS_PORT, RDS_PASSWORD)

	var err error
	DB.DBCon, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func main() {
	// serve static files
	http.Handle("/static", http.FileServer(http.Dir("./static")))
	
	// endpoints
	http.HandleFunc("/about", about)
	http.HandleFunc("/login", customUser.Login)
	http.HandleFunc("/signup", customUser.SignUp)
	http.HandleFunc("/logout", customUser.Logout)
	http.HandleFunc("/delete", customUser.Delete)
	http.HandleFunc("/secret", customUser.Secret)

	// initialize database connection
	initDB()

	// initialize user options
	customUser.InitUser()

	// initialize templates
	templates.InitTemplates()

	// start the server on given $PORT
	log.Print("starting app on port ", os.Getenv("PORT"))
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(PORT, nil))
}
