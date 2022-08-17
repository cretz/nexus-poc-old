package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cretz/nexus-poc/sdk/go/nexus/backend/backendpb"
	"github.com/cretz/nexus-poc/sdk/temporal-go/temporalnexus"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	ctx := context.TODO()

	// Create Temporal client
	c, err := client.Dial(client.Options{HostPort: "nexus-backend.example.com:7233"})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Start worker with our workflow
	worker := worker.New(c, "my-parent-task-queue", worker.Options{})
	worker.RegisterWorkflow(MyParentWorkflow)
	if err := worker.Start(); err != nil {
		log.Fatal(err)
	}
	defer worker.Stop()

	// Run the parent workflow and confirm result
	var resultStr string
	run, err := c.ExecuteWorkflow(
		ctx, client.StartWorkflowOptions{TaskQueue: "my-parent-task-queue"}, MyParentWorkflow, "Nexus")
	if err != nil {
		log.Fatal(err)
	} else if err = run.Get(ctx, &resultStr); err != nil {
		log.Fatal(err)
	} else if resultStr != "Howdy, Nexus!" {
		log.Fatal(fmt.Errorf("expected normal greeting, got: %q", resultStr))
	}
}

// Example of a workflow calling ALOs and sync calls
func MyParentWorkflow(ctx workflow.Context, name string) (string, error) {
	// Create client. Note this doesn't take a context. It could technically be
	// created _outside_ of the workflow.
	client := temporalnexus.NewWorkflowClient(temporalnexus.WorkflowClientOptions{Endpoint: "http://localhost:8080"})

	// Start greeting ALO
	alo, err := client.StartALO(ctx, temporalnexus.WorkflowStartALOOptions{
		Name: "greeting",
		ID:   "my-id",
		Arg:  name,
	})
	if err != nil {
		return "", err
	}

	// Load info and confirm it is running
	if info, err := alo.LoadInfo(ctx); err != nil {
		return "", err
	} else if info.Status != backendpb.AloInfo_RUNNING {
		return "", fmt.Errorf("expected running, got status: %v", info.Status)
	}

	// Send a query and confirm it is the default greeting
	var resultStr string
	res, err := client.CallSync(ctx, temporalnexus.WorkflowCallSyncOptions{Name: "greeting/get", ALOID: "my-id"})
	if err != nil {
		return "", err
	} else if err = res.Get(&resultStr); err != nil {
		return "", err
	} else if resultStr != "Hello, Nexus!" {
		return "", fmt.Errorf("expected normal greeting, got: %q", resultStr)
	}

	// Send a signal to update the prefix, then another signal to finish the
	// workflow
	_, err = client.CallSync(ctx, temporalnexus.WorkflowCallSyncOptions{
		Name:  "greeting/update-prefix",
		ALOID: "my-id",
		Arg:   "Howdy",
	})
	if err != nil {
		return "", err
	}
	_, err = client.CallSync(ctx, temporalnexus.WorkflowCallSyncOptions{Name: "greeting/finish", ALOID: "my-id"})
	if err != nil {
		return "", err
	}

	// Now return the result
	if res, err := alo.WaitForResult(ctx); err != nil {
		return "", err
	} else if err = res.Get(&resultStr); err != nil {
		return "", err
	}
	return resultStr, nil
}
