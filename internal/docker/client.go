package docker

import (
	"log"

	"github.com/moby/moby/client"
)

func CreateClient() (*client.Client, error) {
	apiClient, err := client.New(client.FromEnv, client.WithUserAgent("docklens/1.0.0"))

	if err != nil {
		log.Fatal(err)
	}

	return apiClient, nil
}
