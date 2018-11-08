package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

const signinURL = "https://account.bbc.com/signin"
const nonceRegexp = `(?m)/signin\?nonce=(\b[^&"]*)`

const bbcCookieJarUrl = "https://bbc.co.uk"
const authTokenCookieName = "ckns_atkn"

//Auth will do the required dance to authenticate with the bbc
//this function may stop working as it relies on extracting the
//nonce form the html of the signin page!
func Auth(username string, password string) (string, error) {

	//Create a http client that preserves cookies and follows all 302 requests
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	client := &http.Client{
		Timeout: time.Second * 5,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
		Jar: jar,
	}

	//
	// get nonce
	//
	req, err := http.NewRequest("GET", "https://account.bbc.com/signin", nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	regex := regexp.MustCompile(nonceRegexp)
	matches := regex.FindStringSubmatch(string(body))
	if len(matches) < 1 {
		return "", errors.New("Nonce not found")
	}
	ckns_nonce := matches[1]

	//
	// Get auth token
	//
	formData := url.Values{}
	formData.Add("jsEnabled", "true")
	formData.Add("username", username)
	formData.Add("password", password)
	formData.Add("attempts", "0")
	req, err = http.NewRequest("POST", "https://account.bbc.com/signin?nonce="+ckns_nonce+"&ptrt=https%3A%2F%2Fwww.bbc.co.uk%2F", strings.NewReader(formData.Encode()))
	req.Header.Set("Referer", "https://account.bbc.com/signin")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en-GB,en;q=0.9,pl;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36")
	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	//
	//Get the ckns_atkn token value from the cookie
	//
	u, _ := url.Parse(bbcCookieJarUrl)
	cookies := jar.Cookies(u)
	var token string
	for _, c := range cookies {
		if c.Name == authTokenCookieName {
			token = c.Value
			break
		}
	}

	if token == "" {
		fmt.Println("\n\n\n Error Token not found in cookies bbc.co.uk cookies after POST signin -----------------------------------")
		printCookies(jar.Cookies(u))

		return "", errors.New("Token not found in cookies")
	}

	return token, nil
}

func printCookies(cookies []*http.Cookie) {
	for _, c := range cookies {
		fmt.Println(c.Name, " = ", c.Value)
	}
}
