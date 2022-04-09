package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func Do(method, url string) string {
	client := http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Printf("Could not create request to get prices: %s\n", err)
		return ""
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Do() request failed: %s\n", err)
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Do() body parsing failed: %s\n", err)
		return ""
	}

	return string(body)
}
