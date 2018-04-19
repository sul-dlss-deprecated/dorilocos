package main

import (
	"flag"
	"log"
	"time"

	apiclient "github.com/sul-dlss-labs/taco/generated/client"
	"github.com/sul-dlss-labs/taco/generated/client/operations"
)

var iterations int
var concurrency int
var host string

func init() {
	flag.IntVar(&iterations, "n", 10000, "Total number of requests to make")
	flag.IntVar(&concurrency, "c", 20, "Total number of concurrent requests")
	flag.StringVar(&host, "s", "taco-demo.sul.stanford.edu", "Server to connect to")
}

func main() {
	flag.Parse()
	benchmark(iterations, concurrency,
		func(client *apiclient.Taco) error {
			return healthCheckRequest(client)
		})
}

func benchmark(n int, concurrency int, method func(client *apiclient.Taco) error) {
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

func worker(n int, results chan<- time.Duration, errors chan<- error, method func(client *apiclient.Taco) error) {
	// create the transport
	transport := apiclient.DefaultTransportConfig().
		WithHost(host)
	// create the API client
	client := apiclient.NewHTTPClientWithConfig(nil, transport)
	//apiKeyHeaderAuth := httptransport.APIKeyAuth("On-Behalf-Of", "header", os.Getenv("API_KEY"))

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

func healthCheckRequest(client *apiclient.Taco) error {
	_, err := client.Operations.HealthCheck(operations.NewHealthCheckParams()) //, apiKeyHeaderAuth)
	if err != nil {
		return err
	}
	// log.Printf("Taco says: %s", resp.Payload.Status)
	return nil
}
