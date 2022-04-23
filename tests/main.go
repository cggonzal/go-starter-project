package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func main() {
	resp, _ := http.Get("http://127.0.0.1:8000/")
	if resp.StatusCode != http.StatusOK {
		log.Fatal("expected status 200 got status ", resp.StatusCode)
	}

	resp, _ = http.Get("http://127.0.0.1:8000/login")
	if resp.StatusCode != http.StatusOK {
		log.Fatal("expected status 200 got status ", resp.StatusCode)
	}

	resp, _ = http.Get("http://127.0.0.1:8000/signup")
	if resp.StatusCode != http.StatusOK {
		log.Fatal("expected status 200 got status ", resp.StatusCode)
	}

	// try to login with credentials that don't exist
	fakeCreds := url.Values{"email": {"kbjsdytp2oiqw423ejdas@gmail.com"}, "password": {"pass"}}
	resp, _ = http.PostForm("http://127.0.0.1:8000/login", fakeCreds)
	if resp.StatusCode != http.StatusUnauthorized {
		log.Fatal("expected status 401 got status: ", resp.StatusCode)
	}

	// try to access secret without being authenticated
	resp, _ = http.Get("http://127.0.0.1:8000/secret")
	if resp.StatusCode != http.StatusForbidden {
		log.Fatal("expected status 401 got status: ", resp.StatusCode)
	}

	// user credentials to be tested
	creds := url.Values{"email": {"testuser@gmail.com"}, "password": {"testpass"}}

	// create client so it stores cookies that are used for authentication
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// test signing up user
	resp, _ = client.PostForm("http://127.0.0.1:8000/signup", creds)
	if resp.StatusCode != http.StatusOK {
		log.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// test logging in recently created user
	resp, _ = client.PostForm("http://127.0.0.1:8000/login", creds)
	if resp.StatusCode != http.StatusOK {
		log.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// access secret now that user is authenticated
	resp, _ = client.Get("http://127.0.0.1:8000/secret")
	if resp.StatusCode != http.StatusOK {
		log.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// TODO: delete user that was created so it doesn't fill up the db with test data

	// TODO: make test that signs up a user twice and check that the second sign up fails since emails need to be unique

	// log user out
	resp, _ = http.Get("http://127.0.0.1:8000/logout")
	if resp.StatusCode != http.StatusOK {
		log.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	log.Print("All tests passed!")
}
