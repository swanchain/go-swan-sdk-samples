package main

import (
	"log"
	"time"

	"github.com/swanchain/go-swan-sdk"
)

// MusicGen is an example showing how to deploy a lagrange yaml repo & manual pay for it.

const (
	// apiKey serves for authentication and authorization.
	apiKey = ""

	// privatekey is used for automatic payments when deploying.
	// If not set, task just creates but not deploys, call `PayAndDeployTask` can continue to pay and deploy it.
	privateKey = ""

	// wallet is the wallet addr for payment, just use for non-auto-pay
	walletAddr = ""

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
	repoUri = "https://lagrange.computer/spaces/0x231fe9090f4d45413474BDE53a1a0A3Bd5C0ef03/MusicGen/app" // repo with yaml
)

func main() {
	apiClient, err := swan.NewAPIClient(apiKey)
	if err != nil {
		log.Fatalf("failed to init swan client, error: %v \n", err)
	}

	resources, err := apiClient.InstanceResources(true)
	if err != nil {
		log.Fatalf("failed to get instance resources, error: %v \n", err)
	}

	if len(resources) == 0 {
		log.Fatalf("not found available resources")
	}

	startTimeout := 600 // unit seconds
	resp, err := apiClient.CreateTask(&swan.CreateTaskReq{
		Duration:      time.Hour, // running hours for your deployment
		RepoUri:       repoUri,
		InstanceType:  resources[0].Type,
		StartIn:       startTimeout,
		WalletAddress: walletAddr,
	})
	if err != nil {
		log.Fatalf("create task failed, error: %v \n", err)
	}

	taskUUUID := resp.TaskUuid
	log.Printf("create task successfully, task_uuid: %s\n", taskUUUID)

	result, err := apiClient.PayAndDeployTask(taskUUUID, privateKey, time.Hour, resources[0].Type)
	if err != nil {
		log.Fatalf("task %s pay and deploy failed, error: %v \n", taskUUUID, err)
	}
	log.Printf("task start deploying, task_uuid: %s, tx hash: %s\n", taskUUUID, result.TxHash)

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
