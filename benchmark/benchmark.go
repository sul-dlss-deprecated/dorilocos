package benchmark

import (
	"log"
	"net/http"
	"time"
)

// NewBenchmark creates a new HTTP benchmark
func Run(n int, concurrency int, method func(*http.Client) error) {
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
	if len(errors) > 0 {
		log.Printf("Error report: ")
		for err := range errors {
			log.Printf("\t%s", err)
		}
	}

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
