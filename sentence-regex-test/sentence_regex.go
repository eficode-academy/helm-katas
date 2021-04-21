package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {
	// check that the users has provided an endpoint as the first argument
	if len(os.Args[1:]) < 1 {
		log.Fatalln("Error: You must pass the endpoint to query as the first argument.")
	}
	// parse the endpoint to a variable
	var endpoint = os.Args[1]

	// query the endpoint
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	// parse the response body to a string
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sentence := string(body)

	// create the regex matcher
	var valid_regex = `^[A-Z][a-z]+\ is\ \d+\ years$`
	var valid_response = regexp.MustCompile(valid_regex)

	// check that the received body matches the regex
	if valid_response.MatchString(sentence) {
		log.Println("response: '", sentence, "' is valid.")
		// exit succesfully if matched
		os.Exit(0)
	} else {
		log.Println("response: '", sentence, "' is not valid.")
		// exit unsuccesfully if not matched
		os.Exit(1)
	}
}
