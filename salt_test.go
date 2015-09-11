package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Salt
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewSalt(server.URL)
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	for key, want := range values {
		if v := r.FormValue(key); v != want {
			t.Errorf("Request parameter %v = %v, want %v", key, v, want)
		}
	}
}

// TestLogin ensures that we call the correct URL with the correct arguments
func TestLogin(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/login",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			v := values{
				"username": "user.name",
				"password": "Password123",
			}
			testFormValues(t, r, v)
			fmt.Fprint(w, `{"return": [{"perms": [".*"], "token": "c8960686c8ab40dd5d1bd9e72f2c550f654dd323"}]}`)
		},
	)

	err := client.Login("user.name", "Password123", "ldap")
	if err != nil {
		t.Error(err)
	}
}

// TestLoginWithInvalidCredentials ensures we handle invalid logins correctly
func TestLoginWithInvalidCredentials(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/login",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			v := values{
				"username": "user.name",
				"password": "Password123",
			}
			testFormValues(t, r, v)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Could not authenticate using provided credentials")
		},
	)
	err := client.Login("user.name", "Password123", "ldap")
	if err == nil {
		t.Errorf("Error is nil, want not nil.")
	}
}

// TestLoginWithParseError ensures we handle responses from the server that have an error
func TestLoginWithParseError(t *testing.T) {
	setup()
	defer teardown()
	client.Hostname = "somethingbad"
	err := client.Login("user.name", "Password123", "ldap")
	if err == nil {
		t.Errorf("Error is %v, want not nil.", err)
	}
}

func TestRun(t *testing.T) {
	setup()
	defer teardown()
	// client.Login("user.name", "Password123", "ldap")
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			v := values{
				"arg": "",
				"tgt": "*",
			}
			testFormValues(t, r, v)
			fmt.Fprint(w, `{"return": []}`)
		},
	)

	resp, err := client.Run("*", "test.ping", "")
	if err != nil {
		t.Errorf("Error is %v, want nil", err)
	}

	if !strings.Contains(resp, "return") {
		t.Errorf("Reponse does not contain return, want return")

	}
}

// TestRunWithParseError ensures we handle responses from the server that have an error
func TestRunWithParseError(t *testing.T) {
	setup()
	defer teardown()
	client.Hostname = "somethingbad"
	_, err := client.Run("user.name", "Password123", "ldap")
	if err == nil {
		t.Errorf("Error is %v, want not nil.", err)
	}
}
