package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	benchmark(iterations, concurrency,
		func(client *http.Client) error {
			return healthCheckRequest(client)
		})

	log.Println("Deposit resource")
	benchmark(iterations, concurrency,
		func(client *http.Client) error {
			return despositResourceRequest(client)
		})
}

func benchmark(n int, concurrency int, method func(*http.Client) error) {
	results := make(chan time.Duration, concurrency)
	errors := make(chan error, concurrency)

	requestsPerWorker := n / concurrency
	log.Printf("Running %v requests on each of %v different workers", requestsPerWorker, concurrency)

	start := time.Now()
	for w := 1; w <= concurrency; w++ {
		go worker(requestsPerWorker, results, errors, method)
	}

	for a := 1; a <= concurrency; a++ {
		// Wait for each of the workers to report being done.
		<-results
	}
	finish := time.Now()
	elapsed := finish.Sub(start)

	log.Printf("Total time was %s", elapsed)
	close(errors)
	log.Printf("Errors: %v", len(errors))
	log.Printf("%v requests per second", float64(n)/elapsed.Seconds())

}

func worker(n int, results chan<- time.Duration, errors chan<- error, method func(*http.Client) error) {
	client := &http.Client{}

	start := time.Now()
	for i := 0; i < n; i++ {
		err := method(client)
		if err != nil {
			errors <- err
		}
	}
	finish := time.Now()
	elapsed := finish.Sub(start)
	results <- elapsed
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
