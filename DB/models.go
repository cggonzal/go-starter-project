package DB

// Every table in the database corresponds to a struct in this file.
// It is your responsibility to keep this file up to date with the database.

// user credentials
type Users struct {
	Email    string
	Password string
}
