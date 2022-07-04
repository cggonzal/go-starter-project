package templates

import (
	"html/template"
	"starterProject/logger"
)

/* NOTE: in order to add a new template there are 3 steps
   1. Add the name to the list of var
   2. Create a new struct for the data that gets passed into the template
   3. Add the file(s) that get parsed to the init templates function using ParseFiles() as shown below
*/

// 1.
// template variable names
var (
	AboutTemplate  *template.Template
	HomeTemplate   *template.Template
	SignUpTemplate *template.Template
	LoginTemplate  *template.Template
	DeleteTemplate *template.Template
)

// 2.
// template data definitions. One per template.
type AboutData struct {
	UserID    string
	UserImage string
}

type HomeData struct {
	UserID    string
	UserImage string
}

type LoginData struct {
	PasswordIncorrect bool
}

type DeleteData struct {
	UserDoesNotExist bool
}

type SignUpData struct {
	UserAlreadyExists bool
}

// 3.
// initialize templates. Store them in global variables so that files don't have to be parsed on every request
func InitTemplates() {
	var err error

	// NOTE: base.html must come first. Templates are "inherited" from left to right
	HomeTemplate, err = template.ParseFiles("templates/base.html", "templates/home.html")
	if err != nil {
		logger.Logger.Fatal(err)
	}

	AboutTemplate, err = template.ParseFiles("templates/base.html", "templates/about.html")
	if err != nil {
		logger.Logger.Fatal(err)
	}

	SignUpTemplate, err = template.ParseFiles("templates/signup.html")
	if err != nil {
		logger.Logger.Fatal(err)
	}

	LoginTemplate, err = template.ParseFiles("templates/login.html")
	if err != nil {
		logger.Logger.Fatal(err)
	}

	DeleteTemplate, err = template.ParseFiles("templates/delete.html")
	if err != nil {
		logger.Logger.Fatal(err)
	}
}
