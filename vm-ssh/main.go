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
	repoUri = "https://github.com/swanchain/awesome-swanchain/tree/main/vm-ssh"
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

	startTimeout := 300 // unit seconds
	resp, err := apiClient.CreateTask(&swan.CreateTaskReq{
		Duration:     time.Hour, // running hours for your deployment
		PrivateKey:   privateKey,
		RepoUri:      repoUri,
		InstanceType: resources[0].Type,
		StartIn:      startTimeout,
	})
	if err != nil {
		log.Fatalf("create task with auto pay and deploy, error: %v \n", err)
	}

	taskUUUID := resp.TaskUuid
	log.Printf("create task with auto-pay, task_uuid: %s, tx_hash: %s \n", taskUUUID, resp.TxHash)

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
