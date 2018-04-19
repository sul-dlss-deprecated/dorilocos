package main

import (
	"log"
	"os"

	apiclient "github.com/sul-dlss-labs/taco/generated/client"
	"github.com/sul-dlss-labs/taco/generated/client/operations"
)

func main() {
	// create the transport
	// transport := httptransport.New(os.Getenv("TACO_HOST"), "", nil)

	transport := apiclient.DefaultTransportConfig().
		WithHost(os.Getenv("TACO_HOST"))
	// create the API client
	client := apiclient.NewHTTPClientWithConfig(nil, transport)

	log.Println("Starting Test")
	//apiKeyHeaderAuth := httptransport.APIKeyAuth("On-Behalf-Of", "header", os.Getenv("API_KEY"))
	resp, err := client.Operations.HealthCheck(operations.NewHealthCheckParams()) //, apiKeyHeaderAuth)
	if err != nil {
		panic(err)
	}
	log.Printf("Taco says: %s", resp.Payload.Status)
}
