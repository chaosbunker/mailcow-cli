package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Response []struct {
	Type      string        `json:"type"`
	Log       []interface{} `json:"log"`
	Msg       []string      `json:"msg"`
	SpamAlias string        `json:"spam_alias"`}

type ErrorResponse struct {
	Type string   `json:"type"`
	Msg  string `json:"msg"`
}

func main() {
	host := flag.String("host", "", "Your mailcow domain, e.g. mail.example.com")
	apiKey := flag.String("api-key", "", "Your mailcow API Key")
	username := flag.String("username", "", "E-Mail address to forward to. Must be a valid mailbox.")
	validity := flag.String("validity", "", "Validity in hours. 1, 6, 24, ...")

	/*baseUrl := fmt.Sprintf("https://%s/api/v1",*host)*/
	// var getEndpoint = baseUrl + "/get/time_limited_aliases"
	addEndpoint := "/add/time_limited_alias"

	flag.Parse()

	requestBody, err := json.Marshal(map[string]string{
		"username": *username,
		"validity": *validity,
	})

	if err != nil {
		log.Fatalln(err)
	}

	timeout := time.Duration((10 * time.Second))
	client := http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/v1%s",*host,addEndpoint), bytes.NewBuffer(requestBody))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", *apiKey)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var checkForError ErrorResponse
	json.Unmarshal(body, &checkForError)

	badResponse := checkForError.Type
	if badResponse != "" {
		fmt.Printf("\n%s: %s\n", checkForError.Type, checkForError.Msg)
	} else {
		var responseObject Response
		json.Unmarshal(body, &responseObject)
		fmt.Printf("\n%s: %s\n", responseObject[0].Type, responseObject[0].SpamAlias)
	}
}