# Starter Project

Contains boilerplate code that supports user authentication through the signup, login, and logout endpoints. Also contains testing infrastructure for those endpoints.

run `export PORT=8000` then `go run main.go` then go to `http://localhost:8000/` to get started

## CustomUsers
contains the code that does user authentication

## Testing
Open a terminal window and run the following from the root directory:
```
source local_env_vars.sh  # manually create a database that matches the variables in this script on first setup
go run application.go
```

Open a second terminal window and run:
```
go run tests/main.go
```
