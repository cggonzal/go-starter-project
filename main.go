package main

import (
    "starterProject/DB"
    "starterProject/customUser"
)

func main() {
	http.HandleFunc("/signin", customUser.Signin)
	http.HandleFunc("/signup", customUser.Signup)
	// initialize our database connection
	initDB()
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func initDB(){
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
    DB.DBCon, err = sql.Open("postgres", "dbname=mydb sslmode=disable") // TODO: change to use environment variables
	if err != nil {
		panic(err)
	}
}
