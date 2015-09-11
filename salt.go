package main

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/google/go-querystring/query"
)

// Request represents a single request made to the Salt API
type Request struct {
	Client         string `url:"client"`
	Target         string `url:"tgt"`
	Function       string `url:"fun"`
	Arguments      string `url:"arg,omitempty"`
	ExpressionForm string `url:"expr_form,omitempty"`
}

// Salt represents a connection to the Salt API
type Salt struct {
	Hostname string
	client   http.Client
}

// NewSalt helps create a new Salt object
func NewSalt(hostname string) *Salt {
	s := new(Salt)
	s.Hostname = hostname
	cookieJar, _ := cookiejar.New(nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	s.client = http.Client{
		Jar:       cookieJar,
		Transport: tr,
	}
	return s
}

// Login does a POST request against the Salt API to get an authentication cookie.
func (s *Salt) Login(username string, password string, eauth string) error {
	values := url.Values{
		"username": {username},
		"password": {password},
		"eauth":    {eauth},
	}
	resp, err := s.client.PostForm(s.Hostname+"/login", values)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return errors.New("Could not authenticate using provided credentials")
	}
	return nil
}

// Run sends a Request to the Salt API
func (s *Salt) Run(target string, function string, arguments string) (string, error) {
	request := Request{
		Client:    "local",
		Target:    target,
		Function:  function,
		Arguments: arguments,
	}
	values, _ := query.Values(request)

	resp, err := s.client.PostForm(s.Hostname, values)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}
