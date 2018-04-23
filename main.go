package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sul-dlss-labs/dorilocos/benchmark"
)

var iterations int
var concurrency int
var host string
var key string

func init() {
	flag.IntVar(&iterations, "n", 10000, "Total number of requests to make")
	flag.IntVar(&concurrency, "c", 20, "Total number of concurrent requests")
	flag.StringVar(&host, "s", "taco-demo.sul.stanford.edu", "Server to connect to")
	flag.StringVar(&key, "k", "lmcrae@stanford.edu", "API key")

}

func main() {
	flag.Parse()
	log.Println("Health check")
	benchmark.Run(iterations, concurrency,
		func(client *http.Client) error {
			return healthCheckRequest(client)
		})

	log.Println("Deposit resource")
	benchmark.Run(iterations, concurrency,
		func(client *http.Client) error {
			return despositResourceRequest(client)
		})
}

func healthCheckRequest(client *http.Client) error {

	url := fmt.Sprintf("http://%s/v1/healthcheck", host)

	resp, err := client.Get(url)

	if err != nil {
		log.Printf("Taco err: %s", err)

		return err
	}
	switch resp.StatusCode {
	case 200, 201, 204:
		// do nothing
	default:
		return fmt.Errorf("Bad response: %v", resp.StatusCode)
	}

	return nil
}

func despositResourceRequest(client *http.Client) error {
	byt, err := ioutil.ReadFile("examples/request.json")
	if err != nil {
		panic(err)
	}
	url := fmt.Sprintf("http://%s/v1/resource", host)

	req, err := http.NewRequest("POST", url, bytes.NewReader(byt))
	if err != nil {
		panic(err)
	}
	req.Header.Add("On-Behalf-Of", key)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Taco err: %s", err)

		return err
	}
	switch resp.StatusCode {
	case 200, 201, 204:
		// do nothing
	default:
		return fmt.Errorf("Bad response: %v", resp.StatusCode)
	}
	return nil
}
