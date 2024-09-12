package main

import (
	"log"
	"time"

	"github.com/swanchain/go-swan-sdk"
)

const (
	// testnet sets sdk environment
	testnet = true

	// apiKey serves for authentication and authorization.
	apiKey = ""

	// privatekey is used for automatic payments when deploying.
	// If not set, task just creates but not deploys, call `PayAndDeployTask` can continue to pay and deploy it.
	privateKey = ""

	// taskUUID need to renew
	taskUUID = ""

	// autoPay whether pay automatically when renew task
	autoPay = false
)

func main() {
	apiClient, err := swan.NewAPIClient(apiKey, testnet)
	if err != nil {
		log.Fatalf("failed to init swan client, error: %v \n", err)
	}

	duration := time.Hour // duration to renew task

	var paidTxHash []string
	if !autoPay {
		txHash, err := apiClient.RenewPayment(taskUUID, duration, privateKey)
		if err != nil {
			log.Fatalf("failed to renew payment, error: %v \n", err)
		}
		paidTxHash = append(paidTxHash, txHash)
	}

	resp, err := apiClient.RenewTask(taskUUID, duration, privateKey, paidTxHash...)
	if err != nil {
		log.Fatalf("failed to renew task, error: %v \n", err)
	}
	log.Printf("task renew completed, task end at: %d \n", resp.Task.EndAt)
}
