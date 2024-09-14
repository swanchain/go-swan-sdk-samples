package main

import (
	"github.com/swanchain/go-swan-sdk"
	"log"
)

const (
	// apiKey serves for authentication and authorization.
	apiKey = ""

	// taskUUID need to renew
	taskUUID = ""
)

func main() {
	apiClient, err := swan.NewAPIClient(apiKey)
	if err != nil {
		log.Fatalf("failed to init swan client, error: %v \n", err)
	}

	terminateTaskResp, err := apiClient.TerminateTask(taskUUID)
	if err != nil {
		log.Fatalf("failed to terminate task, error: %v \n", err)
	}

	log.Printf("task terminate completed, result: %+v \n", terminateTaskResp)
}
