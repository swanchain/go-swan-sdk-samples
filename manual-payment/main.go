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

	//
	walletAddress = "" // manual payment (required)

	/*
		1. Allowed Repo URIs
			GitHub Repositories: Must be in the format https://github.com/{owner}/{project}/*.
			Lagrange Repositories: Must be in the format https://lagrange.computer/{owner}/{project}/*.
		2. Supported Formats
			Full Repository: https://github.com/{owner}/{project}
			This format points to the root of the repository. You can optionally specify a branch using RepoBranch.

			Specific Branch and Directory: https://github.com/{owner}/{project}/tree/{branch}/{directory}
			This format specifies a branch and a subdirectory within that branch.
		3. Requirements
			Dockerfile or YAML File:
			A Dockerfile or yaml file is required in the specified repository or directory for the deployment to work.
			These files are typically used to define how the application should be built and deployed.
	*/
	repoUri = "https://github.com/swanchain/awesome-swanchain/tree/main/hello_world"
)

func main() {
	apiClient, err := swan.NewAPIClient(apiKey, testnet)
	if err != nil {
		log.Fatalf("failed to init swan client, error: %v \n", err)
	}

	resources, err := apiClient.InstanceResources(true)
	if err != nil {
		log.Fatalf("failed to init swan client, error: %v \n", err)
	}

	if len(resources) == 0 {
		log.Fatalf("not found available resources")
	}

	durationTime := time.Hour // running hours for your deployment
	startTimeout := 300       // unit seconds
	instanceType := resources[0].Type

	resp, err := apiClient.CreateTask(&swan.CreateTaskReq{
		Duration:      durationTime,
		WalletAddress: walletAddress,
		RepoUri:       repoUri,
		InstanceType:  instanceType,
		StartIn:       startTimeout,
	})
	if err != nil {
		log.Fatalf("failed to create task, error: %v \n", err)
	}

	taskUUUID := resp.TaskUuid
	log.Printf("create task successfully, task_uuid: %s, tx_hash: %s \n", taskUUUID, resp.TxHash)

	// estimate the amount to be paid
	needPayAmount, err := apiClient.EstimatePayment(instanceType, durationTime.Seconds())
	if err != nil {
		log.Fatalf("failed to estimate the amount, error: %v \n", err)
	}
	log.Printf("task_uuid: %s, need to pay the amount: %0.4f \n", taskUUUID, needPayAmount)

	// manual payment and deploy the task
	paymentResult, err := apiClient.PayAndDeployTask(taskUUUID, privateKey, durationTime, instanceType)
	if err != nil {
		log.Fatalf("failed to pay and deploy the task, error: %v \n", err)
	}
	log.Printf("task_uuid: %s, payment result: %+v \n", taskUUUID, paymentResult)

	timer := time.NewTimer(time.Second * time.Duration(startTimeout))
	defer timer.Stop()
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			info, err := apiClient.TaskInfo(taskUUUID)
			if err != nil {
				log.Fatalf("failed to get task info, error: %+v \n", err)
			}
			if info.Task.Status == "completed" {
				appUrls, err := apiClient.GetRealUrl(taskUUUID)
				if err != nil {
					log.Fatalf("failed to get app urls, error: %v \n", err)
				}
				log.Printf("app urls: %v", appUrls)
				return
			}
		case <-timer.C:
			info, err := apiClient.TaskInfo(taskUUUID)
			if err != nil {
				log.Fatalf("failed to get task info, error: %v \n", err)
			}
			log.Fatalf("task deployed timeout, task status: %v", info.Task.Status)
		}
	}
}
