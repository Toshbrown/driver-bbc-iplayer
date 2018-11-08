package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

const recommendationsURL = "https://ibl.api.bbci.co.uk/ibl/v1/user/recommendations?token_source=IDv5&token="

// GetRecommendations will return the logged in users recommendations from BBC Iplayer
func GetRecommendations(authToken string) (string, error) {

	client := &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("GET", recommendationsURL+authToken, nil)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	return string(body), nil

}
