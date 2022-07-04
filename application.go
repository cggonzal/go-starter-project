package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"starterProject/DB"
	"starterProject/customUser"
	"starterProject/logger"
	"starterProject/templates"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func about(w http.ResponseWriter, r *http.Request) {
	data := templates.AboutData{UserID: "user id", UserImage: "/static/images/test.jpg"}
	templates.AboutTemplate.Execute(w, data)
}

func main() {
	// serve landing page
	http.HandleFunc("/", index)

	// serve static files
	http.Handle("/static/", http.FileServer(http.Dir(".")))

	// endpoints
	http.HandleFunc("/about", about)
	http.HandleFunc("/login", customUser.Login)
	http.HandleFunc("/signup", customUser.SignUp)
	http.HandleFunc("/logout", customUser.Logout)
	http.HandleFunc("/delete", customUser.Delete)
	http.HandleFunc("/secret", customUser.Secret)

	// initialize Logger, this has to come before all other initializations since they use the logger
	logger.Logger = logger.InitLogger()

	// initialize database connection
	DB.InitDB()

	// initialize user options
	customUser.InitUser()

	// initialize templates
	templates.InitTemplates()

	// start the server on given $PORT
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Print("starting app on port ", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
