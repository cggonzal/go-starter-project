package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "logger: ", log.LstdFlags|log.Llongfile)

	resp, _ := http.Get("http://127.0.0.1:8000/")
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status ", resp.StatusCode)
	}

	resp, _ = http.Get("http://127.0.0.1:8000/login")
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status ", resp.StatusCode)
	}

	resp, _ = http.Get("http://127.0.0.1:8000/signup")
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status ", resp.StatusCode)
	}

	// try to login with credentials that don't exist
	fakeCreds := url.Values{"email": {"kbjsdytp2oiqw423ejdas@gmail.com"}, "password": {"pass"}}
	resp, _ = http.PostForm("http://127.0.0.1:8000/login", fakeCreds)
	if resp.StatusCode != http.StatusUnauthorized {
		logger.Fatal("expected status 401 got status: ", resp.StatusCode)
	}

	// try to access secret without being authenticated
	resp, _ = http.Get("http://127.0.0.1:8000/secret")
	if resp.StatusCode != http.StatusForbidden {
		logger.Fatal("expected status 403 got status: ", resp.StatusCode)
	}

	// user credentials to be tested
	creds := url.Values{"email": {"testuser123456789@gmail.com"}, "password": {"testpass"}}

	// create client so it stores cookies that are used for authentication
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// test signing up user
	resp, _ = client.PostForm("http://127.0.0.1:8000/signup", creds)
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// test signing up user again, should fail since emails need to be unique
	resp, _ = client.PostForm("http://127.0.0.1:8000/signup", creds)
	if resp.StatusCode != http.StatusForbidden {
		logger.Fatal("expected status 403 got status: ", resp.StatusCode)
	}

	// test logging in recently created user
	resp, _ = client.PostForm("http://127.0.0.1:8000/login", creds)
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// access secret now that user is authenticated
	resp, _ = client.Get("http://127.0.0.1:8000/secret")
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// log user out
	resp, _ = client.Get("http://127.0.0.1:8000/logout")
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// verify user can't access secret after logout
	resp, _ = client.Get("http://127.0.0.1:8000/secret")
	if resp.StatusCode != http.StatusForbidden {
		logger.Fatal("expected status 403 got status: ", resp.StatusCode)
	}

	// delete user that was created
	resp, _ = client.PostForm("http://127.0.0.1:8000/delete", creds)
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// verify that user was deleted
	resp, _ = client.PostForm("http://127.0.0.1:8000/login", creds)
	if resp.StatusCode != http.StatusUnauthorized {
		logger.Fatal("expected status 401 got status: ", resp.StatusCode)
	}

	// upload a file and verify response
	file, err := os.Open("./tests/test.txt")
	if err != nil {
		logger.Fatal("error opening test file to be uploaded: ", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("filename", "test.txt")
	if err != nil {
		logger.Fatal("error creating form file: ", err)
	}
	io.Copy(part, file)
	writer.Close()

	r, _ := http.NewRequest("POST", "http://127.0.0.1:8000/upload/", body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	resp, _ = client.Do(r)
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// download file that was previously uploaded
	resp, _ = client.Get("http://127.0.0.1:8000/upload/test.txt")
	if resp.StatusCode != http.StatusOK {
		logger.Fatal("expected status 200 got status: ", resp.StatusCode)
	}

	// attempt to download file that doesn't exist
	resp, _ = client.Get("http://127.0.0.1:8000/upload/testsjdfdsfhjbashs13123dfabn.txt")
	if resp.StatusCode != http.StatusNotFound {
		logger.Fatal("expected status 404 got status: ", resp.StatusCode)
	}

	logger.Print("All tests passed!")
}
