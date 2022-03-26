package customUser

import (
    "net/http"
    "os"
    "github.com/gorilla/sessions"
    "golang.org/x/crypto/bcrypt"
    "database/sql"
    _ "github.com/lib/pq"
    "starterProject/DB"
)

var (
    // key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
    key = []byte(os.Getenv("SECRET_KEY"))

    store = sessions.NewCookieStore(key)

    // NOTE: change as needed
    BCRYPT_COST = 8
)

// user credentials
type Credentials struct {
    Username string
    Password string
}

func InitUser(){
    // set session to end when browser disconnects
    // docs: https://pkg.go.dev/github.com/gorilla/sessions#CookieStore.MaxAge
    // https://pkg.go.dev/github.com/gorilla/sessions#Options
    store.Options = &sessions.Options{MaxAge: 0}
}

func IsAuthenticated(r *http.Request) bool {
    session, _ := store.Get(r, "cookie-name")
    auth, ok := session.Values["authenticated"].(bool)

    if !auth || !ok {
        return false
    }

    return true
}

func Secret(w http.ResponseWriter, r *http.Request) {
    if !IsAuthenticated(r){
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }

    // user is authenticated, write secret message
    w.Write([]byte("Authentication works!"))
}

// sign up function handler                                                                                                
func Signup(w http.ResponseWriter, r *http.Request){
    // Redirect requests that are not POST
    if r.Method != http.MethodPost{
        http.Redirect(w, r, "/", http.StatusFound)
        return
    }

    err := r.ParseForm()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    // store the request body into a new `Credentials` instance                                                 
    creds := &Credentials{Username: r.PostFormValue("username"), Password: r.PostFormValue("password")}

    // hash the password using the bcrypt algorithm                                                               
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), BCRYPT_COST)

    // Next, insert the username, along with the hashed password into the database                                         
    _, err = DB.DBCon.Exec("INSERT INTO users (username, password) values ($1, $2)", creds.Username, string(hashedPassword))
    if err != nil {
        // If there is any issue with inserting into the database, return a 500 error                                      
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // creates the cookie since it does not exist
    session, _ := store.Get(r, "cookie-name")

    // Set user as authenticated
    session.Values["authenticated"] = true
    session.Save(r, w)

    // credentials stored in the database and user was authenticated, now redirect to landing page
    http.Redirect(w, r, "/", http.StatusFound)
}

func Login(w http.ResponseWriter, r *http.Request) {
    // Redirect requests that are not POST
    if r.Method != http.MethodPost{
        http.Redirect(w, r, "/", http.StatusFound)
        return
    }

    err := r.ParseForm()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("404 page not found"))
        return
    }

    // store the request body into a new `Credentials` instance                                                 
    creds := &Credentials{Username: r.PostFormValue("username"), Password: r.PostFormValue("password")}

    // Get the existing entry present in the database for the given username                                               
    result := DB.DBCon.QueryRow("SELECT password FROM users WHERE username=$1", creds.Username)
    if err != nil {
        // If there is an issue with the database, return a 500 error                                                      
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // We create another instance of `Credentials` to store the credentials we get from the database                       
    storedCreds := &Credentials{}

    // Store the obtained password in `storedCreds`                                                                        
    err = result.Scan(&storedCreds.Password)
    if err != nil {
        // If an entry with the username does not exist, send an "Unauthorized"(401) status                                
        if err == sql.ErrNoRows {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        // If the error is of any other type, send a 500 status                                                            
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // Compare the stored hashed password with the hashed version of the password that was received 
    if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
        // If the two passwords don't match, return a 401 status                                                           
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("Login Error"))
        return
    }

    // If we reach this point, that means the users password was correct, so set the user as authenticated
    session, _ := store.Get(r, "cookie-name")
    session.Values["authenticated"] = true
    session.Save(r, w)

    // redirect to landing page since credentials are correct
    http.Redirect(w, r, "/secret", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "cookie-name")

    // Revoke users authentication
    session.Values["authenticated"] = false
    session.Save(r, w)

    http.Redirect(w, r, "/", http.StatusFound)
}
