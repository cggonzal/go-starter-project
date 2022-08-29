package main

import (
	"fmt"
	"net/http"
	"os"
	"starterProject/DB"
	"starterProject/customLogger"
	"starterProject/customUser"
	"starterProject/templates"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func about(w http.ResponseWriter, _ *http.Request) {
	data := templates.AboutData{UserID: "user id", UserImage: "/static/images/test.jpg"}
	templates.AboutTemplate.Execute(w, data)
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// TODO: change GET method to create a template that returns an html page showing all uploaded files with hyper links...
		http.ServeFile(w, r, "./static/upload.html")
		return
	}

	logger := customLogger.GetLogger()

	// get uploaded file
	file, header, err := r.FormFile("filename")
	if err != nil {
		logger.Fatal("error getting file...", err)
	}

	data := make([]byte, header.Size)
	_, err = file.Read(data)
	if err != nil {
		logger.Fatal("error reading file: ", err)
	}

	err = os.WriteFile("./savedFiles/"+header.Filename, data, 0666)
	if err != nil {
		logger.Fatal("Error saving file:", err)
	}
	logger.Printf("saved file to savedFiles/%s", header.Filename)

	// redirect to landing page
	http.Redirect(w, r, "/", http.StatusFound)
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
	http.HandleFunc("/upload", upload)

	// initialize Logger, this has to come before all other initializations since they use the logger
	customLogger.InitLogger()

	// initialize database connection
	DB.InitDB()

	// initialize user options
	customUser.InitUser()

	// initialize templates
	templates.InitTemplates()

	// check if $PORT environment variable is set
	logger := customLogger.GetLogger()
	if os.Getenv("PORT") == "" {
		logger.Fatal("ERROR... No $PORT environment variable set... Exiting...")
	}

	// start the server on given $PORT
	PORT := fmt.Sprintf(":%s", os.Getenv("PORT"))
	logger.Print("starting app on port ", PORT)
	logger.Fatal(http.ListenAndServe(PORT, nil))
}
